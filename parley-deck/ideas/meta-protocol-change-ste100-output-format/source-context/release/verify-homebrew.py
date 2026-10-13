from pathlib import Path
import datetime,hashlib,json,os,subprocess
root=Path(__file__).resolve().parent
def run(args,**kw):return subprocess.check_output(args,text=True,stderr=subprocess.STDOUT,**kw)
env=dict(os.environ,HOMEBREW_NO_INSTALL_CLEANUP='1',HOMEBREW_NO_AUTO_UPDATE='1',PATH='/opt/homebrew/opt/python@3.14/libexec/bin:'+os.environ['PATH'])
tap=Path('/opt/homebrew/Library/Taps/feci/homebrew-parley')
assert not run(['git','status','--porcelain'],cwd=tap).strip()
expected=json.loads((root/'homebrew-publication.json').read_text())['commit']
run(['git','fetch','origin','main'],cwd=tap);run(['git','merge','--ff-only','origin/main'],cwd=tap)
assert run(['git','rev-parse','HEAD'],cwd=tap).strip()==expected
steps=[]
def step(name,args):
 p=subprocess.run(args,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,env=env)
 (root/'logs'/(name+'.log')).write_text(p.stdout)
 steps.append({'name':name,'command':args,'exit':p.returncode,'log':name+'.log'})
 (root/'homebrew-steps.json').write_text(json.dumps(steps,indent=2)+'\n')
 if p.returncode:raise SystemExit(name+' failed; inspect its retained log')
 return p.stdout
for name in ['parley-deck-cli','parley-deck-skill']:
 step('brew-upgrade-'+name,['brew','upgrade','feci/parley/'+name])
 step('brew-test-'+name,['brew','test','feci/parley/'+name])
step('brew-audit',['brew','audit','--strict','feci/parley/parley-deck-cli','feci/parley/parley-deck-skill'])
assert run(['/opt/homebrew/bin/parley','version']).strip()=='parley 1.55.0'
assert run(['/opt/homebrew/bin/parley-deck-skill','--version']).strip()=='2.19.0'
source=root/'skill/payload/package/skills'
installed=Path(run(['brew','--prefix','parley-deck-skill']).strip())/'libexec/skills'
rows=[]
for p in sorted(source.rglob('*')):
 if not p.is_file():continue
 q=installed/p.relative_to(source)
 assert p.read_bytes()==q.read_bytes(),str(q)
 rows.append({'path':str(q),'sha256':hashlib.sha256(q.read_bytes()).hexdigest()})
result={'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'tap_commit':expected,'steps':steps,'versions':{'cli':'1.55.0','skill':'2.19.0'},'brew_tests':'both passed','strict_audit':'passed','skill_source_bytes_match':True,'payload_files':len(rows),'python':'3.14 task-local PATH','checks':rows}
(root/'homebrew-verification.json').write_text(json.dumps(result,indent=2)+'\n')
print('Both formulae installed and tested. Strict audit passed. All skill payload bytes match.',flush=True)
