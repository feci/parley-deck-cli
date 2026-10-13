from pathlib import Path
import datetime,hashlib,json,os,shutil,subprocess,sys,tarfile
root=Path(__file__).resolve().parent;expected=sys.argv[1];reviewed=sys.argv[2]
source=Path('/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100-skill')
def run(args,**kw):return subprocess.check_output(args,text=True,stderr=subprocess.STDOUT,**kw)
assert run(['git','rev-parse','HEAD'],cwd=source).strip()==expected
assert not run(['git','status','--porcelain'],cwd=source).strip()
assert not run(['git','diff',reviewed,expected,'--'],cwd=source).strip()
assert json.loads((source/'package.json').read_text())['version']=='2.19.0'
dest=root/'skill';dest.mkdir(exist_ok=True);archive=dest/'parley-deck-skill-2.19.0.tgz';assert not archive.exists()
output=run(['npm','pack','--pack-destination',str(dest),'--json'],cwd=source);(root/'logs/skill-pack.log').write_text(output)
assert archive.is_file()
payload=dest/'payload';payload.mkdir(exist_ok=True)
with tarfile.open(archive) as t:t.extractall(payload,filter='data')
pkg=payload/'package';files=[]
for f in sorted(pkg.rglob('*')):
 if not f.is_file():continue
 rel=f.relative_to(pkg);raw=f.read_bytes();assert raw==(source/rel).read_bytes(),str(rel)
 files.append({'path':str(rel),'sha256':hashlib.sha256(raw).hexdigest(),'bytes':len(raw)})
assert any(x['path']=='skills/parley-deck/references/STE100.md' for x in files)
skills=[{'skill':p.parent.name,'sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'bytes':p.stat().st_size} for p in sorted((pkg/'skills').glob('*/SKILL.md'))];assert len(skills)==6
portable=source/'dist/parley-deck-skill-v2.19.0-macos-arm64';assert portable.exists()
assert run([str(portable),'--version']).strip()=='2.19.0'
shutil.copy2(portable,dest/portable.name)
row={'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'version':'2.19.0','merge_commit':expected,'review_candidate':reviewed,'npm_tarball_sha256':hashlib.sha256(archive.read_bytes()).hexdigest(),'skills':skills,'files':len(files),'portable':{'name':portable.name,'sha256':hashlib.sha256(portable.read_bytes()).hexdigest()}}
(dest/'payload-verification.json').write_text(json.dumps(row,indent=2)+'\n');(dest/'payload-files.json').write_text(json.dumps(files,indent=2)+'\n');print(json.dumps(row,indent=2))
