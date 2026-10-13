#!/usr/bin/env python3
"""Build six release binaries from the expected clean merged commit."""
from pathlib import Path
import datetime,hashlib,json,os,subprocess,sys
root=Path(__file__).resolve().parent
expected=sys.argv[1]
source=root/'cli-source'
def run(args,**kw):
    return subprocess.check_output(args,text=True,stderr=subprocess.STDOUT,**kw)
if not source.exists():
    run(['git','clone','--depth','1','--branch','main','https://github.com/feci/parley-deck-cli.git',str(source)])
head=run(['git','rev-parse','HEAD'],cwd=source).strip()
assert head==expected,(head,expected)
assert not run(['git','status','--porcelain'],cwd=source).strip()
impl=source/'parley-deck/ideas/meta-protocol-change-ste100-output-format/IMPLEMENTATION.md'
assert '\nstatus: complete\n' in impl.read_text().split('---',2)[1]+'\n'
assets=root/'cli-assets';assets.mkdir(exist_ok=True)
storage={}
storage['TMPDIR']=storage['GOTMPDIR']=str(root/'tmp')
env=dict(os.environ,**storage,CGO_ENABLED='0')
rows=[];hashes={}
for platform in ['darwin','linux','windows']:
    for arch,label in [('arm64','arm64'),('amd64','x64')]:
        name=f'parley-v1.55.0-{platform}-{label}'+('.exe' if platform=='windows' else '')
        target=assets/name
        output=run(['go','build','-trimpath','-ldflags','-s -w','-o',str(target),'./cmd/parley'],cwd=source,env=dict(env,GOOS=platform,GOARCH=arch))
        (root/'logs'/('build-'+platform+'-'+label+'.log')).write_text(output)
        identity=run(['go','version','-m',str(target)])
        assert 'vcs.revision='+expected in identity and 'vcs.modified=false' in identity,identity
        digest=hashlib.sha256(target.read_bytes()).hexdigest();hashes[name]=digest
        rows.append({'name':name,'sha256':digest,'bytes':target.stat().st_size,'vcs':identity})
        print(name,digest,flush=True)
(assets/'sha256.json').write_text(json.dumps(hashes,indent=2)+'\n')
assert run([str(assets/'parley-v1.55.0-darwin-arm64'),'version']).strip()=='parley 1.55.0'
assert not run(['git','status','--porcelain'],cwd=source).strip()
(root/'cli-build-verification.json').write_text(json.dumps({'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'commit':expected,'clean_source':True,'native_version':'parley 1.55.0','assets':rows},indent=2)+'\n')
