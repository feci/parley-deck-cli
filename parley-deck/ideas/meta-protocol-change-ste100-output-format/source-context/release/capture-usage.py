from pathlib import Path
from datetime import datetime,timezone
import json
root=Path('/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/ste100');delivery=Path(__file__).resolve().parent;slug='meta-protocol-change-ste100-output-format'
rows=[];pending=[]
for p in (root/'.parley-runtime/invocations').glob('*/requested.json'):
 d=json.loads(p.read_text())
 if d.get('metadata',{}).get('idea')!=slug:continue
 term=p.with_name('terminal.json')
 if not term.exists():pending.append(d['invocation_id']);continue
 d=json.loads(term.read_text());m=d['metadata'];o=d['outcome']
 rows.append({'id':d['invocation_id'],'agent':m['agent'],'phase':m['phase'],'started':bool(d['started_at']),'started_at':d['started_at'],'completed_at':d['completed_at'],'seconds':d['duration_ms']/1000,'exit':o['exit_code'],'failure':o.get('failure_class'),'usage':o['usage'],'source':str(term)})
rows.sort(key=lambda x:x['started_at'] or x['completed_at']);agents={}
for row in rows:
 a=agents.setdefault(row['agent'],{'child_starts':0,'successes':0,'failed_children':0,'control_plane_refusals':0,'child_seconds':0})
 if row['started']:
  a['child_starts']+=1;a['successes']+=row['exit']==0;a['failed_children']+=row['exit']!=0;a['child_seconds']+=row['seconds']
 else:a['control_plane_refusals']+=1
native=[]
for p in delivery.glob('*-receipt.json'):
 d=json.loads(p.read_text())
 if d.get('kind')=='configured-native-fallback':native.append(d)
sessions=[]
for source in json.loads((delivery/'codex-usage-sources.json').read_text()):
 latest=None
 with Path(source['path']).open() as f:
  for line in f:
   try:d=json.loads(line)
   except ValueError:continue
   payload=d.get('payload',{});info=payload.get('info') or {}
   if d.get('type')=='event_msg' and payload.get('type')=='token_count' and info.get('total_token_usage'):latest={'timestamp':d.get('timestamp'),'usage':info['total_token_usage']}
 source['latest_cumulative_snapshot']=latest;sessions.append(source)
result={'at':datetime.now(timezone.utc).isoformat(),'idea':slug,'scope':'Measured terminal records that explicitly name this idea. Configured native fallback processes are separate. CLI syntax/alias refusals without an invocation are orchestration gaps, not child starts.','measured_agents':agents,'invocations':rows,'pending':pending,'native_processes':native,'codex_client_sessions':sessions,'interpretation':'Client snapshots are cumulative session counters, not a marginal task bill. Cached input and reasoning output are subsets. Do not add snapshots. Native participant tokens and cost remain unavailable, not zero. No reliable dollar total is claimed.'}
(delivery/'usage-final.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps({'agents':agents,'pending':pending,'native_processes':len(native)},indent=2))
