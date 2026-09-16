from pathlib import Path
import sys,json,os,tempfile,subprocess,time,threading,http.server,hashlib,shutil,datetime
case=sys.argv[1];assert case in ['steps2','steps3','retry1','retry2'];out=Path(__file__).resolve().parent/('kimi-local-cap-'+case);out.mkdir(exist_ok=False)
native=Path(tempfile.mkdtemp(prefix='parley-kimi-cap-local-'));workspace=native/'workspace';workspace.mkdir();data=native/'data';data.mkdir();skills=native/'empty-skills';skills.mkdir();(workspace/'fixture.txt').write_text('local cap fixture\n');agent=native/'fixture-agent.md';agent.write_text('---\nname: local-cap-fixture\ndescription: Deterministic local transport fixture with no real provider\ntools: [Read]\nsubagents: []\n---\nYou are a deterministic local transport test. Use only the supplied fixture.\n')
requests=[]
class Handler(http.server.BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_POST(self):
  raw=self.rfile.read(int(self.headers.get('Content-Length','0')));body=json.loads(raw);requests.append({'path':self.path,'body':body})
  if case.startswith('retry'):
   encoded=json.dumps({'error':{'message':'controlled local transient error','type':'internal_server_error'}}).encode();self.send_response(500);self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(encoded)));self.end_headers();self.wfile.write(encoded);return
  payload={'id':'local-fixture','object':'chat.completion','created':0,'model':'local-cap-fixture','choices':[{'index':0,'message':{'role':'assistant','content':None,'tool_calls':[{'id':'local-read-'+str(len(requests)),'type':'function','function':{'name':'Read','arguments':json.dumps({'path':'fixture.txt'})}}]},'finish_reason':'tool_calls'}],'usage':{'prompt_tokens':1,'completion_tokens':1,'total_tokens':2}}
  if body.get('stream'):
   self.send_response(200);self.send_header('Content-Type','text/event-stream');self.end_headers()
   chunk={'id':'local-fixture','object':'chat.completion.chunk','created':0,'model':'local-cap-fixture','choices':[{'index':0,'delta':{'role':'assistant','tool_calls':[{'index':0,'id':'local-read-'+str(len(requests)),'type':'function','function':{'name':'Read','arguments':json.dumps({'path':'fixture.txt'})}}]},'finish_reason':'tool_calls'}],'usage':payload['usage']};self.wfile.write(('data: '+json.dumps(chunk)+'\n\ndata: [DONE]\n\n').encode())
  else:
   encoded=json.dumps(payload).encode();self.send_response(200);self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(encoded)));self.end_headers();self.wfile.write(encoded)
server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler);port=server.server_address[1];threading.Thread(target=server.serve_forever,daemon=True).start()
profile=native/'task.sb';profile.write_text('(version 1)\n(allow default)\n(deny network*)\n(allow network-outbound (remote ip "localhost:'+str(port)+'"))\n(deny file-write* (subpath "/Volumes/My Shared Files/AI_WORKSPACE"))\n(deny file-write* (subpath "/Users/tomasfecko"))\n(deny file-write* (subpath "'+str(native.parent.resolve())+'"))\n(allow file-write* (subpath "'+str(native.resolve())+'"))\n')
canary='import socket,sys\ns=socket.socket();s.settimeout(2);s.connect(("127.0.0.1",int(sys.argv[1])));s.close();print("local endpoint allowed")\ns=socket.socket();s.settimeout(2)\ntry:s.connect(("1.1.1.1",443))\nexcept PermissionError:print("external endpoint denied")\nelse:raise RuntimeError("external traffic was not denied")\n'
c=subprocess.run(['/usr/bin/sandbox-exec','-f',str(profile),'/usr/bin/python3','-c',canary,str(port)],capture_output=True);(native/'canary.log').write_bytes(c.stdout+c.stderr)
if c.returncode:raise RuntimeError((c.stdout+c.stderr).decode())
env={k:v for k,v in os.environ.items() if not k.startswith(('KIMI_','MOONSHOT_'))};env.update(KIMI_CODE_HOME=str(data),KIMI_MODEL_NAME='local-cap-fixture',KIMI_MODEL_PROVIDER_TYPE='kimi',KIMI_MODEL_API_KEY='local-fixture-dummy',KIMI_MODEL_BASE_URL='http://127.0.0.1:'+str(port)+'/v1',KIMI_MODEL_MAX_COMPLETION_TOKENS='17',KIMI_LOOP_MAX_STEPS_PER_TURN=('3' if case=='steps3' else '2'),KIMI_LOOP_MAX_ATTEMPTS_PER_STEP=('2' if case=='retry2' else '1'),KIMI_CODE_INFINITE_RETRY='false',KIMI_CODE_BACKGROUND_PRINT_BACKGROUND_MODE='exit',KIMI_CODE_BACKGROUND_KEEP_ALIVE_ON_EXIT='false',KIMI_CODE_NO_AUTO_UPDATE='true',KIMI_DISABLE_TELEMETRY='true',KIMI_CODE_BUILTIN_PRODUCT_SKILLS='false')
cmd=['/usr/bin/sandbox-exec','-f',str(profile),'/Users/tomasfecko/.kimi-code/bin/kimi','--agent-file',str(agent),'--skills-dir',str(skills),'--output-format','stream-json','-p','Return LOCAL_FIXTURE_DONE. This is a local protocol transport fixture.']
start=time.monotonic();timed_out=False
with (native/'stdout.log').open('wb') as stdout,(native/'stderr.log').open('wb') as stderr:
 proc=subprocess.Popen(cmd,cwd=workspace,env=env,stdout=stdout,stderr=stderr,start_new_session=True)
 try:code=proc.wait(timeout=45)
 except subprocess.TimeoutExpired:
  import signal
  timed_out=True;os.killpg(proc.pid,signal.SIGTERM)
  try:code=proc.wait(timeout=5)
  except subprocess.TimeoutExpired:os.killpg(proc.pid,signal.SIGKILL);code=proc.wait()
 for f in [stdout,stderr]:f.flush();os.fsync(f.fileno())
server.shutdown();(native/'requests.json').write_text(json.dumps(requests,indent=2)+'\n')
r={'completed_at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'native':str(native),'exit':code,'timeout':timed_out,'seconds':round(time.monotonic()-start,3),'requests':len(requests),'request_summaries':[{'path':x['path'],'model':x['body'].get('model'),'max_completion_tokens':x['body'].get('max_completion_tokens'),'max_tokens':x['body'].get('max_tokens'),'tools':[t.get('function',{}).get('name') for t in x['body'].get('tools',[])]} for x in requests],'case':case,'observed_tool_result':any(m.get('role')=='tool' and 'local cap fixture' in str(m.get('content')) for x in requests for m in x['body'].get('messages',[])),'scope':'local fake endpoint and CLI wire configuration only; no external provider, no treatment, no billing proof','env':{k:v for k,v in env.items() if k.startswith('KIMI_') and k!='KIMI_MODEL_API_KEY'},'command':cmd}
for name in ['stdout.log','stderr.log','requests.json','canary.log','task.sb','fixture-agent.md']:
 shutil.copy2(native/name,out/name);assert (native/name).read_bytes()==(out/name).read_bytes();r[name+'_sha256']=hashlib.sha256((out/name).read_bytes()).hexdigest()
(out/'result.json').write_text(json.dumps(r,indent=2)+'\n');print(json.dumps(r),flush=True)
