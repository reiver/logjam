import{aa as Hr,d1 as Mr,d2 as Ur,d3 as Tr,d4 as Lr,d5 as Cr,d6 as Gr,cO as br,d7 as Jr,P as k,d8 as Kr}from"./index-f6fefe02.js";import{s as Vr,c as Wr}from"./callBound-45536e65.js";const Xr=Object.freeze(Object.defineProperty({__proto__:null,BigNumber:Hr,FixedFormat:Mr,FixedNumber:Ur,_base16To36:Tr,_base36To16:Lr,formatFixed:Cr,parseFixed:Gr},Symbol.toStringTag,{value:"Module"})),je=br(Jr),$e=br(Xr);var Er={exports:{}},nr=Object.prototype.toString,wr=function(e){var t=nr.call(e),o=t==="[object Arguments]";return o||(o=t!=="[object Array]"&&e!==null&&typeof e=="object"&&typeof e.length=="number"&&e.length>=0&&nr.call(e.callee)==="[object Function]"),o},X,ir;function Yr(){if(ir)return X;ir=1;var r;if(!Object.keys){var e=Object.prototype.hasOwnProperty,t=Object.prototype.toString,o=wr,f=Object.prototype.propertyIsEnumerable,u=!f.call({toString:null},"toString"),g=f.call(function(){},"prototype"),c=["toString","toLocaleString","valueOf","hasOwnProperty","isPrototypeOf","propertyIsEnumerable","constructor"],m=function(v){var y=v.constructor;return y&&y.prototype===v},E={$applicationCache:!0,$console:!0,$external:!0,$frame:!0,$frameElement:!0,$frames:!0,$innerHeight:!0,$innerWidth:!0,$onmozfullscreenchange:!0,$onmozfullscreenerror:!0,$outerHeight:!0,$outerWidth:!0,$pageXOffset:!0,$pageYOffset:!0,$parent:!0,$scrollLeft:!0,$scrollTop:!0,$scrollX:!0,$scrollY:!0,$self:!0,$webkitIndexedDB:!0,$webkitStorageInfo:!0,$window:!0},B=function(){if(typeof window>"u")return!1;for(var v in window)try{if(!E["$"+v]&&e.call(window,v)&&window[v]!==null&&typeof window[v]=="object")try{m(window[v])}catch{return!0}}catch{return!0}return!1}(),D=function(v){if(typeof window>"u"||!B)return m(v);try{return m(v)}catch{return!1}};r=function(y){var H=y!==null&&typeof y=="object",N=t.call(y)==="[object Function]",z=o(y),J=H&&t.call(y)==="[object String]",A=[];if(!H&&!N&&!z)throw new TypeError("Object.keys called on a non-object");var I=g&&N;if(J&&y.length>0&&!e.call(y,0))for(var M=0;M<y.length;++M)A.push(String(M));if(z&&y.length>0)for(var S=0;S<y.length;++S)A.push(String(S));else for(var q in y)!(I&&q==="prototype")&&e.call(y,q)&&A.push(String(q));if(u)for(var _=D(y),$=0;$<c.length;++$)!(_&&c[$]==="constructor")&&e.call(y,c[$])&&A.push(c[$]);return A}}return X=r,X}var Zr=Array.prototype.slice,Qr=wr,or=Object.keys,G=or?function(e){return or(e)}:Yr(),fr=Object.keys;G.shim=function(){if(Object.keys){var e=function(){var t=Object.keys(arguments);return t&&t.length===arguments.length}(1,2);e||(Object.keys=function(o){return Qr(o)?fr(Zr.call(o)):fr(o)})}else Object.keys=G;return Object.keys||G};var xr=G,re=xr,Or=Vr(),Sr=Wr,sr=Object,ee=Sr("Array.prototype.push"),ar=Sr("Object.prototype.propertyIsEnumerable"),te=Or?Object.getOwnPropertySymbols:null,ne=function(e,t){if(e==null)throw new TypeError("target must be an object");var o=sr(e);if(arguments.length===1)return o;for(var f=1;f<arguments.length;++f){var u=sr(arguments[f]),g=re(u),c=Or&&(Object.getOwnPropertySymbols||te);if(c)for(var m=c(u),E=0;E<m.length;++E){var B=m[E];ar(u,B)&&ee(g,B)}for(var D=0;D<g.length;++D){var v=g[D];if(ar(u,v)){var y=u[v];o[v]=y}}}return o},Y=ne,ie=function(){if(!Object.assign)return!1;for(var r="abcdefghijklmnopqrst",e=r.split(""),t={},o=0;o<e.length;++o)t[e[o]]=e[o];var f=Object.assign({},t),u="";for(var g in f)u+=g;return r!==u},oe=function(){if(!Object.assign||!Object.preventExtensions)return!1;var r=Object.preventExtensions({1:2});try{Object.assign(r,"xy")}catch{return r[1]==="y"}return!1},fe=function(){return!Object.assign||ie()||oe()?Y:Object.assign},jr={},se=function(e){return e&&typeof e=="object"&&typeof e.copy=="function"&&typeof e.fill=="function"&&typeof e.readUInt8=="function"},Z={exports:{}};typeof Object.create=="function"?Z.exports=function(e,t){e.super_=t,e.prototype=Object.create(t.prototype,{constructor:{value:e,enumerable:!1,writable:!0,configurable:!0}})}:Z.exports=function(e,t){e.super_=t;var o=function(){};o.prototype=t.prototype,e.prototype=new o,e.prototype.constructor=e};var ae=Z.exports;(function(r){var e=/%[sdj%]/g;r.format=function(n){if(!I(n)){for(var i=[],s=0;s<arguments.length;s++)i.push(f(arguments[s]));return i.join(" ")}for(var s=1,p=arguments,w=p.length,h=String(n).replace(e,function(d){if(d==="%%")return"%";if(s>=w)return d;switch(d){case"%s":return String(p[s++]);case"%d":return Number(p[s++]);case"%j":try{return JSON.stringify(p[s++])}catch{return"[Circular]"}default:return d}}),l=p[s];s<w;l=p[++s])z(l)||!_(l)?h+=" "+l:h+=" "+f(l);return h},r.deprecate=function(n,i){if(S(k.process))return function(){return r.deprecate(n,i).apply(this,arguments)};if(process.noDeprecation===!0)return n;var s=!1;function p(){if(!s){if(process.throwDeprecation)throw new Error(i);process.traceDeprecation?console.trace(i):console.error(i),s=!0}return n.apply(this,arguments)}return p};var t={},o;r.debuglog=function(n){if(S(o)&&(o={}.NODE_DEBUG||""),n=n.toUpperCase(),!t[n])if(new RegExp("\\b"+n+"\\b","i").test(o)){var i=process.pid;t[n]=function(){var s=r.format.apply(r,arguments);console.error("%s %d: %s",n,i,s)}}else t[n]=function(){};return t[n]};function f(n,i){var s={seen:[],stylize:g};return arguments.length>=3&&(s.depth=arguments[2]),arguments.length>=4&&(s.colors=arguments[3]),N(i)?s.showHidden=i:i&&r._extend(s,i),S(s.showHidden)&&(s.showHidden=!1),S(s.depth)&&(s.depth=2),S(s.colors)&&(s.colors=!1),S(s.customInspect)&&(s.customInspect=!0),s.colors&&(s.stylize=u),m(s,n,s.depth)}r.inspect=f,f.colors={bold:[1,22],italic:[3,23],underline:[4,24],inverse:[7,27],white:[37,39],grey:[90,39],black:[30,39],blue:[34,39],cyan:[36,39],green:[32,39],magenta:[35,39],red:[31,39],yellow:[33,39]},f.styles={special:"cyan",number:"yellow",boolean:"yellow",undefined:"grey",null:"bold",string:"green",date:"magenta",regexp:"red"};function u(n,i){var s=f.styles[i];return s?"\x1B["+f.colors[s][0]+"m"+n+"\x1B["+f.colors[s][1]+"m":n}function g(n,i){return n}function c(n){var i={};return n.forEach(function(s,p){i[s]=!0}),i}function m(n,i,s){if(n.customInspect&&i&&L(i.inspect)&&i.inspect!==r.inspect&&!(i.constructor&&i.constructor.prototype===i)){var p=i.inspect(s,n);return I(p)||(p=m(n,p,s)),p}var w=E(n,i);if(w)return w;var h=Object.keys(i),l=c(h);if(n.showHidden&&(h=Object.getOwnPropertyNames(i)),T(i)&&(h.indexOf("message")>=0||h.indexOf("description")>=0))return B(i);if(h.length===0){if(L(i)){var d=i.name?": "+i.name:"";return n.stylize("[Function"+d+"]","special")}if(q(i))return n.stylize(RegExp.prototype.toString.call(i),"regexp");if($(i))return n.stylize(Date.prototype.toString.call(i),"date");if(T(i))return B(i)}var b="",P=!1,C=["{","}"];if(H(i)&&(P=!0,C=["[","]"]),L(i)){var Ir=i.name?": "+i.name:"";b=" [Function"+Ir+"]"}if(q(i)&&(b=" "+RegExp.prototype.toString.call(i)),$(i)&&(b=" "+Date.prototype.toUTCString.call(i)),T(i)&&(b=" "+B(i)),h.length===0&&(!P||i.length==0))return C[0]+b+C[1];if(s<0)return q(i)?n.stylize(RegExp.prototype.toString.call(i),"regexp"):n.stylize("[Object]","special");n.seen.push(i);var W;return P?W=D(n,i,s,l,h):W=h.map(function(Rr){return v(n,i,s,l,Rr,P)}),n.seen.pop(),y(W,b,C)}function E(n,i){if(S(i))return n.stylize("undefined","undefined");if(I(i)){var s="'"+JSON.stringify(i).replace(/^"|"$/g,"").replace(/'/g,"\\'").replace(/\\"/g,'"')+"'";return n.stylize(s,"string")}if(A(i))return n.stylize(""+i,"number");if(N(i))return n.stylize(""+i,"boolean");if(z(i))return n.stylize("null","null")}function B(n){return"["+Error.prototype.toString.call(n)+"]"}function D(n,i,s,p,w){for(var h=[],l=0,d=i.length;l<d;++l)tr(i,String(l))?h.push(v(n,i,s,p,String(l),!0)):h.push("");return w.forEach(function(b){b.match(/^\d+$/)||h.push(v(n,i,s,p,b,!0))}),h}function v(n,i,s,p,w,h){var l,d,b;if(b=Object.getOwnPropertyDescriptor(i,w)||{value:i[w]},b.get?b.set?d=n.stylize("[Getter/Setter]","special"):d=n.stylize("[Getter]","special"):b.set&&(d=n.stylize("[Setter]","special")),tr(p,w)||(l="["+w+"]"),d||(n.seen.indexOf(b.value)<0?(z(s)?d=m(n,b.value,null):d=m(n,b.value,s-1),d.indexOf(`
`)>-1&&(h?d=d.split(`
`).map(function(P){return"  "+P}).join(`
`).substr(2):d=`
`+d.split(`
`).map(function(P){return"   "+P}).join(`
`))):d=n.stylize("[Circular]","special")),S(l)){if(h&&w.match(/^\d+$/))return d;l=JSON.stringify(""+w),l.match(/^"([a-zA-Z_][a-zA-Z_0-9]*)"$/)?(l=l.substr(1,l.length-2),l=n.stylize(l,"name")):(l=l.replace(/'/g,"\\'").replace(/\\"/g,'"').replace(/(^"|"$)/g,"'"),l=n.stylize(l,"string"))}return l+": "+d}function y(n,i,s){var p=n.reduce(function(w,h){return h.indexOf(`
`)>=0,w+h.replace(/\u001b\[\d\d?m/g,"").length+1},0);return p>60?s[0]+(i===""?"":i+`
 `)+" "+n.join(`,
  `)+" "+s[1]:s[0]+i+" "+n.join(", ")+" "+s[1]}function H(n){return Array.isArray(n)}r.isArray=H;function N(n){return typeof n=="boolean"}r.isBoolean=N;function z(n){return n===null}r.isNull=z;function J(n){return n==null}r.isNullOrUndefined=J;function A(n){return typeof n=="number"}r.isNumber=A;function I(n){return typeof n=="string"}r.isString=I;function M(n){return typeof n=="symbol"}r.isSymbol=M;function S(n){return n===void 0}r.isUndefined=S;function q(n){return _(n)&&K(n)==="[object RegExp]"}r.isRegExp=q;function _(n){return typeof n=="object"&&n!==null}r.isObject=_;function $(n){return _(n)&&K(n)==="[object Date]"}r.isDate=$;function T(n){return _(n)&&(K(n)==="[object Error]"||n instanceof Error)}r.isError=T;function L(n){return typeof n=="function"}r.isFunction=L;function kr(n){return n===null||typeof n=="boolean"||typeof n=="number"||typeof n=="string"||typeof n=="symbol"||typeof n>"u"}r.isPrimitive=kr,r.isBuffer=se;function K(n){return Object.prototype.toString.call(n)}function V(n){return n<10?"0"+n.toString(10):n.toString(10)}var Fr=["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"];function Nr(){var n=new Date,i=[V(n.getHours()),V(n.getMinutes()),V(n.getSeconds())].join(":");return[n.getDate(),Fr[n.getMonth()],i].join(" ")}r.log=function(){console.log("%s - %s",Nr(),r.format.apply(r,arguments))},r.inherits=ae,r._extend=function(n,i){if(!i||!_(i))return n;for(var s=Object.keys(i),p=s.length;p--;)n[s[p]]=i[s[p]];return n};function tr(n,i){return Object.prototype.hasOwnProperty.call(n,i)}})(jr);var ue=fe();/*!
 * The buffer module from node.js, for the browser.
 *
 * @author   Feross Aboukhadijeh <feross@feross.org> <http://feross.org>
 * @license  MIT
 */function ur(r,e){if(r===e)return 0;for(var t=r.length,o=e.length,f=0,u=Math.min(t,o);f<u;++f)if(r[f]!==e[f]){t=r[f],o=e[f];break}return t<o?-1:o<t?1:0}function U(r){return k.Buffer&&typeof k.Buffer.isBuffer=="function"?k.Buffer.isBuffer(r):!!(r!=null&&r._isBuffer)}var j=jr,le=Object.prototype.hasOwnProperty,lr=Array.prototype.slice,$r=function(){return(function(){}).name==="foo"}();function cr(r){return Object.prototype.toString.call(r)}function gr(r){return U(r)||typeof k.ArrayBuffer!="function"?!1:typeof ArrayBuffer.isView=="function"?ArrayBuffer.isView(r):r?!!(r instanceof DataView||r.buffer&&r.buffer instanceof ArrayBuffer):!1}var a=Er.exports=Ar,ce=/\s*function\s+([^\(\s]*)\s*/;function Br(r){if(j.isFunction(r)){if($r)return r.name;var e=r.toString(),t=e.match(ce);return t&&t[1]}}a.AssertionError=function(e){this.name="AssertionError",this.actual=e.actual,this.expected=e.expected,this.operator=e.operator,e.message?(this.message=e.message,this.generatedMessage=!1):(this.message=ge(this),this.generatedMessage=!0);var t=e.stackStartFunction||O;if(Error.captureStackTrace)Error.captureStackTrace(this,t);else{var o=new Error;if(o.stack){var f=o.stack,u=Br(t),g=f.indexOf(`
`+u);if(g>=0){var c=f.indexOf(`
`,g+1);f=f.substring(c+1)}this.stack=f}}};j.inherits(a.AssertionError,Error);function pr(r,e){return typeof r=="string"?r.length<e?r:r.slice(0,e):r}function yr(r){if($r||!j.isFunction(r))return j.inspect(r);var e=Br(r),t=e?": "+e:"";return"[Function"+t+"]"}function ge(r){return pr(yr(r.actual),128)+" "+r.operator+" "+pr(yr(r.expected),128)}function O(r,e,t,o,f){throw new a.AssertionError({message:t,actual:r,expected:e,operator:o,stackStartFunction:f})}a.fail=O;function Ar(r,e){r||O(r,!0,e,"==",a.ok)}a.ok=Ar;a.equal=function(e,t,o){e!=t&&O(e,t,o,"==",a.equal)};a.notEqual=function(e,t,o){e==t&&O(e,t,o,"!=",a.notEqual)};a.deepEqual=function(e,t,o){R(e,t,!1)||O(e,t,o,"deepEqual",a.deepEqual)};a.deepStrictEqual=function(e,t,o){R(e,t,!0)||O(e,t,o,"deepStrictEqual",a.deepStrictEqual)};function R(r,e,t,o){if(r===e)return!0;if(U(r)&&U(e))return ur(r,e)===0;if(j.isDate(r)&&j.isDate(e))return r.getTime()===e.getTime();if(j.isRegExp(r)&&j.isRegExp(e))return r.source===e.source&&r.global===e.global&&r.multiline===e.multiline&&r.lastIndex===e.lastIndex&&r.ignoreCase===e.ignoreCase;if((r===null||typeof r!="object")&&(e===null||typeof e!="object"))return t?r===e:r==e;if(gr(r)&&gr(e)&&cr(r)===cr(e)&&!(r instanceof Float32Array||r instanceof Float64Array))return ur(new Uint8Array(r.buffer),new Uint8Array(e.buffer))===0;if(U(r)!==U(e))return!1;o=o||{actual:[],expected:[]};var f=o.actual.indexOf(r);return f!==-1&&f===o.expected.indexOf(e)?!0:(o.actual.push(r),o.expected.push(e),pe(r,e,t,o))}function hr(r){return Object.prototype.toString.call(r)=="[object Arguments]"}function pe(r,e,t,o){if(r==null||e===null||e===void 0)return!1;if(j.isPrimitive(r)||j.isPrimitive(e))return r===e;if(t&&Object.getPrototypeOf(r)!==Object.getPrototypeOf(e))return!1;var f=hr(r),u=hr(e);if(f&&!u||!f&&u)return!1;if(f)return r=lr.call(r),e=lr.call(e),R(r,e,t);var g=vr(r),c=vr(e),m,E;if(g.length!==c.length)return!1;for(g.sort(),c.sort(),E=g.length-1;E>=0;E--)if(g[E]!==c[E])return!1;for(E=g.length-1;E>=0;E--)if(m=g[E],!R(r[m],e[m],t,o))return!1;return!0}a.notDeepEqual=function(e,t,o){R(e,t,!1)&&O(e,t,o,"notDeepEqual",a.notDeepEqual)};a.notDeepStrictEqual=qr;function qr(r,e,t){R(r,e,!0)&&O(r,e,t,"notDeepStrictEqual",qr)}a.strictEqual=function(e,t,o){e!==t&&O(e,t,o,"===",a.strictEqual)};a.notStrictEqual=function(e,t,o){e===t&&O(e,t,o,"!==",a.notStrictEqual)};function dr(r,e){if(!r||!e)return!1;if(Object.prototype.toString.call(e)=="[object RegExp]")return e.test(r);try{if(r instanceof e)return!0}catch{}return Error.isPrototypeOf(e)?!1:e.call({},r)===!0}function ye(r){var e;try{r()}catch(t){e=t}return e}function _r(r,e,t,o){var f;if(typeof e!="function")throw new TypeError('"block" argument must be a function');typeof t=="string"&&(o=t,t=null),f=ye(e),o=(t&&t.name?" ("+t.name+").":".")+(o?" "+o:"."),r&&!f&&O(f,t,"Missing expected exception"+o);var u=typeof o=="string",g=!r&&j.isError(f),c=!r&&f&&!t;if((g&&u&&dr(f,t)||c)&&O(f,t,"Got unwanted exception"+o),r&&f&&t&&!dr(f,t)||!r&&f)throw f}a.throws=function(r,e,t){_r(!0,r,e,t)};a.doesNotThrow=function(r,e,t){_r(!1,r,e,t)};a.ifError=function(r){if(r)throw r};function Pr(r,e){r||O(r,!0,e,"==",Pr)}a.strict=ue(Pr,a,{equal:a.strictEqual,deepEqual:a.deepStrictEqual,notEqual:a.notStrictEqual,notDeepEqual:a.notDeepStrictEqual});a.strict.strict=a.strict;var vr=Object.keys||function(r){var e=[];for(var t in r)le.call(r,t)&&e.push(t);return e},Be=Er.exports,F={},he=k&&k.__importDefault||function(r){return r&&r.__esModule?r:{default:r}};Object.defineProperty(F,"__esModule",{value:!0});F.getLength=F.decode=F.encode=void 0;var de=he(Kr);function Dr(r){if(Array.isArray(r)){for(var e=[],t=0;t<r.length;t++)e.push(Dr(r[t]));var o=Buffer.concat(e);return Buffer.concat([mr(o.length,192),o])}else{var f=er(r);return f.length===1&&f[0]<128?f:Buffer.concat([mr(f.length,128),f])}}F.encode=Dr;function Q(r,e){if(r[0]==="0"&&r[1]==="0")throw new Error("invalid RLP: extra zeros");return parseInt(r,e)}function mr(r,e){if(r<56)return Buffer.from([r+e]);var t=rr(r),o=t.length/2,f=rr(e+55+o);return Buffer.from(f+t,"hex")}function ve(r,e){if(e===void 0&&(e=!1),!r||r.length===0)return Buffer.from([]);var t=er(r),o=x(t);if(e)return o;if(o.remainder.length!==0)throw new Error("invalid remainder");return o.data}F.decode=ve;function me(r){if(!r||r.length===0)return Buffer.from([]);var e=er(r),t=e[0];if(t<=127)return e.length;if(t<=183)return t-127;if(t<=191)return t-182;if(t<=247)return t-191;var o=t-246,f=Q(e.slice(1,o).toString("hex"),16);return o+f}F.getLength=me;function x(r){var e,t,o,f,u,g=[],c=r[0];if(c<=127)return{data:r.slice(0,1),remainder:r.slice(1)};if(c<=183){if(e=c-127,c===128?o=Buffer.from([]):o=r.slice(1,e),e===2&&o[0]<128)throw new Error("invalid rlp encoding: byte must be less 0x80");return{data:o,remainder:r.slice(e)}}else if(c<=191){if(t=c-182,r.length-1<t)throw new Error("invalid RLP: not enough bytes for string length");if(e=Q(r.slice(1,t).toString("hex"),16),e<=55)throw new Error("invalid RLP: expected string length to be greater than 55");if(o=r.slice(t,e+t),o.length<e)throw new Error("invalid RLP: not enough bytes for string");return{data:o,remainder:r.slice(e+t)}}else if(c<=247){for(e=c-191,f=r.slice(1,e);f.length;)u=x(f),g.push(u.data),f=u.remainder;return{data:g,remainder:r.slice(e)}}else{t=c-246,e=Q(r.slice(1,t).toString("hex"),16);var m=t+e;if(m>r.length)throw new Error("invalid rlp: total length is larger than the data");if(f=r.slice(t,m),f.length===0)throw new Error("invalid rlp, List has a invalid length");for(;f.length;)u=x(f),g.push(u.data),f=u.remainder;return{data:g,remainder:r.slice(m)}}}function zr(r){return r.slice(0,2)==="0x"}function be(r){return typeof r!="string"?r:zr(r)?r.slice(2):r}function rr(r){if(r<0)throw new Error("Invalid integer as argument, must be unsigned!");var e=r.toString(16);return e.length%2?"0"+e:e}function Ee(r){return r.length%2?"0"+r:r}function we(r){var e=rr(r);return Buffer.from(e,"hex")}function er(r){if(!Buffer.isBuffer(r)){if(typeof r=="string")return zr(r)?Buffer.from(Ee(be(r)),"hex"):Buffer.from(r);if(typeof r=="number"||typeof r=="bigint")return r?we(r):Buffer.from([]);if(r==null)return Buffer.from([]);if(r instanceof Uint8Array)return Buffer.from(r);if(de.default.isBN(r))return Buffer.from(r.toArray());throw new Error("invalid type")}return r}export{$e as a,Be as b,F as d,je as r};
