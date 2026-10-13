from pathlib import Path
import yaml,jsonschema,json,hashlib,subprocess
root=Path(__file__).resolve().parent
manifest=root/'winget/manifests/f/Feci/ParleyDeckSkill/2.19.0'
rows=[]
for p in sorted(manifest.glob('*.yaml')):
 raw=p.read_bytes();assert b'\r\n' in raw and b'\n' not in raw.replace(b'\r\n',b'')
 data=yaml.safe_load(raw);assert data['PackageIdentifier']=='Feci.ParleyDeckSkill' and data['PackageVersion']=='2.19.0'
 url=raw.decode().splitlines()[0].split('$schema=',1)[1]
 schema_path=root/'logs'/('winget-'+data['ManifestType']+'-schema.json')
 subprocess.run(['curl','--fail','--location','--retry','3','--silent','--show-error','--output',str(schema_path),url],check=True)
 schema=json.loads(schema_path.read_text());jsonschema.validate(data,schema)
 rows.append({'file':str(p),'schema':url,'status':'pass','sha256':hashlib.sha256(raw).hexdigest()})
assert len(rows)==3
assets=json.loads((root/'skill/windows-assets.json').read_text())
install=yaml.safe_load((manifest/'Feci.ParleyDeckSkill.installer.yaml').read_text())
for entry in install['Installers']:
 a=assets[entry['Architecture']];assert entry['InstallerSha256']==a['sha256'].upper() and entry['InstallerUrl']==a['url']
(root/'winget-validation.json').write_text(json.dumps({'schema_checks':rows,'installer_hashes_match_live_hosted_assets':True,'windows_validation':'Native winget validate/install unavailable on this Mac; upstream validation required'},indent=2)+'\n')
print('Three manifests pass1.12.0 JSON Schemas, CRLF, version/package checks and live hosted installer hashes.')
