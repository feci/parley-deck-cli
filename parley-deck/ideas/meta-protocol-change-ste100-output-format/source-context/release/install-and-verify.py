#!/usr/bin/env python3
"""After attended close only: install exact reviewed package and record all SKILL hashes."""
from pathlib import Path
import datetime,hashlib,json,os,subprocess,tarfile
root=Path(__file__).resolve().parent
worktree=Path('/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100')
impl=worktree/'parley-deck/ideas/meta-protocol-change-ste100-output-format/IMPLEMENTATION.md'
if '\nstatus: complete\n' not in impl.read_text().split('---',2)[1]+'\n':
    raise SystemExit('Attended close is required before runtime installation')
expected=json.loads((root/'skill/payload-verification.json').read_text())
archive=root/'skill/parley-deck-skill-2.19.0.tgz'
assert hashlib.sha256(archive.read_bytes()).hexdigest()==expected['npm_tarball_sha256']
payload=root/'skill/payload'
payload.mkdir(exist_ok=True)
with tarfile.open(archive) as t:
    t.extractall(payload,filter='data')
cli=['/opt/homebrew/bin/parley-deck-skill']
assert subprocess.check_output(cli+['--version'],text=True).strip()=='2.19.0'
os.environ['PATH']='/opt/homebrew/opt/python@3.14/libexec/bin:'+os.environ['PATH']
logs=root/'logs'
def command(name,args):
    p=subprocess.run(cli+args,text=True,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
    (logs/(name+'.json')).write_text(p.stdout)
    (logs/(name+'.stderr')).write_text(p.stderr)
    if p.returncode: raise SystemExit(f'{name} exited {p.returncode}: {p.stderr}')
    return json.loads(p.stdout)
generic=[('hermes-'+p,Path('/Users/tomasfecko/.hermes/profiles')/p/'skills/parley-deck') for p in ['ldx','librade','testprofile']]
generic+=[('opencode-config',Path('/Users/tomasfecko/.config/opencode/skills/parley-deck'))]
inventories=[command('runtime-all-before-release',['status','--target','all','--include-undetected','--json'])]
for name,dest in generic:
    inventories.append(command(name+'-before-release',['status','--target','generic','--dest',str(dest),'--json']))
destinations=[]
for inv in inventories:
    for runtime in inv['runtimeInstalls']:
        for s in runtime['skills']:
            if s['selected']:
                destinations.append(Path(s['dest']))
assert destinations and len(set(destinations))==len(destinations)
managed_count=len(inventories[0]["runtimeInstalls"])
assert managed_count==15, managed_count
backup=root/'runtime-skills-before-2.19.0.tar.gz'
if backup.exists(): raise SystemExit('Backup already exists; inspect before re-running installer')
with tarfile.open(backup,'w:gz') as t:
    for dest in destinations:
        assert dest.is_dir() and not dest.is_symlink(),str(dest)
        t.add(dest,arcname=str(dest).lstrip('/'))
os.chmod(backup,0o600)
command('runtime-all-install',['install','--target','all','--include-undetected','--force','--json'])
for name,dest in generic:
    command(name+'-install',['install','--target','generic','--dest',str(dest),'--force','--json'])
command('runtime-all-doctor',['doctor','--target','all','--include-undetected','--json'])
for name,dest in generic:
    command(name+'-doctor',['doctor','--target','generic','--dest',str(dest),'--json'])
checks=[]
expected_skills={s['skill']:s['sha256'] for s in expected['skills']}
for dest in destinations:
    path=dest/'SKILL.md'
    digest=hashlib.sha256(path.read_bytes()).hexdigest()
    row={'skill':dest.name,'path':str(path),'sha256':digest,'matches':digest==expected_skills[dest.name]}
    checks.append(row)
assert all(x['matches'] for x in checks)
result={'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'version':'2.19.0','managed_targets':managed_count,'generic_targets':len(generic),'skill_files':len(checks),'checks':checks,'backup':str(backup),'backup_sha256':hashlib.sha256(backup.read_bytes()).hexdigest()}
(root/'runtime-verification.json').write_text(json.dumps(result,indent=2)+'\n')
print(f'All {len(checks)} SKILL.md hashes match across {managed_count+len(generic)} targets; doctor passed.')
