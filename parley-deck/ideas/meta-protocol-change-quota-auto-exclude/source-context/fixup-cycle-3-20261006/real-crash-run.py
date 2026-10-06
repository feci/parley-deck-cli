from pathlib import Path
import datetime, hashlib, json, os, signal, subprocess, sys, time
source=Path.cwd(); evidence=Path(sys.argv[1]).resolve(); evidence.mkdir(parents=True,exist_ok=False)
root=evidence/'workspace'; root.mkdir(); private=evidence/'parley-home'; private.mkdir()
env=os.environ.copy(); env['PARLEY_HOME']=str(private)
guard=source/'.parley-runtime/quota-implementation/fixup-2/host/no-provider-bin'
env['PATH']=str(guard)+os.pathsep+env['PATH']
subprocess.run(['go','run','./.parley-runtime/quota-implementation/fixup-3/real-crash/init.go',str(root)],env=env,check=True)
worker=evidence/'local-waiter.py'; worker.write_text('import time\ntime.sleep(900)\n')
stub=evidence/'local-agent.sh';stub.write_text('#!/bin/sh\nif [ "$1" = "--version" ]; then echo local-crash-stub; exit 0; fi\nexec '+sys.executable+' '+"'"+str(worker)+"'"+'\n');stub.chmod(0o700)
config=''
for name in ['crash-alpha','crash-beta']:
 config+='[agents.'+name+']\ncommand = '+json.dumps(str(stub))+'\nprompt_mode = "stdin"\nexternal_backend = "local"\ntimeout_ms = 900000\n'
(root/'parley-deck/agents.local.toml').write_text(config)
binary=source/'.parley-runtime/quota-implementation/fixup-3/host/parley'
args=[str(binary),'run','--dir',str(root),'--no-tui','--no-auto','--no-ping','--no-preflight','--yes','--quota-auto-exclude=true','--participants','crash-alpha,crash-beta','Native supervisor crash audit']
(evidence/'command.json').write_text(json.dumps(args,indent=2)+'\n')
stream=(evidence/'launch.log').open('wb'); proc=subprocess.Popen(args,env=env,stdout=stream,stderr=subprocess.STDOUT,start_new_session=True)
records=[]; writer_pids=[]; deadline=time.monotonic()+40
try:
 while time.monotonic()<deadline:
  records=[]
  for p in (root/'.parley-runtime/invocations').glob('*/started.json'):
   try:d=json.loads(p.read_text())
   except (ValueError,OSError):continue
   # Inspect the actual emitted record below; its writer identity must bind to our process.
   if d.get('writer',{}).get('supervisor_pid')==proc.pid: records.append((p,d))
  if len(records)>=2:break
  if proc.poll() is not None:raise RuntimeError('parley exited before two local writers: '+str(proc.returncode))
  time.sleep(.05)
 if len(records)!=2:raise RuntimeError('did not observe exactly two owned native writer records')
 writer_pids=[]
 for p,d in records:
  pid=d['pid'];w=d['writer']; assert pid==w['process_group'],d
  until=time.monotonic()+5
  while True:
   cmd=subprocess.check_output(['ps','-p',str(pid),'-o','command='],text=True)
   if str(worker) in cmd: break
   assert str(stub) in cmd,cmd
   if time.monotonic()>=until: raise RuntimeError('local stub did not exec its worker: '+cmd)
   time.sleep(.02)
  assert os.getpgid(pid)==pid
  writer_pids.append(pid)
 # Kill only this supervisor and the two verified local stub process groups it created.
 proc.kill();proc.wait(timeout=10)
 for pid in writer_pids:os.killpg(pid,signal.SIGKILL)
 for pid in writer_pids:
  until=time.monotonic()+10
  while time.monotonic()<until:
   try:os.kill(pid,0)
   except ProcessLookupError:break
   time.sleep(.05)
 for p,d in records:assert not (p.parent/'terminal.json').exists(),'unexpected normal terminal'
 manifest_paths=list((root/'parley-deck/runs').glob('*/run.json'))
 assert len(manifest_paths)==1,manifest_paths
 manifest=json.loads(manifest_paths[0].read_text());run=manifest_paths[0].parent.name
 slug=manifest['idea_slug']; recover=[str(binary),'quota','recover','--dir',str(root),'--idea',slug,'--run',run]
 results=[];snapshots=[]
 for i in range(2):
  r=subprocess.run(recover,env=env,capture_output=True,text=True,timeout=60)
  (evidence/f'recover-{i+1}.log').write_text(r.stdout+r.stderr)
  results.append({'command':recover,'exit':r.returncode})
  assert r.returncode==0,r.stdout+r.stderr
  receipts={str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in (root/'.parley-runtime/invocations').glob('*/crash-settlement.json')}
  assert len(receipts)==2,receipts
  snapshots.append(receipts)
  for p,d in records:assert not (p.parent/'terminal.json').exists(),'fabricated terminal'
 assert snapshots[0]==snapshots[1],'settlement replay changed bytes'
 (evidence/'result.json').write_text(json.dumps({'utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'supervisor_pid':proc.pid,'writer_pids':writer_pids,'results':results,'receipt_hashes':snapshots[0],'normal_terminals_created':False},indent=2)+'\n')
 print(json.dumps({'status':'passed','evidence':str(evidence),'writers':len(records)}))
finally:
 if proc.poll() is None:proc.kill();proc.wait(timeout=10)
 # Clean up only local workers with a persisted binding to this supervisor and
 # their exact private script in ps. Never signal a guessed or reused PID/group.
 for p in (root/'.parley-runtime/invocations').glob('*/started.json'):
  try:
   d=json.loads(p.read_text()); w=d.get('writer',{}); pid=d['pid']
   if w.get('supervisor_pid')!=proc.pid or w.get('process_group')!=pid: continue
   cmd=subprocess.check_output(['ps','-p',str(pid),'-o','command='],text=True,stderr=subprocess.DEVNULL)
   if (str(worker) in cmd or str(stub) in cmd) and os.getpgid(pid)==pid: os.killpg(pid,signal.SIGKILL)
  except (OSError,ValueError,KeyError,subprocess.CalledProcessError): pass
 stream.close()
