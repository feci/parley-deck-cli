# Supplemental zcode structured-channel evidence (codex-1, 2026-10-04)

This is implementer source evidence for independent review, not a verification verdict.
The review/round-01/claude-1.md stderr findings remain binding. No source or configured
invocation was changed, no zcode participant or provider was invoked, and no credentials
were read. Source SHA256: `3e3433d90fa502e5d02498dfde6c2090df898331359bcfe5f3dbc9a1d00b685f`.

## Located native record

`/Users/tomasfecko/.zcode/cli/log/zcode-2026-10-02.jsonl:16560` contains an existing
structured `turn.failed` / `core.runtime` record at `2026-10-02T16:33:48.366Z`.
The directory is **log**, singular; round-01 R6 names searches of **logs**.
A scoped `rg -l -F 'Weekly/Monthly Limit Exhausted'` found this file. Only matching
records were inspected; unrelated messages and configurations were not copied.

`zcode-native-turn-failed-scrubbed.json` preserves the JSON keys/nesting of that actual
line. Correlation identifiers were consistently pseudonymized and response-header
values scrubbed. It is reformatted, not byte-identical; no error fields were invented.
The native record says 53h 41m 9s, not AC2's 49h 8m 50s. Its provenance/adapter version
at capture is not established solely by the current installed-source hash. These
limitations require reviewer assessment; the file does not itself prove AC2 complete.

The native error includes `error.cause.context.source=provider`, provider/model identity,
and the upstream message in the cause chain. Its response-body summary lists reset
field names without necessarily preserving their values. Check whether this loses a
contradictory machine reset, which must remain a gate under FINAL.

## Candidate implementation to assess

Retain the configured `zcode --prompt ... --mode yolo --cwd ...` invocation. Capture its
already-emitted JSONL diagnostics in a private per-invocation directory using the
supported `ZCODE_LOG_DIR` environment variable. Bind a terminal `turn.failed` to the
exact invocation/root turn and require source-verified provider attribution and all
stated reset values; nested/tool/unrelated events cannot supply authority. This is a
candidate for FINAL §13.8's adapter support work, not approval to enable a classifier.

Assess whether this remains an implementation choice under FINAL or needs the owner's
transport/scope decision required by MAJOR-1. If it cannot provide all necessary data,
keep the major finding open. Do not manufacture eligibility from the native fixture.
The app-server event mapping also has ids, but it is a larger alternative; no such
transport has been selected or implemented.

## Located source snippets

Offsets refer to the installed vendor/zcode.cjs with the SHA above. Snippets are bounded
and some end in the middle of a function; inspect the original surrounding source.

### logger_dir

UTF-8 byte offset 7542547; character offset 7542547.

```javascript
function c2(e={}){let t=e.minLevel??Yti(e.env),r=!1,o=e.logDir??e.env?.ZCODE_LOG_DIR??Qti(),n=typeof e.console=="object"?e.console.stream:e.console===!0||e.env?.ZCODE_LOG_CONSOLE==="1"?process.stderr:void 0,i=e.redactor??new iH,a=s((u,l={})=>new hCt({category:u,defaultContext:l,getMinLevel:s(()=>t,"getMinLevel"),logDir:o,consoleStream:n,includeErrorStack:e.includeErrorStack,redactor:i}),"create");return{createLogger(u){return a(u)},withContext(u){return a("root",u)},setLevel(u){t=u},getLogDir(){return o},scheduleLogRetentionCleanup(u={}){if(!r)return r=!0,dCt({...u,logDir:o,logger:u.logger??a("zcode",{module:"adapters.logging"})})}}}function 
```

### terminal_logger

UTF-8 byte offset 11352977; character offset 11352971.

```javascript
async function PL(e,t){let{coreError:r,events:o,traceContext:n,turnPhase:i,inputId:a}=t,u=r.type===ye.TurnCancelled,l=ZY(r),c=u?e.createEvent(K.TurnComplete,{response:"",tokenCount:0,usage:oh(o),toolCallCount:0,historyRoundCount:t.historyRoundCount??0,duration:t.durationMs,resultType:"cancelled",inputId:a,...t.preserveQueueAutoDrainOnCancel?{preserveQueueAutoDrainOnCancel:!0}:{}},n):e.createEvent(K.TurnError,{error:{type:l?.code??r.type,...EA(r,t.fallbackMessage),stack:r.stack},turnPhase:i,inputId:a},n);await e.appendEvent(c,n),o.push(c),e.logger?.error(`${t.logLabel} failed`,r,{...be(n),event:t.logEvent,module:"core.runtime",status:u?"cancelled":"failed",turnPhase:i})}function Ia(e,t){if(ZY(t?.reason)||ZY(e))return!1;if(t?.aborted)return!0;let r=e,o=new WeakSet;for(let n=0;n<=6;n+=1){if(r==null||typeof r!="object"||o.has(r))return!1;if(o.add(r),yi(r)&&r.type===ye.TurnCancelled)return!0;
```

### log_serializer

UTF-8 byte offset 7540265; character offset 7540265.

```javascript
function WVr(e,t){let r={timestamp:e.timestamp.toISOString(),level:e.levelName.toLowerCase(),event:e.event,module:e.module,message:e.message,traceId:e.traceId,spanId:e.spanId,parentSpanId:e.parentSpanId,sessionId:e.sessionId,turnId:e.turnId,toolCallId:e.toolCallId,durationMs:e.durationMs,status:e.status,context:t.redact(e.context),error:t.redact(e.error)};return Object.fromEntries(Object.entries(r).filter(([,o])=>o!==void 0))}function HVr(e){let t=e.traceId?` trace=${e.traceId.slice(0,8)}`:"",r=e.event?` event=${e.event}`:"";return`${e.levelName.toLowerCase()} [${e.module??"log"}]${t}${r} ${e.message}`}function KVr(e){return e==="started"||e=
```

### error_chain

UTF-8 byte offset 7540975; character offset 7540975.

```javascript
function mCt(e,t,r=new WeakSet,o=0){if(o>8)return{name:"ErrorCauseDepthLimit",message:"Error cause chain exceeded the serialization depth limit."};if(e===null||typeof e!="object")return{name:"UnknownError",message:String(e)};if(r.has(e))return{name:"ErrorCauseCircularReference",message:"Error cause chain contained a circular reference."};r.add(e);let n=e,i=n.cause,a={name:typeof n.name=="string"?n.name:"Error",message:typeof n.message=="string"?n.message:String(e)},u=VVr(n,"code"),l=VVr(n,"type"),c=n.context;return u&&(a.code=u),l&&(a.type=l),t&&typeof n.stack=="string"&&(a.stack=n.stack),Jti(c)&&(a.context=c),i!==void 0&&(a.cause=mCt(i,t,r,o+1)),a}function VVr(e,t){let r=e[t];return typeof r=="string"&&r.length>0?r:void 0}function Jti(e){if(e===null||typeof e!="object")return!1;let t=Object.getPrototypeOf(e);return t===Object.prototype||t===null}var iH,fCt=S(()=>{"use strict";iH=class{static{s(this,"DefaultLogRedact
```

### framed_write

UTF-8 byte offset 7544602; character offset 7544602.

```javascript
log(t,r,o,n){if(t<this.getMinLevel())return;let i={...this.defaultContext,...n},a=this.createEntry(t,r,i,o),u=WVr(a,this.redactor),l=JSON.stringify(u);try{eri(this.logDir);let c=(0,k2e.join)(this.logDir,tri());Cu({operation:"appendFile",path:c}),(0,Jz.appendFileSync)(c,`${l}
`,"utf8")}catch{}this.consoleStream&&this.consoleStream.write(`${HVr(a)}
`)}createEntry(t,r,o,n){return{timestamp:new Date,level:t,levelName:gtr[t],event:typeof o.event=="string"?o.event:void 0,module:typ
```

### adapter_unwrap

UTF-8 byte offset 7173272; character offset 7173272.

```javascript
N instanceof OP)throw A=!0,N.adapterError;let L=Date.now(),H=!p&&q(N);H&&(y={...y,maxAttempts:e.retry.maxAttempts+1});let ue=Yb(N,e.request.abortSignal),ce=KG({providerId:String(y.
```

### event_envelope

UTF-8 byte offset 12566647; character offset 12566613.

```javascript
function xqn(e,t,r={}){return{deliveryKind:t,eventId:String(e.id),payload:Qca(e),seq:r.seq??e.sequenceNumber,sessionId:String(e.sessionId),timestamp:e.timestamp.getTime(),traceId:String(e.traceId),turnId:e.turnId?String(e.turnId):void 0,type:Lda(e.type)}}funct
```

### payload_projection

UTF-8 byte offset 8519062; character offset 8519056.

```javascript
function EA(e,t="Turn execution failed"){let r=Zcn(e),o=Vcn(r),n=o?.message??jJ(t)??t,i=tbi(o,r),a=rbi(r,n),u=Wxi(r);return{...u?{attribution:u}:{},...i?{code:i}:{},message:n,...a?{detail:a}:{}}}function Wxi(e){let t=e.flatMap(m=>m.con
```

### payload_chain

UTF-8 byte offset 8520542; character offset 8520536.

```javascript
function Zcn(e){let t=[],r=new WeakSet,o=e;for(let n=0;n<12&&!(!Wcn(o)||r.has(o));n+=1){r.add(o),t.push(ebi(o));let i=o.cause??o.lastError??o.error;if(!i||i===o)break;o=i}return t}
```
