from pathlib import Path
import hashlib,json,subprocess,shutil,os
out=Path('.parley-runtime/quota-implementation/fixup-2/checks');p=Path('/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs');raw=p.read_bytes();s=raw.decode();sha=lambda b:hashlib.sha256(b).hexdigest()
slices={'stream-default-error-sink':(4900280,4901600),'stream-options':(7119632,7121400),'undefined-option-filter':(7123374,7123630),'adapter-call':(7167600,7171000),'adapter-runtime':(7185750,7187500),'own-retry-defaults':(4442400,4445100),'sdk-retry':(4854200,4857100),'stream-retry':(7179790,7185860),'run-error-sinks':(13079000,13086800),'console-constructor':(12436500,12439000),'console-assignment':(7542100,7543000),'depth-unrelated-1':(10659100,10661000),'depth-unrelated-2':(11967800,11969200)}
index={'source':str(p),'source_sha256':sha(raw),'coordinate':'zero-based Unicode character offsets, end exclusive','slices':{}}
for name,(a,b) in slices.items():
 q=out/('g14-'+name+'.txt');q.write_text(s[a:b]);index['slices'][name]={'start':a,'end':b,'sha256':sha(q.read_bytes()),'path':str(q)}
for name,a,b in [('api-class','rwr="vercel.ai.error"',',awr="AI_EmptyResponseBodyError"'),('retry-class','UBr="AI_RetryError"',';s(SP')]:
 i=s.index(a);j=s.index(b,i);q=out/('g14-'+name+'.txt');q.write_text(s[i:j]);index['slices'][name]={'start':i,'end':j,'sha256':sha(q.read_bytes()),'path':str(q)}
index['literal_counts']={n:s.count(n) for n in ['inspect.defaultOptions','inspectOptions','console.error=','depth:null','console.error(q)']}
(out/'g14-extraction.json').write_text(json.dumps(index,indent=2))
node=shutil.which('node');env={'LC_ALL':'C','LANG':'C','TZ':'UTC'};results=[]
for mode in ['retained-retries','retained-top-level','backoff-retries','default-console']:
 cmd=[node,'internal/telemetry/testdata/quota/sdk-framing-harness.cjs',mode];r=subprocess.run(cmd,env=env,capture_output=True)
 q=out/('g14-sterile-'+mode+'.stderr');q.write_bytes(r.stdout);(out/('g14-sterile-'+mode+'.log')).write_bytes(r.stderr)
 fixture=Path('internal/telemetry/testdata/quota/zcode-'+mode+'-source-derived.stderr')
 results.append({'argv':cmd,'env':env,'exit':r.returncode,'stdout_sha256':sha(r.stdout),'fixture_equal':r.stdout==fixture.read_bytes(),'source_label_stderr':r.stderr.decode()})
r=subprocess.run([node,'-e','process.stdout.write(JSON.stringify({node:process.versions.node,inspectDepth:require("node:util").inspect.defaultOptions.depth}))'],env=env,capture_output=True);results.append({'node_default_console_depth':r.stdout.decode(),'exit':r.returncode,'env':env});(out/'g14-sterile-results.json').write_text(json.dumps(results,indent=2))
# Verify committed reviewer sources equal the scripts read and run (except our output-only relocation).
base=Path('parley-deck/ideas/meta-protocol-change-quota-auto-exclude/source-context/review-round-04-relaunch-20261005/independent-probes');matches=[]
for name,src,other in [('zretry','zretry/main.go.txt','.parley-runtime/claude1-r4/zretry/main.go'),('r4extra','r4extra/main.go.txt','.parley-runtime/claude1-r4/r4extra/main.go'),('tail','tail/main.go.txt','.parley-runtime/claude1-r4/tail/main.go'),('rerun','probes/rerun.sh','.parley-runtime/claude1-r4/probes/rerun.sh')]:
 a=(base/src).read_bytes();b=Path(other).read_bytes();matches.append({'name':name,'committed':str(base/src),'runtime':other,'equal':a==b,'committed_sha256':sha(a),'runtime_sha256':sha(b)})
(out/'reviewer-source-comparison.json').write_text(json.dumps(matches,indent=2))
print(json.dumps({'sterile':results,'reviewer_sources':matches},indent=2))
