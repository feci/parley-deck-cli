from pathlib import Path
import datetime,hashlib,json,struct,subprocess,yaml
root=Path(__file__).resolve().parent
hosted=root/'hosted-assets';hosted.mkdir(exist_ok=True)
rows=[]
merged=json.loads((root/'merge-commits.json').read_text())
for repo,tag in [('parley-deck-cli','v1.55.0'),('parley-deck-skill','v2.19.0')]:
 release=json.loads(subprocess.check_output(['gh','api',f'repos/feci/{repo}/releases/tags/{tag}']))
 assert not release['draft'] and not release['prerelease']
 actual_commit=subprocess.check_output(['gh','api',f'repos/feci/{repo}/commits/{tag}','--jq','.sha'],text=True).strip()
 assert actual_commit==merged[repo],(repo,actual_commit,merged[repo])
 (root/(repo+'-live-release.json')).write_text(json.dumps(release,indent=2)+'\n')
 for asset in release['assets']:
  name=asset['name'];target=hosted/name
  assert Path(name).name==name
  subprocess.run(['curl','--fail','--location','--retry','3','--silent','--show-error','--output',str(target),asset['browser_download_url']],check=True)
  digest=hashlib.sha256(target.read_bytes()).hexdigest()
  assert asset['digest']=='sha256:'+digest,(name,asset['digest'],digest)
  local=root/('skill' if repo=='parley-deck-skill' else 'cli-assets')/name
  if local.is_file():assert hashlib.sha256(local.read_bytes()).hexdigest()==digest
  rows.append({'repo':repo,'tag':tag,'name':name,'sha256':digest,'bytes':target.stat().st_size,'url':asset['browser_download_url'],'github_digest_matches':True})
expected=json.loads((root/'cli-assets/sha256.json').read_text())
assert len(expected)==6
for name,digest in expected.items():
 assert hashlib.sha256((hosted/name).read_bytes()).hexdigest()==digest
assert json.loads((hosted/'sha256.json').read_text())==expected
payload=json.loads((root/'skill/payload-verification.json').read_text())
assert hashlib.sha256((hosted/'parley-deck-skill-2.19.0.tgz').read_bytes()).hexdigest()==payload['npm_tarball_sha256']
windows={}
for arch,machine in [('x64',0x8664),('arm64',0xaa64)]:
 name=f'parley-deck-skill-v2.19.0-windows-{arch}.exe';raw=(hosted/name).read_bytes()
 assert raw[:2]==b'MZ'
 offset=struct.unpack_from('<I',raw,0x3c)[0]
 assert raw[offset:offset+4]==b'PE\0\0' and struct.unpack_from('<H',raw,offset+4)[0]==machine
 row=next(x for x in rows if x['name']==name)
 windows[arch]={'url':row['url'],'sha256':row['sha256'],'bytes':len(raw),'pe_machine':hex(machine)}
(root/'skill/windows-assets.json').write_text(json.dumps(windows,indent=2)+'\n')
path=root/'winget/manifests/f/Feci/ParleyDeckSkill/2.19.0/Feci.ParleyDeckSkill.installer.yaml'
body=path.read_text();data=yaml.safe_load(body)
for entry in data['Installers']:
 body=body.replace(entry['InstallerSha256'],windows[entry['Architecture']]['sha256'].upper())
path.write_bytes(body.replace('\r\n','\n').replace('\n','\r\n').encode())
(root/'hosted-assets-verification.json').write_text(json.dumps({'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'assets':rows,'cli_expected_matches':True,'skill_tarball_matches':True,'windows_pe_architectures_match':True},indent=2)+'\n')
print(f'Verified {len(rows)} hosted asset digests; CLI build hashes, skill tarball and both Windows PE architectures match.')
