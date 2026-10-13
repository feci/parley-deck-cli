from pathlib import Path
import datetime,hashlib,json,re,subprocess,urllib.request
root=Path(__file__).resolve().parent;tap=root/'tap-source'
def run(args,**kw):return subprocess.check_output(args,text=True,stderr=subprocess.STDOUT,**kw)
if not tap.exists():run(['git','clone','--depth','1','https://github.com/feci/homebrew-parley.git',str(tap)])
assert not run(['git','status','--porcelain'],cwd=tap).strip()
before=run(['git','rev-parse','HEAD'],cwd=tap).strip();rows=[]
for repo,version,formula in [('parley-deck-cli','1.55.0','parley-deck-cli.rb'),('parley-deck-skill','2.19.0','parley-deck-skill.rb')]:
 url=f'https://github.com/feci/{repo}/archive/refs/tags/v{version}.tar.gz';archive=root/(repo+'-v'+version+'.tar.gz')
 subprocess.run(['curl','--fail','--location','--retry','3','--silent','--show-error','--output',str(archive),url],check=True)
 digest=hashlib.sha256(archive.read_bytes()).hexdigest();p=tap/'Formula'/formula;s=p.read_text()
 s,n=re.subn(r'^  url ".*"$',f'  url "{url}"',s,count=1,flags=re.M);assert n==1
 s,n=re.subn(r'^  sha256 "[0-9a-f]+"$',f'  sha256 "{digest}"',s,count=1,flags=re.M);assert n==1;p.write_text(s)
 rows.append({'formula':formula,'url':url,'sha256':digest})
(root/'logs/homebrew-style.log').write_text(run(['brew','style','Formula/parley-deck-cli.rb','Formula/parley-deck-skill.rb'],cwd=tap))
changes=run(['git','diff','--numstat'],cwd=tap);assert len(changes.strip().splitlines())==2
for line in changes.splitlines():assert line.startswith('2\t2\tFormula/'),line
run(['git','add','Formula/parley-deck-cli.rb','Formula/parley-deck-skill.rb'],cwd=tap)
run(['git','commit','-m','[codex-1] meta-protocol-change-ste100-output-format: update both release formulae'],cwd=tap)
run(['git','push','origin','main'],cwd=tap);head=run(['git','rev-parse','HEAD'],cwd=tap).strip()
assert run(['gh','api','repos/feci/homebrew-parley/commits/main','--jq','.sha']).strip()==head
(root/'homebrew-publication.json').write_text(json.dumps({'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'base':before,'commit':head,'formulae':rows,'style':'pass','changes':'two URL and two hash lines only'},indent=2)+'\n');print(head)
