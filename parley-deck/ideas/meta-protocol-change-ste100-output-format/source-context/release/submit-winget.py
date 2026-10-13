from pathlib import Path
import subprocess,json,base64,datetime,hashlib
root=Path(__file__).resolve().parent
branch='feci-parleydeckskill-2.19.0'
repo='feci/winget-pkgs'
def api(endpoint,body=None):
 args=['gh','api',endpoint]
 if body is not None:args+=['--method','POST','--input','-']
 p=subprocess.run(args,input=json.dumps(body) if body is not None else None,capture_output=True,text=True,check=True)
 return json.loads(p.stdout)
validation=json.loads((root/'winget-validation.json').read_text());assert validation['installer_hashes_match_live_hosted_assets']
workflow=json.loads((root/'skill/portable-workflow.json').read_text());assert workflow['conclusion']=='success' and workflow['status']=='completed'
workflow_id=str(workflow['databaseId'])
existing=json.loads(subprocess.check_output(['gh','pr','list','--repo','microsoft/winget-pkgs','--state','open','--search','Feci.ParleyDeckSkill 2.19.0 in:title','--json','number,url']))
assert not existing,existing
base=api('repos/microsoft/winget-pkgs/git/ref/heads/master')['object']['sha']
commit=api('repos/'+repo+'/git/commits/'+base)
files=sorted((root/'winget/manifests/f/Feci/ParleyDeckSkill/2.19.0').glob('*.yaml'));assert len(files)==3
entries=[]
for p in files:
 raw=p.read_bytes();blob=api('repos/'+repo+'/git/blobs',{'content':base64.b64encode(raw).decode(),'encoding':'base64'})
 entries.append({'path':str(p.relative_to(root/'winget')),'mode':'100644','type':'blob','sha':blob['sha']})
tree=api('repos/'+repo+'/git/trees',{'base_tree':commit['tree']['sha'],'tree':entries})
new=api('repos/'+repo+'/git/commits',{'message':'[codex-1] meta-protocol-change-ste100-output-format: update Feci.ParleyDeckSkill to 2.19.0','tree':tree['sha'],'parents':[base]})
api('repos/'+repo+'/git/refs',{'ref':'refs/heads/'+branch,'sha':new['sha']})
body=f'''## 📖 Description

Update Feci.ParleyDeckSkill to 2.19.0. This PR adds one package/version manifest set with three files.
The installers are the final Windows x64 and arm64 assets from the published skill release.
CLI WinGet remains held and is outside this PR.

The three manifests pass the 1.12.0 schemas, CRLF checks, and package/version checks.
Both hosted installer hashes match the live GitHub asset digests.
The release workflow passed macOS and Ubuntu tests and cross-built both Windows assets.
Its job named “windows” runs on Ubuntu. It is not native Windows test evidence.

Native `winget validate` and `winget install` are unavailable on this Mac.
The upstream Windows validation remains required.

## ✅ Checklist

- [ ] Signed the [Contributor License Agreement](https://cla.opensource.microsoft.com). Existing contributor status is subject to the CLA check.
- [ ] Linked to an issue (if applicable). Not applicable to this routine manifest update.

## 📦 Manifest Checklist

- [x] Checked for another open PR for the same package version; none found.
- [x] This PR modifies one package/version manifest set.
- [ ] Validated locally with `winget validate --manifest <path>`.
- [ ] Tested locally with `winget install --manifest <path>`.
- [x] Manifest conforms to the 1.12.0 schemas.

Release: https://github.com/feci/parley-deck-skill/releases/tag/v2.19.0
Portable build: https://github.com/feci/parley-deck-skill/actions/runs/{workflow_id}
'''
(root/'winget-pr-body.md').write_text(body)
pr=api('repos/microsoft/winget-pkgs/pulls',{'title':'Update: Feci.ParleyDeckSkill to 2.19.0','head':'feci:'+branch,'base':'master','body':body,'maintainer_can_modify':True})
changes=api('repos/microsoft/winget-pkgs/pulls/'+str(pr['number'])+'/files')
assert sorted(x['filename'] for x in changes)==sorted(x['path'] for x in entries)
for e in entries:
 f=next(x for x in changes if x['filename']==e['path']);assert f['sha']==e['sha'] and f['status']=='added'
result={'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'url':pr['html_url'],'number':pr['number'],'state':pr['state'],'branch':branch,'head':new['sha'],'base':base,'files':changes,'exact_three_file_diff_verified':True}
(root/'winget-submission.json').write_text(json.dumps(result,indent=2)+'\n')
print(pr['html_url']);print(new['sha'])
