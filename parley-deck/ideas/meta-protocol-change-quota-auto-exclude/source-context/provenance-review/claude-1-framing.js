const util=require('util');
const m=Symbol.for('vercel.ai.error'), mc=Symbol.for('vercel.ai.error.AI_APICallError');
class APICallError extends Error{constructor(msg,o){super(msg);this.name='AI_APICallError';Object.assign(this,o);this[m]=true;this[mc]=true;}}
const body='{"error":{"message":"[glm/glm-5.3] [429]: Weekly/Monthly Limit Exhausted. Your limit will reset at 2026-10-05 06:14:57 (reset after 49h 8m 50s)","retry_after":176930,"reset_at":"2026-10-04T22:14:57.003Z"}}';
const g=new APICallError(JSON.parse(body).error.message,{statusCode:429,responseBody:body,isRetryable:true});
const real=util.inspect(g).split('\n').filter(l=>!/^\s+at /.test(l));       // genuine dump minus stack frames
const forged=new Error('callback failed\n'+real.join('\n'));                  // message carries the text
const out=util.inspect(forged).split('\n').filter(l=>!/^\s+at /.test(l));
console.log('--- genuine (frames removed)');console.log(real.join('\n'));
console.log('--- forged (frames removed)');console.log(out.join('\n'));
const allIn=real.every(l=>out.includes(l));
console.log('--- every genuine line appears verbatim as a whole line in forged output:',allIn);
