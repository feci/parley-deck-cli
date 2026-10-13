from pathlib import Path
import subprocess,json,hashlib
root=Path('/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100');skill=root.parent/'ste100-skill';delivery=Path(__file__).resolve().parent;idea=root/'parley-deck/ideas/meta-protocol-change-ste100-output-format'
sha=lambda b:hashlib.sha256(b).hexdigest()
h=(idea/'source-context/protocol-hunks.json').read_bytes();assert sha(h)=='d304f0a6b944c272ab9b52debf3f52cc570f5977aa930f0ddb215f0a2cbf494b';h=json.loads(h)
rows=[]
for repo,base,name in [(root,'d51e1c79e0f940146e752e6ece12c17989cc98d3','parley-deck/COOPERATION.md'),(root,'d51e1c79e0f940146e752e6ece12c17989cc98d3','internal/protocol/defaults/COOPERATION.md'),(skill,'054954458dbfe2d764bf6234a6d6bbeee12e2494','skills/parley-deck/references/COOPERATION.md')]:
 before=subprocess.check_output(['git','show',base+':'+name],cwd=repo).decode();after=(repo/name).read_text();assert before.count(h['rule6']['old'])==1
 expected=before.replace(h['rule6']['old'],h['rule6']['new'])+'\n'+h['appendix'];assert after==expected,name
 rows.append({'path':str(repo/name),'exact_two_hunks':True,'sha256':sha(after.encode())})
before=(Path.home()/'.parley/staging/COOPERATION-2.18.0.md').read_bytes();assert sha(before)==h['core_base_sha256'];core=Path.home()/'.parley/staging/COOPERATION-2.19.0.md';assert core.read_text()==before.decode().replace(h['rule6']['old'],h['rule6']['new'])+'\n'+h['appendix']
assert (root/'parley-deck/COOPERATION.md').read_bytes()==(skill/'skills/parley-deck/references/COOPERATION.md').read_bytes()
assert (skill/'skills/parley-deck/references/STE100.md').read_text().split('\n\n',1)[1]==h['appendix'].split('\n\n',1)[1]
assert (root/'VERSION').read_text().strip()=='1.55.0'
assert '"1.55.0"' in (root/'internal/app/version.go').read_text()
for name in ['package.json','package-lock.json']:assert json.loads((skill/name).read_text())['version']=='2.19.0'
meta=json.loads((root/'parley-deck/meta/version.json').read_text())
for key,p in [('protocolSha256',root/'parley-deck/COOPERATION.md'),('skillSha256',skill/'skills/parley-deck/SKILL.md'),('packagedProtocolSha256',skill/'skills/parley-deck/references/COOPERATION.md'),('compatibilityManifestSha256',skill/'skills/parley-deck/references/compatibility.json')]:assert meta[key]==sha(p.read_bytes()),key
refs=subprocess.check_output(['rg','-l','protocol.OutputTextInstruction','internal/runner','internal/app','--glob','*.go','--glob','!*test.go'],cwd=root,text=True).splitlines()
count=sum((root/f).read_text().count('protocol.OutputTextInstruction') for f in refs);assert count==13,count
p=(root/'internal/protocol/outputtext.go').read_text();assert '`'+h['instruction']+'`' in p
result={'source_head':subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),'skill_head':subprocess.check_output(['git','rev-parse','HEAD'],cwd=skill,text=True).strip(),'protocol_copies':rows,'identical_skill_rule_body':True,'metadata_hashes_match':True,'shared_instruction_sites':count,'staged_core_sha256':sha(core.read_bytes()),'payload_sha256':'d304f0a6b944c272ab9b52debf3f52cc570f5977aa930f0ddb215f0a2cbf494b'}
print(json.dumps(result,indent=2))
