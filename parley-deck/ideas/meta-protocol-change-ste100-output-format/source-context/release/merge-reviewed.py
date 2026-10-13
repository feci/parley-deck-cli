from pathlib import Path
import datetime,json,subprocess
root=Path(__file__).resolve().parent
close=json.loads((root/'source-close-receipt.json').read_text())
assert close['all_four_conditions'] and close['windows_final_recorded']
cli=Path('/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100')
skill=cli.parent/'ste100-skill'
def run(args,**kw):return subprocess.check_output(args,text=True,stderr=subprocess.STDOUT,**kw)
merges={};rows=[]
for repo,pr,source,reviewed in [('parley-deck-cli','84',cli,'99d04b3ecd0740fb834e70ec109c9384c7b7cd4d'),('parley-deck-skill','12',skill,'92e3dc78da2f1c55d35c905cabe57fc69ede0670')]:
 assert not run(['git','status','--porcelain'],cwd=source).strip()
 head=run(['git','rev-parse','HEAD'],cwd=source).strip()
 paths=run(['git','diff','--name-only',reviewed,head],cwd=source).splitlines()
 if repo=='parley-deck-skill':assert not paths,paths
 else:
  allowed=['parley-deck/ideas/meta-protocol-change-ste100-output-format/','parley-deck/inbox/codex-1-to-all_meta-protocol-change-ste100-output-format_','parley-deck/runs/']
  assert all(any(p.startswith(a) for a in allowed) for p in paths),paths
 state=json.loads(run(['gh','pr','view',pr,'--repo','feci/'+repo,'--json','headRefOid,isDraft,state']))
 assert state['headRefOid']==head and state['state']=='OPEN',state
 if state['isDraft']:run(['gh','pr','ready',pr,'--repo','feci/'+repo])
 run(['gh','pr','merge',pr,'--repo','feci/'+repo,'--merge','--match-head-commit',head])
 state=json.loads(run(['gh','pr','view',pr,'--repo','feci/'+repo,'--json','mergeCommit,mergedAt,state,url']))
 assert state['state']=='MERGED',state
 merged=state['mergeCommit']['oid']
 parents=json.loads(run(['gh','api','repos/feci/'+repo+'/commits/'+merged]))['parents'];assert len(parents)==2 and parents[1]['sha']==head
 run(['git','fetch','origin','main'],cwd=source);run(['git','merge','--ff-only','origin/main'],cwd=source)
 assert run(['git','rev-parse','HEAD'],cwd=source).strip()==merged
 assert not run(['git','diff',head,merged,'--'],cwd=source).strip()
 merges[repo]=merged;rows.append({'repo':repo,'reviewed_product':reviewed,'pr_head':head,'merge_commit':merged,'url':state['url'],'tree_identical_to_pr_head':True,'mergedAt':state['mergedAt']})
 (root/'merge-commits.json').write_text(json.dumps(merges,indent=2)+'\n')
 (root/'merge-verification.json').write_text(json.dumps({'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'repositories':rows},indent=2)+'\n')
 print(repo,merged,flush=True)
