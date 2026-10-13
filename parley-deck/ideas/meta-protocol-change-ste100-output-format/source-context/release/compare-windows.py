from pathlib import Path
import datetime,hashlib,json,re,subprocess
root=Path(__file__).resolve().parent
source=Path('/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100')
base=Path('/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-10-09-driver-unstall/logs/windows-fixup-pr.log')
pattern=r'--- FAIL: ([^\s]+)'
base_bytes=base.read_bytes();base_names=set(re.findall(pattern,base_bytes.decode(errors='replace')))
assert base_names
product_diff=subprocess.check_output(['git','diff','--name-only','134ac40cd178ebe0318a838d9357c4e6f561f925','d51e1c79e0f940146e752e6ece12c17989cc98d3','--','.',' :(exclude)parley-deck/ideas'.strip(),':(exclude)parley-deck/inbox',':(exclude)parley-deck/runs'],cwd=source,text=True)
assert not product_diff.strip(),product_diff
rows=[]
for run,job in [('38406035407','115325766440'),('38406028243','115325741335')]:
 state=json.loads(subprocess.check_output(['gh','api',f'repos/feci/parley-deck-cli/actions/jobs/{job}']))
 assert state['status']=='completed',state['status']
 p=root/'logs'/('windows-99d04b3-'+job+'.log')
 data=subprocess.check_output(['gh','run','view',run,'--repo','feci/parley-deck-cli','--job',job,'--log'])
 p.write_bytes(data);text=data.decode(errors='replace');names=set(re.findall(pattern,text))
 added=sorted(names-base_names)
 rows.append({'run':run,'job':job,'conclusion':state['conclusion'],'url':state['html_url'],'log':str(p),'sha256':hashlib.sha256(data).hexdigest(),'direct_failing_names':len(names),'added_names':added,'removed_names':sorted(base_names-names),'timeout':bool(re.search(r'panic: test timed out',text)),'evidence_build_failure':'internal/evidence [build failed]' in text,'new_policy_or_transport_failures':sorted(n for n in names if any(x in n for x in ['OutputText','LargeHeadlessPrompt','HeadlessArgumentBounds','LargeStdin','LargePromptDelivery','LiveSTEAppendix']))})
result={'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'candidate':'99d04b3ecd0740fb834e70ec109c9384c7b7cd4d','baseline_product':'134ac40cd178ebe0318a838d9357c4e6f561f925','base_product_equal_to_design_baseline':True,'baseline_log':str(base),'baseline_sha256':hashlib.sha256(base_bytes).hexdigest(),'baseline_direct_failing_names':len(base_names),'extraction_pattern':pattern,'jobs':rows,'limits':'Failure-name comparison does not establish equal causes or prove tests ran before timeout. Native Windows remains unresolved and experimental. CLI WinGet stays held.'}
(root/'windows-final-comparison.json').write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps(result,indent=2))
