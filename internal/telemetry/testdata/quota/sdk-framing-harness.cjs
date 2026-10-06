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
const mode=process.argv[2]||'api';
const retained=mode==='retained-retries'||mode==='retained-top-level';
const backoff=mode==='backoff-retries';
const seconds=retained?[176930,176890]:backoff?[176936,176934,176930]:[176930];
const errors=seconds.map(seconds=>{
 const b=JSON.parse(JSON.stringify(body));
 const hours=Math.floor(seconds/3600), minutes=Math.floor(seconds%3600/60), sec=seconds%60;
 b.error.retry_after=seconds;
 b.error.message=b.error.message.replace(/reset after [0-9hms ]+/,`reset after ${hours}h ${minutes}m ${sec}s`);
 const request=mode==='default-console'?{model:'fixture-model',messages:[{role:'user',content:[{type:'text',text:'offline fixture'}]}]}:{model:'fixture-model',messages:[]};
 const e=new context.On({message:b.error.message,url:'https://provider.invalid/v1/messages',requestBodyValues:request,statusCode:429,responseHeaders:{'content-type':'application/json','retry-after':String(seconds)},responseBody:JSON.stringify(b)});
 e.stack=e.stack.split('\n')[0]+'\n    at offlineFixture (sdk-framing-harness.cjs:1:1)';
 return e;
});
let result=errors[0];
if(mode==='retry'||mode==='retained-retries'||backoff) {
 const attempts=mode==='retry'?[result,result,result]:errors;
 result=new context.A_e({message:'Failed after '+attempts.length+' attempts. Last error: '+attempts.at(-1).message,reason:'maxRetriesExceeded',errors:attempts});
 result.stack=result.stack.split('\n')[0]+'\n    at offlineFixture (sdk-framing-harness.cjs:1:1)';
}
process.stderr.write('Source SHA256 '+crypto.createHash('sha256').update(src).digest('hex')+'; SOURCE-DERIVED '+mode+' framing, not native capture\n');
if(mode==='default-console') {
 // Capture the source-located console.error sink's default depth without loading zcode.
 const chunks=[];
 const {Console}=require('node:console');
 const {Writable}=require('node:stream');
 const sink=new Writable({write(chunk,encoding,cb){chunks.push(chunk);cb();}});
 new Console({stdout:sink,stderr:sink}).error(result);
 process.stdout.write(Buffer.concat(chunks));
} else if(mode==='retained-top-level') {
 for(const e of errors) console.log(util.inspect(e,{depth:10,colors:false}));
} else console.log(util.inspect(result,{depth:10,colors:false}));
console.log('Error: Turn execution failed (traceId: source-derived-fixture)');
