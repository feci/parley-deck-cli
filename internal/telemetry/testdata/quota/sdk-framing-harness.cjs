// Offline, SOURCE-DERIVED framing harness. NOT a captured zcode invocation.
// Evaluates only the installed SDK's located error class definitions in an
// isolated VM. No SDK bootstrap, credentials, provider configuration or network.
const fs = require('node:fs');
const vm = require('node:vm');
const util = require('node:util');
const crypto = require('node:crypto');
const sourcePath = '/opt/homebrew/lib/node_modules/zcode-app-cli/vendor/zcode.cjs';
const src = fs.readFileSync(sourcePath, 'utf8');
function section(a,b) { const i=src.indexOf(a), j=src.indexOf(b,i); if(i<0||j<0)throw Error('installed SDK changed');return src.slice(i,j); }
const api=section('rwr="vercel.ai.error"', ',awr="AI_EmptyResponseBodyError"');
const retry=section('UBr="AI_RetryError"',';s(SP');
const context=vm.createContext({});
vm.runInContext('var Pbr,Cbr,Rbr,Abr,qBr; var s=(x,n)=>Object.defineProperty(x,"name",{value:n,configurable:true});'+api+';'+retry+';',context);
const body=JSON.parse(fs.readFileSync(__dirname+'/zcode-recorded-response.json','utf8'));
const message=body.error.message;
const make=()=>new context.On({message,url:'https://provider.invalid/v1/messages',requestBodyValues:{model:'fixture-model',messages:[]},statusCode:429,responseHeaders:{'content-type':'application/json','retry-after':String(body.error.retry_after)},responseBody:JSON.stringify(body)});
// Deterministic source-derived stack location; actual error object/framing is SDK.
const e=make();e.stack=e.stack.split('\n')[0]+'\n    at offlineFixture (sdk-framing-harness.cjs:1:1)';
const mode=process.argv[2]||'api';
let result=e;
if(mode==='retry') {result=new context.A_e({message:'Failed after 3 attempts. Last error: '+message,reason:'maxRetriesExceeded',errors:[e,e,e]});result.stack=result.stack.split('\n')[0]+'\n    at offlineFixture (sdk-framing-harness.cjs:1:1)';}
process.stderr.write('Source SHA256 '+crypto.createHash('sha256').update(src).digest('hex')+'; SOURCE-DERIVED '+mode+' framing, not native capture\n');
console.log(util.inspect(result,{depth:10,colors:false}));
console.log('Error: Turn execution failed (traceId: source-derived-fixture)');
