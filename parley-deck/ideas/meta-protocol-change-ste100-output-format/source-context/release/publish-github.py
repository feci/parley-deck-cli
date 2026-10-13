from pathlib import Path
import json, subprocess, datetime
root=Path(__file__).resolve().parent
merges=json.loads((root/'merge-commits.json').read_text())
cli=json.loads((root/'cli-build-verification.json').read_text())
skill=json.loads((root/'skill/payload-verification.json').read_text())
assert cli['commit']==merges['parley-deck-cli'] and cli['clean_source']
assert skill['merge_commit']==merges['parley-deck-skill']
rows=[]
for repo,tag,title,notes,assets in [
 ('parley-deck-cli','v1.55.0','Parley Deck CLI 1.55.0',root/'cli-release-notes.md',sorted((root/'cli-assets').iterdir())),
 ('parley-deck-skill','v2.19.0','Parley Deck Skill 2.19.0',root/'skill-release-notes.md',[root/'skill/parley-deck-skill-2.19.0.tgz',root/'skill/parley-deck-skill-v2.19.0-macos-arm64'])
]:
 assert all(p.is_file() for p in assets)
 if repo=='parley-deck-cli':assert len(assets)==7
 remote='feci/'+repo
 exists=subprocess.run(['gh','release','view',tag,'--repo',remote,'--json','tagName'],capture_output=True,text=True)
 if exists.returncode==0:raise SystemExit('Release exists; inspect before any resume: '+remote+' '+tag)
 assert 'release not found' in exists.stderr,exists.stderr
 cmd=['gh','release','create',tag,*map(str,assets),'--repo',remote,'--target',merges[repo],'--title',title,'--notes-file',str(notes)]
 result=subprocess.check_output(cmd,text=True,stderr=subprocess.STDOUT)
 release=json.loads(subprocess.check_output(['gh','release','view',tag,'--repo',remote,'--json','tagName,url,publishedAt,isDraft,isPrerelease,targetCommitish,assets']))
 assert not release['isDraft'] and not release['isPrerelease']
 assert subprocess.check_output(['gh','api','repos/'+remote+'/commits/'+tag,'--jq','.sha'],text=True).strip()==merges[repo]
 rows.append(release)
 (root/(repo+'-published-release.json')).write_text(json.dumps(release,indent=2)+'\n')
 print(result,flush=True)
(root/'github-publication.json').write_text(json.dumps({'at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'releases':rows},indent=2)+'\n')
