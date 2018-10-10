/* Array.js v0.0.2 | https://github.com/iwind/Array.js */
Array.$nil={};Array.prototype.$contains=function(a){var c=this;if(c==null){return false}for(var b=0;b<c.length;b++){if(c[b]==a){return true}}return false};Array.prototype.$include=function(a){var b=this;if(b==null){return false}return b.$contains(a)};Array.prototype.$removeValue=function(b){var d=this;if(d==null){return true}var a=[];for(var c=0;c<d.length;c++){if(d[c]!=b){a.push(d[c])}}d.$clear();d.$pushAll(a);return true};Array.prototype.$remove=function(a){var b=this;if(b==null){return true}b.splice(a,1);return true};Array.prototype.$removeIf=function(b){var c=this;if(c==null){return 0}var a=c.length;var d=c.$reject(b);c.$replace(d);return a-c.length};Array.prototype.$keepIf=function(b){var c=this;if(c==null){return 0}var a=c.length;var d=c.$findAll(b);c.$replace(d);return a-c.length};Array.prototype.$replace=function(a){var b=this;if(b==null){return false}if(!Array.isArray(a)){return false}b.splice.apply(b,[0,b.length].concat(a));return true};Array.prototype.$clear=function(){var a=this;if(a==null){return true}if(a.length==0){return true}a.splice(0,a.length);return true};Array.prototype.$each=function(b){var d=this;if(d==null){return true}if(typeof(b)!="function"){return true}var c=d.length;for(var a=0;a<c;a++){b.call(d,a,d[a])}return true};Array.prototype.$unique=function(d){var e=this;if(e==null){return true}var a=[];var b=[];e.$each(function(h,g){if(typeof(d)=="function"){g=d.call(e,h,g)}if(!a.$contains(g)){a.push(g);b.push(h)}});var f=e.$copy();e.$clear();for(var c=0;c<b.length;c++){e.push(f[b[c]])}return true};Array.prototype.$get=function(a){var b=this;if(b==null){return null}if(a>b.length-1){return null}return b[a]};Array.prototype.$getAll=function(e){var d=this;if(d==null){return[]}var b=[];for(var c=0;c<arguments.length;c++){var a=arguments[c];if(Array.$isArray(a)){b.$pushAll(d.$getAll.apply(d,a))}else{if(typeof(a)=="number"&&a<d.length){b.$push(d.$get(a))}else{if(typeof(a)=="string"&&/^\\d+$/.test(a)){a=parseInt(a);if(a<d.length){b.$push(d.$get(a))}}}}}return b};Array.prototype.$set=function(a,c){var b=this;if(b==null){return false}if(a>b.length-1){return false}b[a]=c;return true};Array.prototype.$copy=function(){var c=this;if(c==null){return c}var a=[];for(var b=0;b<c.length;b++){a.push(c[b])}return a};Array.prototype.$isEmpty=function(){var a=this;if(a==null){return true}return(a.length==0)};Array.prototype.$all=function(b){var c=this;if(c==null){return false}for(var a=0;a<c.length;a++){if(!b.call(c,a,c[a])){return false}}return true};Array.prototype.$any=function(b){var c=this;if(c==null){return false}for(var a=0;a<c.length;a++){if(b.call(c,a,c[a])){return true}}return false};Array.prototype.$map=function(d){var e=this;if(e==null){return[]}var b=[];for(var c=0;c<e.length;c++){var a=d.call(e,c,e[c]);if(a===Array.$nil){continue}b.push(a)}return b};Array.prototype.$reduce=function(a){var b=this;if(b==null){return null}var c=null;b.$each(function(e,d){c=a.call(b,e,d,c)});return c};Array.prototype.$collect=function(a){var b=this;if(b==null){return[]}return b.$map(a)};Array.prototype.$find=function(c){var d=this;if(d==null){return -1}if(typeof(c)=="undefined"){return d.$get(0)}var b=-1;var a=null;d.$each(function(f,e){if(b>-1){return}if(c.call(d,f,e)){b=f;a=e}});return a};Array.prototype.$findAll=function(b){var c=this;if(c==null){return[]}if(typeof(b)=="undefined"){return c.$copy()}var a=[];c.$each(function(e,d){if(b.call(c,e,d)){a.push(d)}});return a};Array.prototype.$filter=function(a){var b=this;if(b==null){return[]}return b.$findAll(a)};Array.prototype.$reject=function(b){var c=this;if(c==null){return[]}if(typeof(b)=="undefined"){return[]}var a=[];c.$each(function(e,d){if(!b.call(c,e,d)){a.push(d)}});return a};Array.prototype.$grep=function(b){var a=this;if(a==null){return[]}return a.$findAll(function(d,c){if(c==null){return false}return b.test(c.toString())})};Array.prototype.$keys=function(e,a){var d=this;if(d==null){return[]}if(arguments.length==0){return Array.$range(0,d.length-1)}var c=[];if(typeof(a)=="undefined"){a=false}for(var b=0;b<d.length;b++){if((a&&e===d[b])||(!a&&e==d[b])){c.push(b)}}return c};Array.prototype.$indexesOf=function(c,a){var b=this;if(b==null){return[]}if(arguments.length==0){return Array.$range(0,b.length-1)}return b.$keys(c,a)};Array.prototype.$sort=function(a){var b=this;if(b==null){return false}if(typeof(a)=="undefined"){a=function(d,c){if(d>c){return 1}else{if(d==c){return 0}else{return -1}}}}b.sort(a);return true};Array.prototype.$rsort=function(a){var b=this;if(b==null){return false}this.$sort(a);b.reverse();return true};Array.prototype.$asort=function(a){var e=this;if(e==null){return[]}var c=[];for(var d=0;d<e.length;d++){c.push(d)}if(typeof(a)=="undefined"){a=function(g,f){if(g<f){return -1}if(g>f){return 1}return 0}}for(d=0;d<e.length;d++){for(var b=0;b<e.length;b++){if(b>0&&a(e[b-1],e[b])>0){e.$swap(b,b-1);c.$swap(b,b-1)}}}return c};Array.prototype.$arsort=function(a){var c=this;if(c==null){return[]}var b=c.$asort(a);c.reverse();b.reverse();return b};Array.prototype.$diff=function(c){var b=this;if(b==null){return[]}var a=[];b.$each(function(e,d){if(!c.$contains(d)){a.push(d)}});return a};Array.prototype.$intersect=function(c){var b=this;if(b==null){return[]}var a=[];b.$each(function(e,d){if(c.$contains(d)){a.push(d)}});return a};Array.prototype.$max=function(a){var c=this;if(c==null){return null}if(c.length>0){var b=c.$copy();b.$rsort(a);return b.$get(0)}return null};Array.prototype.$min=function(a){var c=this;if(c==null){return null}if(c.length>0){var b=c.$copy();b.$sort(a);return b.$get(0)}return null};Array.prototype.$swap=function(e,d){var c=this;if(c==null){return false}var b=c.$get(e);var a=c.$get(d);c.$set(e,a);c.$set(d,b);return true};Array.prototype.$sum=function(b){var c=this;if(c==null){return 0}var a=0;c.$each(function(e,d){if(typeof(b)=="function"){d=b.call(c,e,d)}if(typeof(d)=="number"){a+=d}else{if(typeof(d)=="string"){var f=parseFloat(d);if(!isNaN(f)){a+=f}}}});return a};Array.prototype.$product=function(b){var c=this;if(c==null){return 0}var a=1;c.$each(function(e,d){if(typeof(b)=="function"){d=b.call(c,e,d)}if(typeof(d)=="number"){a*=d}else{if(typeof(d)=="string"){var f=parseFloat(d);if(!isNaN(f)){a*=f}}}});return a};Array.prototype.$chunk=function(c){var d=this;if(d==null){return[]}if(typeof(c)=="undefined"){c=1}c=parseInt(c);if(isNaN(c)||c<1){return[]}var a=[];for(var b=0;b<d.length/c;b++){a.$push(d.slice(b*c,(b+1)*c))}return a};Array.prototype.$combine=function(f){var e=this;if(e==null){return[]}var b=e.$chunk(1);for(var d=0;d<arguments.length;d++){var a=arguments[d];if(Array.$isArray(a)){for(var c=0;c<e.length;c++){b[c].$push(a.$get(c))}}}return b};Array.prototype.$pad=function(d,b){var c=this;if(c==null){return false}if(typeof(b)=="undefined"){b=1}if(b<1){return false}for(var a=0;a<b;a++){c.push(d)}return true};Array.prototype.$fill=function(e,d){var c=this;if(c==null){return false}if(typeof(d)=="undefined"){d=c.length}if(d<c.length){return false}if(d==c.length){return true}var b=d-c.length;for(var a=0;a<b;a++){c.push(e)}return true};Array.prototype.$shuffle=function(){var a=this;if(a==null){return false}a.$sort(function(){return Math.random()-0.5});return true};Array.prototype.$rand=function(a){var b=this;if(b==null){return false}if(typeof(a)=="undefined"){a=1}var c=b.$copy();c.$shuffle();return c.slice(0,a)};Array.prototype.$size=function(){var a=this;if(a==null){return 0}return a.length};Array.prototype.$count=function(){var a=this;if(a==null){return 0}return a.length};Array.prototype.$first=function(){var a=this;if(a==null){return null}if(a.length==0){return null}return a.$get(0)};Array.prototype.$last=function(){var a=this;if(a==null){return null}if(a.length==0){return null}return a[a.length-1]};Array.prototype.$push=function(){var a=this;if(a==null){return 0}return Array.prototype.push.apply(a,arguments)};Array.prototype.$pushAll=function(b){var a=this;if(a==null){return 0}return Array.prototype.push.apply(a,b)};Array.prototype.$insert=function(b,e){var d=this;if(d==null){return false}var a=[];if(arguments.length==0){return false}for(var c=1;c<arguments.length;c++){a.push(arguments[c])}if(b<0){b=d.length+b+1}d.splice.apply(d,[b,0].concat(a));return true};Array.prototype.$asc=function(b){var a=this;if(a==null){return false}return a.$sort(function(d,c){if(typeof(d)=="object"&&typeof(c)=="object"){if(d[b]>c[b]){return 1}if(d[b]==c[b]){return 0}return -1}return 0})};Array.prototype.$desc=function(b){var a=this;if(a==null){return false}return a.$sort(function(d,c){if(typeof(d)=="object"&&typeof(c)=="object"){if(d[b]>c[b]){return -1}if(d[b]==c[b]){return 0}return 1}return 0})};Array.prototype.$equal=function(c){var b=this;if(b==null){return false}if(!Array.$isArray(c)){return false}if(b.length!=c.length){return false}for(var a=0;a<b.length;a++){if(b[a]!=c[a]){return false}}return true};Array.prototype.$loop=function(a){var b=this;if(b==null){return false}if(b.length==0){return false}a.call(b,0,b[0],{index:0,next:function(){this.index++;if(this.index>b.length-1){this.index=0}a.call(b,this.index,b[this.index],this);return this.index},sleep:function(c){var d=this;setTimeout(function(){d.next()},c)}});return true};Array.prototype.$asJSON=function(){return JSON.stringify(this)};Array.$range=function(e,a,c){var d=[];if(typeof(c)=="undefined"){c=1}if(e<a){for(var b=e;b<=a;b+=c){d.push(b)}}else{for(var b=e;b>=a;b-=c){d.push(b)}}return d};Array.$isArray=function(a){return Object.prototype.toString.call(a)==="[object Array]"};if(!Array.from){Array.from=(function(){var d=Object.prototype.toString;var e=function(g){return typeof g==="function"||d.call(g)==="[object Function]"};var c=function(h){var g=Number(h);if(isNaN(g)){return 0}if(g===0||!isFinite(g)){return g}return(g>0?1:-1)*Math.floor(Math.abs(g))};var b=Math.pow(2,53)-1;var a=function(h){var g=c(h);return Math.min(Math.max(g,0),b)};return function f(p){var g=this;var o=Object(p);if(p==null){throw new TypeError("Array.from requires an array-like object - not null or undefined")}var m=arguments.length>1?arguments[1]:void undefined;var i;if(typeof m!=="undefined"){if(!e(m)){throw new TypeError("Array.from: when provided, the second argument must be a function")}if(arguments.length>2){i=arguments[2]}}var n=a(o.length);var h=e(g)?Object(new g(n)):new Array(n);var j=0;var l;while(j<n){l=o[j];if(m){h[j]=typeof i==="undefined"?m(l,j):m.call(i,l,j)}else{h[j]=l}j+=1}h.length=n;return h}}())};

/* axios v0.18.0 | (c) 2018 by Matt Zabriskie */
!function(e,t){"object"==typeof exports&&"object"==typeof module?module.exports=t():"function"==typeof define&&define.amd?define([],t):"object"==typeof exports?exports.axios=t():e.axios=t()}(this,function(){return function(e){function t(r){if(n[r])return n[r].exports;var o=n[r]={exports:{},id:r,loaded:!1};return e[r].call(o.exports,o,o.exports,t),o.loaded=!0,o.exports}var n={};return t.m=e,t.c=n,t.p="",t(0)}([function(e,t,n){e.exports=n(1)},function(e,t,n){"use strict";function r(e){var t=new s(e),n=i(s.prototype.request,t);return o.extend(n,s.prototype,t),o.extend(n,t),n}var o=n(2),i=n(3),s=n(5),u=n(6),a=r(u);a.Axios=s,a.create=function(e){return r(o.merge(u,e))},a.Cancel=n(23),a.CancelToken=n(24),a.isCancel=n(20),a.all=function(e){return Promise.all(e)},a.spread=n(25),e.exports=a,e.exports.default=a},function(e,t,n){"use strict";function r(e){return"[object Array]"===R.call(e)}function o(e){return"[object ArrayBuffer]"===R.call(e)}function i(e){return"undefined"!=typeof FormData&&e instanceof FormData}function s(e){var t;return t="undefined"!=typeof ArrayBuffer&&ArrayBuffer.isView?ArrayBuffer.isView(e):e&&e.buffer&&e.buffer instanceof ArrayBuffer}function u(e){return"string"==typeof e}function a(e){return"number"==typeof e}function c(e){return"undefined"==typeof e}function f(e){return null!==e&&"object"==typeof e}function p(e){return"[object Date]"===R.call(e)}function d(e){return"[object File]"===R.call(e)}function l(e){return"[object Blob]"===R.call(e)}function h(e){return"[object Function]"===R.call(e)}function m(e){return f(e)&&h(e.pipe)}function y(e){return"undefined"!=typeof URLSearchParams&&e instanceof URLSearchParams}function w(e){return e.replace(/^\s*/,"").replace(/\s*$/,"")}function g(){return("undefined"==typeof navigator||"ReactNative"!==navigator.product)&&("undefined"!=typeof window&&"undefined"!=typeof document)}function v(e,t){if(null!==e&&"undefined"!=typeof e)if("object"!=typeof e&&(e=[e]),r(e))for(var n=0,o=e.length;n<o;n++)t.call(null,e[n],n,e);else for(var i in e)Object.prototype.hasOwnProperty.call(e,i)&&t.call(null,e[i],i,e)}function x(){function e(e,n){"object"==typeof t[n]&&"object"==typeof e?t[n]=x(t[n],e):t[n]=e}for(var t={},n=0,r=arguments.length;n<r;n++)v(arguments[n],e);return t}function b(e,t,n){return v(t,function(t,r){n&&"function"==typeof t?e[r]=E(t,n):e[r]=t}),e}var E=n(3),C=n(4),R=Object.prototype.toString;e.exports={isArray:r,isArrayBuffer:o,isBuffer:C,isFormData:i,isArrayBufferView:s,isString:u,isNumber:a,isObject:f,isUndefined:c,isDate:p,isFile:d,isBlob:l,isFunction:h,isStream:m,isURLSearchParams:y,isStandardBrowserEnv:g,forEach:v,merge:x,extend:b,trim:w}},function(e,t){"use strict";e.exports=function(e,t){return function(){for(var n=new Array(arguments.length),r=0;r<n.length;r++)n[r]=arguments[r];return e.apply(t,n)}}},function(e,t){function n(e){return!!e.constructor&&"function"==typeof e.constructor.isBuffer&&e.constructor.isBuffer(e)}function r(e){return"function"==typeof e.readFloatLE&&"function"==typeof e.slice&&n(e.slice(0,0))}/*!
	 * Determine if an object is a Buffer
	 *
	 * @author   Feross Aboukhadijeh <https://feross.org>
	 * @license  MIT
	 */
    e.exports=function(e){return null!=e&&(n(e)||r(e)||!!e._isBuffer)}},function(e,t,n){"use strict";function r(e){this.defaults=e,this.interceptors={request:new s,response:new s}}var o=n(6),i=n(2),s=n(17),u=n(18);r.prototype.request=function(e){"string"==typeof e&&(e=i.merge({url:arguments[0]},arguments[1])),e=i.merge(o,{method:"get"},this.defaults,e),e.method=e.method.toLowerCase();var t=[u,void 0],n=Promise.resolve(e);for(this.interceptors.request.forEach(function(e){t.unshift(e.fulfilled,e.rejected)}),this.interceptors.response.forEach(function(e){t.push(e.fulfilled,e.rejected)});t.length;)n=n.then(t.shift(),t.shift());return n},i.forEach(["delete","get","head","options"],function(e){r.prototype[e]=function(t,n){return this.request(i.merge(n||{},{method:e,url:t}))}}),i.forEach(["post","put","patch"],function(e){r.prototype[e]=function(t,n,r){return this.request(i.merge(r||{},{method:e,url:t,data:n}))}}),e.exports=r},function(e,t,n){"use strict";function r(e,t){!i.isUndefined(e)&&i.isUndefined(e["Content-Type"])&&(e["Content-Type"]=t)}function o(){var e;return"undefined"!=typeof XMLHttpRequest?e=n(8):"undefined"!=typeof process&&(e=n(8)),e}var i=n(2),s=n(7),u={"Content-Type":"application/x-www-form-urlencoded"},a={adapter:o(),transformRequest:[function(e,t){return s(t,"Content-Type"),i.isFormData(e)||i.isArrayBuffer(e)||i.isBuffer(e)||i.isStream(e)||i.isFile(e)||i.isBlob(e)?e:i.isArrayBufferView(e)?e.buffer:i.isURLSearchParams(e)?(r(t,"application/x-www-form-urlencoded;charset=utf-8"),e.toString()):i.isObject(e)?(r(t,"application/json;charset=utf-8"),JSON.stringify(e)):e}],transformResponse:[function(e){if("string"==typeof e)try{e=JSON.parse(e)}catch(e){}return e}],timeout:0,xsrfCookieName:"XSRF-TOKEN",xsrfHeaderName:"X-XSRF-TOKEN",maxContentLength:-1,validateStatus:function(e){return e>=200&&e<300}};a.headers={common:{Accept:"application/json, text/plain, */*"}},i.forEach(["delete","get","head"],function(e){a.headers[e]={}}),i.forEach(["post","put","patch"],function(e){a.headers[e]=i.merge(u)}),e.exports=a},function(e,t,n){"use strict";var r=n(2);e.exports=function(e,t){r.forEach(e,function(n,r){r!==t&&r.toUpperCase()===t.toUpperCase()&&(e[t]=n,delete e[r])})}},function(e,t,n){"use strict";var r=n(2),o=n(9),i=n(12),s=n(13),u=n(14),a=n(10),c="undefined"!=typeof window&&window.btoa&&window.btoa.bind(window)||n(15);e.exports=function(e){return new Promise(function(t,f){var p=e.data,d=e.headers;r.isFormData(p)&&delete d["Content-Type"];var l=new XMLHttpRequest,h="onreadystatechange",m=!1;if("undefined"==typeof window||!window.XDomainRequest||"withCredentials"in l||u(e.url)||(l=new window.XDomainRequest,h="onload",m=!0,l.onprogress=function(){},l.ontimeout=function(){}),e.auth){var y=e.auth.username||"",w=e.auth.password||"";d.Authorization="Basic "+c(y+":"+w)}if(l.open(e.method.toUpperCase(),i(e.url,e.params,e.paramsSerializer),!0),l.timeout=e.timeout,l[h]=function(){if(l&&(4===l.readyState||m)&&(0!==l.status||l.responseURL&&0===l.responseURL.indexOf("file:"))){var n="getAllResponseHeaders"in l?s(l.getAllResponseHeaders()):null,r=e.responseType&&"text"!==e.responseType?l.response:l.responseText,i={data:r,status:1223===l.status?204:l.status,statusText:1223===l.status?"No Content":l.statusText,headers:n,config:e,request:l};o(t,f,i),l=null}},l.onerror=function(){f(a("Network Error",e,null,l)),l=null},l.ontimeout=function(){f(a("timeout of "+e.timeout+"ms exceeded",e,"ECONNABORTED",l)),l=null},r.isStandardBrowserEnv()){var g=n(16),v=(e.withCredentials||u(e.url))&&e.xsrfCookieName?g.read(e.xsrfCookieName):void 0;v&&(d[e.xsrfHeaderName]=v)}if("setRequestHeader"in l&&r.forEach(d,function(e,t){"undefined"==typeof p&&"content-type"===t.toLowerCase()?delete d[t]:l.setRequestHeader(t,e)}),e.withCredentials&&(l.withCredentials=!0),e.responseType)try{l.responseType=e.responseType}catch(t){if("json"!==e.responseType)throw t}"function"==typeof e.onDownloadProgress&&l.addEventListener("progress",e.onDownloadProgress),"function"==typeof e.onUploadProgress&&l.upload&&l.upload.addEventListener("progress",e.onUploadProgress),e.cancelToken&&e.cancelToken.promise.then(function(e){l&&(l.abort(),f(e),l=null)}),void 0===p&&(p=null),l.send(p)})}},function(e,t,n){"use strict";var r=n(10);e.exports=function(e,t,n){var o=n.config.validateStatus;n.status&&o&&!o(n.status)?t(r("Request failed with status code "+n.status,n.config,null,n.request,n)):e(n)}},function(e,t,n){"use strict";var r=n(11);e.exports=function(e,t,n,o,i){var s=new Error(e);return r(s,t,n,o,i)}},function(e,t){"use strict";e.exports=function(e,t,n,r,o){return e.config=t,n&&(e.code=n),e.request=r,e.response=o,e}},function(e,t,n){"use strict";function r(e){return encodeURIComponent(e).replace(/%40/gi,"@").replace(/%3A/gi,":").replace(/%24/g,"$").replace(/%2C/gi,",").replace(/%20/g,"+").replace(/%5B/gi,"[").replace(/%5D/gi,"]")}var o=n(2);e.exports=function(e,t,n){if(!t)return e;var i;if(n)i=n(t);else if(o.isURLSearchParams(t))i=t.toString();else{var s=[];o.forEach(t,function(e,t){null!==e&&"undefined"!=typeof e&&(o.isArray(e)?t+="[]":e=[e],o.forEach(e,function(e){o.isDate(e)?e=e.toISOString():o.isObject(e)&&(e=JSON.stringify(e)),s.push(r(t)+"="+r(e))}))}),i=s.join("&")}return i&&(e+=(e.indexOf("?")===-1?"?":"&")+i),e}},function(e,t,n){"use strict";var r=n(2),o=["age","authorization","content-length","content-type","etag","expires","from","host","if-modified-since","if-unmodified-since","last-modified","location","max-forwards","proxy-authorization","referer","retry-after","user-agent"];e.exports=function(e){var t,n,i,s={};return e?(r.forEach(e.split("\n"),function(e){if(i=e.indexOf(":"),t=r.trim(e.substr(0,i)).toLowerCase(),n=r.trim(e.substr(i+1)),t){if(s[t]&&o.indexOf(t)>=0)return;"set-cookie"===t?s[t]=(s[t]?s[t]:[]).concat([n]):s[t]=s[t]?s[t]+", "+n:n}}),s):s}},function(e,t,n){"use strict";var r=n(2);e.exports=r.isStandardBrowserEnv()?function(){function e(e){var t=e;return n&&(o.setAttribute("href",t),t=o.href),o.setAttribute("href",t),{href:o.href,protocol:o.protocol?o.protocol.replace(/:$/,""):"",host:o.host,search:o.search?o.search.replace(/^\?/,""):"",hash:o.hash?o.hash.replace(/^#/,""):"",hostname:o.hostname,port:o.port,pathname:"/"===o.pathname.charAt(0)?o.pathname:"/"+o.pathname}}var t,n=/(msie|trident)/i.test(navigator.userAgent),o=document.createElement("a");return t=e(window.location.href),function(n){var o=r.isString(n)?e(n):n;return o.protocol===t.protocol&&o.host===t.host}}():function(){return function(){return!0}}()},function(e,t){"use strict";function n(){this.message="String contains an invalid character"}function r(e){for(var t,r,i=String(e),s="",u=0,a=o;i.charAt(0|u)||(a="=",u%1);s+=a.charAt(63&t>>8-u%1*8)){if(r=i.charCodeAt(u+=.75),r>255)throw new n;t=t<<8|r}return s}var o="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=";n.prototype=new Error,n.prototype.code=5,n.prototype.name="InvalidCharacterError",e.exports=r},function(e,t,n){"use strict";var r=n(2);e.exports=r.isStandardBrowserEnv()?function(){return{write:function(e,t,n,o,i,s){var u=[];u.push(e+"="+encodeURIComponent(t)),r.isNumber(n)&&u.push("expires="+new Date(n).toGMTString()),r.isString(o)&&u.push("path="+o),r.isString(i)&&u.push("domain="+i),s===!0&&u.push("secure"),document.cookie=u.join("; ")},read:function(e){var t=document.cookie.match(new RegExp("(^|;\\s*)("+e+")=([^;]*)"));return t?decodeURIComponent(t[3]):null},remove:function(e){this.write(e,"",Date.now()-864e5)}}}():function(){return{write:function(){},read:function(){return null},remove:function(){}}}()},function(e,t,n){"use strict";function r(){this.handlers=[]}var o=n(2);r.prototype.use=function(e,t){return this.handlers.push({fulfilled:e,rejected:t}),this.handlers.length-1},r.prototype.eject=function(e){this.handlers[e]&&(this.handlers[e]=null)},r.prototype.forEach=function(e){o.forEach(this.handlers,function(t){null!==t&&e(t)})},e.exports=r},function(e,t,n){"use strict";function r(e){e.cancelToken&&e.cancelToken.throwIfRequested()}var o=n(2),i=n(19),s=n(20),u=n(6),a=n(21),c=n(22);e.exports=function(e){r(e),e.baseURL&&!a(e.url)&&(e.url=c(e.baseURL,e.url)),e.headers=e.headers||{},e.data=i(e.data,e.headers,e.transformRequest),e.headers=o.merge(e.headers.common||{},e.headers[e.method]||{},e.headers||{}),o.forEach(["delete","get","head","post","put","patch","common"],function(t){delete e.headers[t]});var t=e.adapter||u.adapter;return t(e).then(function(t){return r(e),t.data=i(t.data,t.headers,e.transformResponse),t},function(t){return s(t)||(r(e),t&&t.response&&(t.response.data=i(t.response.data,t.response.headers,e.transformResponse))),Promise.reject(t)})}},function(e,t,n){"use strict";var r=n(2);e.exports=function(e,t,n){return r.forEach(n,function(n){e=n(e,t)}),e}},function(e,t){"use strict";e.exports=function(e){return!(!e||!e.__CANCEL__)}},function(e,t){"use strict";e.exports=function(e){return/^([a-z][a-z\d\+\-\.]*:)?\/\//i.test(e)}},function(e,t){"use strict";e.exports=function(e,t){return t?e.replace(/\/+$/,"")+"/"+t.replace(/^\/+/,""):e}},function(e,t){"use strict";function n(e){this.message=e}n.prototype.toString=function(){return"Cancel"+(this.message?": "+this.message:"")},n.prototype.__CANCEL__=!0,e.exports=n},function(e,t,n){"use strict";function r(e){if("function"!=typeof e)throw new TypeError("executor must be a function.");var t;this.promise=new Promise(function(e){t=e});var n=this;e(function(e){n.reason||(n.reason=new o(e),t(n.reason))})}var o=n(23);r.prototype.throwIfRequested=function(){if(this.reason)throw this.reason},r.source=function(){var e,t=new r(function(t){e=t});return{token:t,cancel:e}},e.exports=r},function(e,t){"use strict";e.exports=function(e){return function(t){return e.apply(null,t)}}}])});
//# sourceMappingURL=axios.min.map

!function(t,e){"object"==typeof exports&&"undefined"!=typeof module?module.exports=e():"function"==typeof define&&define.amd?define(e):t.ES6Promise=e()}(this,function(){"use strict";function t(t){var e=typeof t;return null!==t&&("object"===e||"function"===e)}function e(t){return"function"==typeof t}function n(t){B=t}function r(t){G=t}function o(){return function(){return process.nextTick(a)}}function i(){return"undefined"!=typeof z?function(){z(a)}:c()}function s(){var t=0,e=new J(a),n=document.createTextNode("");return e.observe(n,{characterData:!0}),function(){n.data=t=++t%2}}function u(){var t=new MessageChannel;return t.port1.onmessage=a,function(){return t.port2.postMessage(0)}}function c(){var t=setTimeout;return function(){return t(a,1)}}function a(){for(var t=0;t<W;t+=2){var e=V[t],n=V[t+1];e(n),V[t]=void 0,V[t+1]=void 0}W=0}function f(){try{var t=Function("return this")().require("vertx");return z=t.runOnLoop||t.runOnContext,i()}catch(e){return c()}}function l(t,e){var n=this,r=new this.constructor(p);void 0===r[Z]&&O(r);var o=n._state;if(o){var i=arguments[o-1];G(function(){return P(o,r,i,n._result)})}else E(n,r,t,e);return r}function h(t){var e=this;if(t&&"object"==typeof t&&t.constructor===e)return t;var n=new e(p);return g(n,t),n}function p(){}function v(){return new TypeError("You cannot resolve a promise with itself")}function d(){return new TypeError("A promises callback cannot return that same promise.")}function _(t){try{return t.then}catch(e){return nt.error=e,nt}}function y(t,e,n,r){try{t.call(e,n,r)}catch(o){return o}}function m(t,e,n){G(function(t){var r=!1,o=y(n,e,function(n){r||(r=!0,e!==n?g(t,n):S(t,n))},function(e){r||(r=!0,j(t,e))},"Settle: "+(t._label||" unknown promise"));!r&&o&&(r=!0,j(t,o))},t)}function b(t,e){e._state===tt?S(t,e._result):e._state===et?j(t,e._result):E(e,void 0,function(e){return g(t,e)},function(e){return j(t,e)})}function w(t,n,r){n.constructor===t.constructor&&r===l&&n.constructor.resolve===h?b(t,n):r===nt?(j(t,nt.error),nt.error=null):void 0===r?S(t,n):e(r)?m(t,n,r):S(t,n)}function g(e,n){e===n?j(e,v()):t(n)?w(e,n,_(n)):S(e,n)}function A(t){t._onerror&&t._onerror(t._result),T(t)}function S(t,e){t._state===$&&(t._result=e,t._state=tt,0!==t._subscribers.length&&G(T,t))}function j(t,e){t._state===$&&(t._state=et,t._result=e,G(A,t))}function E(t,e,n,r){var o=t._subscribers,i=o.length;t._onerror=null,o[i]=e,o[i+tt]=n,o[i+et]=r,0===i&&t._state&&G(T,t)}function T(t){var e=t._subscribers,n=t._state;if(0!==e.length){for(var r=void 0,o=void 0,i=t._result,s=0;s<e.length;s+=3)r=e[s],o=e[s+n],r?P(n,r,o,i):o(i);t._subscribers.length=0}}function M(t,e){try{return t(e)}catch(n){return nt.error=n,nt}}function P(t,n,r,o){var i=e(r),s=void 0,u=void 0,c=void 0,a=void 0;if(i){if(s=M(r,o),s===nt?(a=!0,u=s.error,s.error=null):c=!0,n===s)return void j(n,d())}else s=o,c=!0;n._state!==$||(i&&c?g(n,s):a?j(n,u):t===tt?S(n,s):t===et&&j(n,s))}function x(t,e){try{e(function(e){g(t,e)},function(e){j(t,e)})}catch(n){j(t,n)}}function C(){return rt++}function O(t){t[Z]=rt++,t._state=void 0,t._result=void 0,t._subscribers=[]}function k(){return new Error("Array Methods must be provided an Array")}function F(t){return new ot(this,t).promise}function Y(t){var e=this;return new e(U(t)?function(n,r){for(var o=t.length,i=0;i<o;i++)e.resolve(t[i]).then(n,r)}:function(t,e){return e(new TypeError("You must pass an array to race."))})}function q(t){var e=this,n=new e(p);return j(n,t),n}function D(){throw new TypeError("You must pass a resolver function as the first argument to the promise constructor")}function K(){throw new TypeError("Failed to construct 'Promise': Please use the 'new' operator, this object constructor cannot be called as a function.")}function L(){var t=void 0;if("undefined"!=typeof global)t=global;else if("undefined"!=typeof self)t=self;else try{t=Function("return this")()}catch(e){throw new Error("polyfill failed because global object is unavailable in this environment")}var n=t.Promise;if(n){var r=null;try{r=Object.prototype.toString.call(n.resolve())}catch(e){}if("[object Promise]"===r&&!n.cast)return}t.Promise=it}var N=void 0;N=Array.isArray?Array.isArray:function(t){return"[object Array]"===Object.prototype.toString.call(t)};var U=N,W=0,z=void 0,B=void 0,G=function(t,e){V[W]=t,V[W+1]=e,W+=2,2===W&&(B?B(a):X())},H="undefined"!=typeof window?window:void 0,I=H||{},J=I.MutationObserver||I.WebKitMutationObserver,Q="undefined"==typeof self&&"undefined"!=typeof process&&"[object process]"==={}.toString.call(process),R="undefined"!=typeof Uint8ClampedArray&&"undefined"!=typeof importScripts&&"undefined"!=typeof MessageChannel,V=new Array(1e3),X=void 0;X=Q?o():J?s():R?u():void 0===H&&"function"==typeof require?f():c();var Z=Math.random().toString(36).substring(2),$=void 0,tt=1,et=2,nt={error:null},rt=0,ot=function(){function t(t,e){this._instanceConstructor=t,this.promise=new t(p),this.promise[Z]||O(this.promise),U(e)?(this.length=e.length,this._remaining=e.length,this._result=new Array(this.length),0===this.length?S(this.promise,this._result):(this.length=this.length||0,this._enumerate(e),0===this._remaining&&S(this.promise,this._result))):j(this.promise,k())}return t.prototype._enumerate=function(t){for(var e=0;this._state===$&&e<t.length;e++)this._eachEntry(t[e],e)},t.prototype._eachEntry=function(t,e){var n=this._instanceConstructor,r=n.resolve;if(r===h){var o=_(t);if(o===l&&t._state!==$)this._settledAt(t._state,e,t._result);else if("function"!=typeof o)this._remaining--,this._result[e]=t;else if(n===it){var i=new n(p);w(i,t,o),this._willSettleAt(i,e)}else this._willSettleAt(new n(function(e){return e(t)}),e)}else this._willSettleAt(r(t),e)},t.prototype._settledAt=function(t,e,n){var r=this.promise;r._state===$&&(this._remaining--,t===et?j(r,n):this._result[e]=n),0===this._remaining&&S(r,this._result)},t.prototype._willSettleAt=function(t,e){var n=this;E(t,void 0,function(t){return n._settledAt(tt,e,t)},function(t){return n._settledAt(et,e,t)})},t}(),it=function(){function t(e){this[Z]=C(),this._result=this._state=void 0,this._subscribers=[],p!==e&&("function"!=typeof e&&D(),this instanceof t?x(this,e):K())}return t.prototype["catch"]=function(t){return this.then(null,t)},t.prototype["finally"]=function(t){var e=this,n=e.constructor;return e.then(function(e){return n.resolve(t()).then(function(){return e})},function(e){return n.resolve(t()).then(function(){throw e})})},t}();return it.prototype.then=l,it.all=F,it.race=Y,it.resolve=h,it.reject=q,it._setScheduler=n,it._setAsap=r,it._asap=G,it.polyfill=L,it.Promise=it,it.polyfill(),it});

/** vue.tea.js **/
(function () {
    var contextFunctions = [];
    var that = this;
    var data = {};

    window.Tea = {};
    window.Tea.context = function (fn) {
        if (typeof(fn) !== "function") {
            throw new Error("Tea.scope(fn) should accept a function argument");
        }

        // 合并context
        contextFunctions.push(fn);
    };

    Vue.config.errorHandler = function (error, vue) {
        var match = error.toString().match(/(\w+) is not defined/);
        if (match != null && match.length === 2) {
            vue[match[1]] = "";
            vue.$set(vue, match[1], "");
            vue.$forceUpdate();

            console.error(error);
            return;
        }

        throw new Error(error.toString());
    };

    this.load = function () {
        if (typeof(window.TEA.ACTION.data) !== "undefined") {
            data = window.TEA.ACTION.data;
        }

        var innerMethods = {
            $delay: Tea.delay,
            $get: function (action) {
                return Tea.action(action).get();
            },
            $post: function (action) {
                return Tea.action(action).post();
            },
            $go: Tea.go,
            $url: Tea.url,
            $find: Tea.element
        };

        var vueElement = document.getElementById("tea-app");
        if (vueElement == null && document.body) {
            var rootNodes = document.body.childNodes;
            for (var i = 0; i < rootNodes.length; i ++) {
                var rootNode = rootNodes[i];
                if (rootNode.nodeType == 1) {
                    vueElement = rootNode;
                    vueElement.setAttribute("tea-root", "generated");
                    break;
                }
            }
        }

        if (contextFunctions.length > 0) {
            var context = {};
            context.Tea = window.Tea;

            // 内置方法
            for (var methodName in innerMethods) {
                if (innerMethods.hasOwnProperty(methodName)) {
                    context[methodName] = innerMethods[methodName];
                }
            }

            for (key in data) {
                if (!data.hasOwnProperty(key)) {
                    continue;
                }
                context[key] = data[key];
            }
            
            for (var i = 0; i < contextFunctions.length; i ++) {
                var contextFn = contextFunctions[i];
                if (typeof(contextFn) != "function") {
                    continue;
                }
                contextFn.call(context);
                for (var key in context) {
                    if (!context.hasOwnProperty(key)) {
                        continue;
                    }
                    if (typeof(key) != "string") {
                        continue;
                    }

                    // 跳过自定义方法
                    if (key.length > 0 && key[0] == "$") {
                        continue;
                    }

                    var value = context[key];
                    if (typeof(value) === "function") {
                        context[key] = function (value) {
                            return function () {
                                if (window.Tea.Vue == null) {
                                    return value.apply(innerMethods, arguments);
                                }
                                else {
                                    return value.apply(window.Tea.Vue, arguments);
                                }
                            };
                        }(value);
                    }
                }
            }

            // 清除context中的预定义变量
            for (var methodName in innerMethods) {
                if (innerMethods.hasOwnProperty(methodName)) {
                    delete(context[methodName]);
                }
            }

            window.Tea.Vue = new Vue({
                el: vueElement,
                data: context,

                // 内置方法
                methods: innerMethods
            });
        }
        else {
            var context = {
                Tea: window.Tea
            };
            for (key in data) {
                if (!data.hasOwnProperty(key)) {
                    continue;
                }
                context[key] = data[key];
            }

            window.Tea.Vue = new Vue({
                el: vueElement,
                data: context,

                // 内置方法
                methods: {
                    $delay: Tea.delay,
                    $get: function (action) {
                        return Tea.action(action).get();
                    },
                    $post: function (action) {
                        return Tea.action(action).post();
                    },
                    $go: Tea.go,
                    $url: Tea.url,
                    $find: Tea.element
                }
            });
        }
    };

    document.addEventListener("DOMContentLoaded", function () {
        that.load();

        if (document.body) {
            Tea.activate(document.body);
        }
    });
})();

/**
 * 序列化参数为可传递的字符串
 *
 * 代码来自jQuery：https://jquery.com/download/
 *
 * @param a 要序列化的参数集
 * @param traditional
 * @returns {*}
 */
window.Tea.serialize = function (a, traditional) {
    var prefix,
        s = [],
        add = function (key, valueOrFunction) {

            // If value is a function, invoke it and use its return value
            var value = (typeof(valueOrFunction) === "function") ?
                valueOrFunction() :
                valueOrFunction;

            s[s.length] = encodeURIComponent(key) + "=" +
                encodeURIComponent(value == null ? "" : value);
        };

    var
        rbracket = /\[]$/;

    var buildParams = function (prefix, obj, traditional, add) {
        var name;
        if (Array.isArray(obj)) {
            // Serialize array item.
            for (var i in obj) {
                if (!obj.hasOwnProperty(i)) {
                    continue;
                }
                var v = obj[i];
                if (traditional || rbracket.test(prefix)) {

                    // Treat each array item as a scalar.
                    add(prefix, v);

                } else {

                    // Item is non-scalar (array or object), encode its numeric index.
                    buildParams(
                        prefix + "[" + (typeof v === "object" && v != null ? i : "") + "]",
                        v,
                        traditional,
                        add
                    );
                }
            }

        } else if (!traditional && typeof(obj) === "object") {

            // Serialize object item.
            for (name in obj) {
                buildParams(prefix + "[" + name + "]", obj[name], traditional, add);
            }

        } else {

            // Serialize scalar item.
            add(prefix, obj);
        }
    };

    // If an array was passed in, assume that it is an array of form elements.
    if (Array.isArray(a)) {
        // Serialize the form elements
        for (key in a) {
            if (!a.hasOwnProperty(key)) {
                continue;
            }
            add(key, a[key]);
        }

    } else {

        // If traditional, encode the "old" way (the way 1.3.2 or older
        // did it), otherwise encode params recursively.
        for (prefix in a) {
            if (!a.hasOwnProperty(prefix)) {
                continue;
            }
            buildParams(prefix, a[prefix], traditional, add);
        }
    }

    // Return the resulting serialization
    return s.join("&");
};

/**
 * 取得Action对应的URL
 *
 * @param action Action
 * @param params 参数
 * @param hashParams 锚点参数
 * @returns {*}
 */
window.Tea.url = function (action, params, hashParams) {
    var config = window.TEA.ACTION;
    var controller = config.parent;
    var module = config.module;
    var base = config.base;
    var actionParam = config.actionParam;

    var url;
    if (action.match(/\//)) {//支持URL
        url = action;

        if (typeof(params) === "object") {
            var query = Tea.serialize(params);
            if (query.length > 0) {
                url += "?" + query;
            }
        }
        if (!url.match(/^(http|https|ftp):/i)) {
            url = base + ((url.substr(0, 1) === "/") ? "" : "/") + url;
        }
    }
    else {
        if (action.substr(0, 2) === "..") {
            var pos = controller.lastIndexOf(".");
            if (pos === -1) {
                action = action.substr(2);
            }
            else {
                action = controller.substr(0, pos) + action.substr(1);
            }
            if (module !== "") {
                action = "@" + module + "." + action;
            }
        }
        else if (action.substr(0, 1) === ".") {
            action = controller + action;
            if (module !== "") {
                action = "@" + module + "." + action;
            }
        }
        else if (module !== "") {
            if (action === "@") {
                action = "@" + module;
            }
            else {
                action = action.replace("@.", "@" + module + ".");
            }
        }
        action = action.replace(/\.$/, "");
        if (actionParam) {
            var path = action.replace(/[.\/]+/g, "/");
            if (path.substr(0, 1) !== "/") {
                path = "/" + path;
            }
            url = base + "?__ACTION__=" + path;
        }
        else {
            url = base + "/" + action.replace(/[.\/]+/g, "/").replace(/^\//, "");
        }
        if (typeof(params) === "object") {
            params = Tea.serialize(params);
            if (params.length > 0) {
                if (url.indexOf("?") === -1) {
                    url += "?" + params;
                }
                else {
                    url += "&" + params;
                }
            }
        }
        if (typeof(hashParams) === "string") {
            url += "#" + hashParams;
        }
        else if (typeof(hashParams) === "object") {
            url += "#" + Tea.serialize(hashParams);
        }
    }
    return url;
};

/**
 * 跳转
 *
 * @param action 要跳转到的action
 * @param params 附带的参数
 * @param hash 附带的锚点参数
 */
window.Tea.go = function (action, params, hash) {
    var url = Tea.url(action, params);
    if (hash && hash.length > 0) {
        url += "#" + hash;
    }

    window.location.href = url;
};

/**
 * 格式化字节数
 *
 * @param bytes 字节数
 * @returns {*}
 */
window.Tea.formatBytes = function (bytes) {
    if (bytes < 1024) {
        return "< 1kb";
    }
    else if (bytes < 1024 * 1024) {
        return Math.round(bytes / 1024 * 100) / 100 + " kb";
    }
    else if (bytes < 1024 * 1024 * 1024) {
        return Math.round(bytes / 1024 / 1024 * 100) / 100 + " mb";
    }
    return Math.round(bytes / 1024 / 1024 / 1024 * 100) / 100 + " gb";
};

/**
 * 版本号对比
 *
 * 代码来自 http://stackoverflow.com/questions/6832596/how-to-compare-software-version-number-using-js-only-number
 *
 * @param a
 * @param b
 * @returns {number}
 */
window.Tea.versionCompare = function compare(a, b) {
    if (a === b) {
        return 0;
    }

    var a_components = a.split(".");
    var b_components = b.split(".");

    var len = Math.min(a_components.length, b_components.length);

    // loop while the components are equal
    for (var i = 0; i < len; i++) {
        // A bigger than B
        if (parseInt(a_components[i]) > parseInt(b_components[i])) {
            return 1;
        }

        // B bigger than A
        if (parseInt(a_components[i]) < parseInt(b_components[i])) {
            return -1;
        }
    }

    // If one's a prefix of the other, the longer one is greater.
    if (a_components.length > b_components.length) {
        return 1;
    }

    if (a_components.length < b_components.length) {
        return -1;
    }

    // Otherwise they are the same.
    return 0;
};


/**
 * 延时执行
 *
 * @param fn 要执行的函数
 * @param ms 延时长度
 */
window.Tea.delay = function (fn, ms) {
    if (typeof(ms) === "undefined") {
        ms = 10;
    }
    setTimeout(function () {
        fn.call(Tea.Vue);
    }, ms);
};

/**
 * 定义Action对象
 *
 * @param action Action
 * @param params 参数集
 * @constructor
 */
window.Tea.Action = function (action, params) {
    var _action = action;
    var _params = params;
    var _successFn;
    var _failFn;
    var _errorFn;
    var _doneFn;
    var _method = "POST";
    var _timeout = 30;
    var _delay = 0;
    var _progressFn;

    this.params = function (params) {
        _params = params;
        return this;
    };

    this.form = function (form) {
        _params = new FormData(form);
        return this;
    };

    this.success = function (successFn) {
        _successFn = successFn;
        return this;
    };

    this.fail = function (failFn) {
        _failFn = failFn;
        return this;
    };

    this.error = function (errorFn) {
        _errorFn = errorFn;
        return this;
    };

    this.done = function (doneFn) {
        _doneFn = doneFn;
        return this;
    };

    this.timeout = function (timeout) {
        _timeout = timeout;
        return this;
    };

    this.delay = function (delay) {
        _delay = delay;
        return this;
    };

    this.progress = function (progressFn) {
        _progressFn = progressFn;
        return this;
    };

    this.post = function () {
        _method = "POST";
        setTimeout(this._post);

        return this;
    };

    this.get = function () {
        _method = "GET";
        setTimeout(this._post);

        return this;
    };

    this._post = function () {
        var params = _params;

        // 参数配置：https://github.com/axios/axios#request-config
        var config = {
            method: _method,
            url: Tea.url(_action),
            timeout: _timeout * 1000,
            headers: {
                "X-Requested-With": "XMLHttpRequest"
            }
        };

        if (_progressFn != null && typeof(_progressFn) == "function") {
            config["onUploadProgress"] = function (event) {
               _progressFn.call(Tea.Vue, event.loaded, event.total, event);
            };
        }

        if (_method === "GET") {
            config["params"] = params;
        }
        else {
            if (typeof(params) === "object" && params instanceof FormData) {
                Array.from(params).$each(function (name, object) {
                    if (object != null && object instanceof File) {
                        if (object.size === 0 && object.name.length === 0) {
                            params.delete(name);
                        }
                    }
                });
                config["data"] = params;
            }
            else {
                var formData = new FormData();
                for (var key in params) {
                    if (!params.hasOwnProperty(key)) {
                        continue;
                    }

                    if (params[key] == null) {
                        formData.append(key, "");
                    }
                    else {
                        formData.append(key, params[key]);
                    }
                }

                config["data"] = formData;
            }
        }

        axios(config)
            .then(function (response) {
                response = response.data;

                setTimeout(function () {
                    if (typeof(response) !== "object" || typeof(response.code) === "undefined") {
                        if (typeof(_errorFn) === "function") {
                            _errorFn.call(Tea.Vue, {});
                        }
                        return;
                    }

                    var code = parseInt(response.code, 10);
                    if (code === 200) {
                        if (typeof(_successFn) === "function") {
                            var result = _successFn.call(Tea.Vue, response);
                            if (typeof(result) === "boolean" && !result) {
                                return;
                            }
                        }

                        if (response.message != null && response.message.length > 0) {
                            alert(response.message);
                        }

                        if (response.next != null && typeof(response.next) === "object") {
                            if (response.next.action === "*refresh") {
                                window.location.reload();
                            }
                            else {
                                Tea.go(response.next.action, response.next.params, response.next.hash);
                            }
                        }
                    }
                    else {
                        if (typeof(_failFn) === "function") {
                            _failFn.call(Tea.Vue, response);
                        }
                        else {
                            Tea.failResponse(response);
                        }
                    }
                });
            })
            .catch(function (error) {
                console.log(error);

                if (typeof(_errorFn) === "function") {
                    _errorFn.call(Tea.Vue, {});
                }
            })
            .then(function () {
                // console.log("done");
                if (typeof(_doneFn) == "function") {
                    _doneFn.call(Tea.Vue, {});
                }
            });
    };
};

/**
 * 取得Action对象
 *
 * @param action Action
 * @returns {Window.Tea.Action}
 */
window.Tea.action = function (action) {
    return new this.Action(action);
};

/**
 * 激活元素中的Action
 *
 * 支持
 * - data-tea-action
 * - data-tea-confirm
 * - data-tea-timeout
 * - data-tea-before
 * - data-tea-success
 * - data-tea-fail
 * - data-tea-error
 * - data-tea-progress
 */
window.Tea.activate = function (element) {
    var nodes = Tea.element("*[data-tea-action]", element);
    if (nodes.length === 0) {
        return;
    }
    for (var i = 0; i < nodes.length; i++) {
        var node = nodes[i];

        if (node.tagName.toUpperCase() === "FORM") {
            Tea.element(node).unbind("submit").bind("submit", function (e) {
                Tea.runActionOn(this);

                e.preventDefault();
                e.stopPropagation();
            });
        }
        else {
            Tea.element(node).unbind("click").bind("click", function (e) {
                Tea.runActionOn(this);

                e.preventDefault();
                e.stopPropagation();

                return false;
            });
        }
    }
};

/**
 * 执行绑定data-tea-*的元素
 *
 * @param element 元素
 */
window.Tea.runActionOn = function (element) {
    var form = Tea.element(element);
    var action = form.attr("data-tea-action");
    var timeout = form.attr("data-tea-timeout");
    var confirm = form.attr("data-tea-confirm");
    var beforeFn = form.attr("data-tea-before");
    var successFn = form.attr("data-tea-success");
    var failFn = form.attr("data-tea-fail");
    var errorFn = form.attr("data-tea-error");
    var progressFn = form.attr("data-tea-progress");
    if (confirm != null && confirm.length > 0 && !window.confirm(confirm)) {
        return;
    }

    //执行前调用beforeFn
    if (beforeFn != null && beforeFn.length > 0) {
        beforeFn = beforeFn.split("(")[0].trim();
        if (typeof(Tea.Vue[beforeFn]) === "function") {
            var result = Tea.Vue[beforeFn].call(Tea.Vue, form);
            if (typeof(result) === "boolean" && !result) {
                return;
            }
        }
    }

    //请求对象
    var actionObject = Tea.action(action)
        .post();

    if (successFn != null && successFn.length > 0) {
        if (typeof(Tea.Vue[successFn]) === "function") {
            actionObject.success(function (resp) {
               Tea.Vue[successFn].call(Tea.Vue, resp);
            });
        }
    }

    if (failFn != null && failFn.length > 0) {
        if (typeof(Tea.Vue[failFn]) === "function") {
            actionObject.fail(function (resp) {
                Tea.Vue[failFn].call(Tea.Vue, resp);
            });
        }
    }

    if (errorFn != null && errorFn.length > 0) {
        if (typeof(Tea.Vue[errorFn]) === "function") {
            actionObject.error(function () {
                Tea.Vue[errorFn].call(Tea.Vue);
            });
        }
    }

    if (progressFn != null && progressFn.length > 0) {
        if (typeof(Tea.Vue[progressFn]) === "function") {
            actionObject.progress(function () {
                Tea.Vue[progressFn].apply(Tea.Vue, arguments);
            });
        }
    }

    //超时时间
    if (timeout != null) {
        timeout = parseFloat(timeout);
        if (!isNaN(timeout)) {
            actionObject.timeout(timeout);
        }
    }

    //参数
    if (element.tagName.toUpperCase() === "FORM") {
        actionObject.form(element);
    }
    else {
        var attributes = element.attributes;
        var params = {};
        for (var i = 0; i < attributes.length; i++) {
            var attr = attributes[i];
            var match = attr.name.toString().match(/^data-(.+)$/);
            if (match && !match[1].match(/^tea-/)) {
                var pieces = match[1].split("-");
                for (var j = 1; j < pieces.length; j++) {
                    pieces[j] = pieces[j][0].toUpperCase() + pieces[j].substr(1);
                }
                var name = pieces.join("");
                params[name] = attr.value;
            }
        }
        actionObject.params(params);
    }
};

var teaEventListeners = {}; // element => { event => [ callback1, ... ] }
function TeaElementObjects(elements) {
    var that = this;

    elements.$each(function (index, element) {
        that[index] = element;
    });

    this.bind = function (event, listener) {
        elements.$each(function (_, element) {
            if (typeof(teaEventListeners[element]) == "undefined") {
                teaEventListeners[element] = {};
            }
            if (typeof(teaEventListeners[element][event]) == "undefined") {
                teaEventListeners[element][event] = [];
            }
            teaEventListeners[element][event].push(listener);
            element.addEventListener(event, listener)
        });

        return this;
    };

    this.unbind = function (event) {
        elements.$each(function (_, element) {
            if (typeof(teaEventListeners[element]) == "undefined") {
                return;
            }
            if (typeof(teaEventListeners[element][event]) == "undefined") {
                return;
            }
            teaEventListeners[element][event].$each(function (_, listener) {
                element.removeEventListener(event, listener);
            });
            teaEventListeners[element][event] = [];
            var hasListeners = false;
            for (var k in teaEventListeners[element]) {
                if (!teaEventListeners[element].hasOwnProperty(k)) {
                    continue;
                }

                if (teaEventListeners[element][k] instanceof Array && teaEventListeners[element][k].length > 0) {
                    hasListeners = true;
                }
            }

            if (!hasListeners) {
                delete(teaEventListeners[element]);
            }
        });

        return this;
    };

    this.first = function () {
        var first = elements.$first;
        if (first != null) {
            return Tea.element(elements.$first());
        }
        return new TeaElementObjects([]);
    };

    this.attrs = function () {
        var first = this.first();
        if (first.length === 0) {
            return {};
        }

        var attrs = {};
        var node = first[0];
        for (var i = 0; i < node.attributes.length; i ++) {
            var attr = node.attributes[i];
            attrs[attr.name] = attr.value;
        }
        return attrs;
    };

    this.attr = function (name, value) {
        if (arguments.length === 0) {
            return "";
        }

        if (arguments.length === 1) {
            var attrs = this.attrs();
            if (typeof(attrs[name]) !== "undefined") {
                return attrs[name];
            }
            return "";
        }

        var first = this.first();
        if (first.length > 0) {
            first[0].setAttribute(name, value);
        }

        return this;
    };

    this.tagName = function () {
        var first = this.first();
        if (first.length === 0) {
            return "";
        }
        return first[0].tagName;
    };

    this.focus = function () {
        var first = this.first();
        if (first.length === 0) {
            return;
        }
        first[0].focus();
    };

    this.each = function (iterator) {
        elements.$each(function (index, element) {
            iterator(index, element);
        });

        return this;
    };

    this.find = function (selector) {
        if (this.length == 0) {
            return new TeaElementObjects([]);
        }
        return Tea.element(selector, this.first()[0]);
    };

    this.hide = function () {
        this.each(function (_, element) {
            element.style.display = "none";
        });
        return this;
    };

    this.show = function () {
        this.each(function (_, element) {
            element.style.display = "block";
        });
        return this;
    };

    this.text = function () {
        if (arguments.length > 0) {
            var text = arguments[0];
            this.each(function (_, element) {
                if (typeof(element.textContent) != "undefined") {
                    element.textContent = text;
                }
                if (typeof(element.innerText) != "undefined") {
                    element.innerText = text;
                }
            });
            return this;
        }

        if (this.length == 0) {
            return "";
        }
        if (typeof(elements[0].textContent) == "string") {
            return elements[0].textContent;
        }
        return elements[0].innerText;
    };

    this.html = function () {
        if (arguments.length > 0) {
            var html = arguments[0];
            this.each(function (_, element) {
                element.innerHTML = html;
            });
            return this;
        }

        if (this.length == 0) {
            return "";
        }
        return elements[0].innerHTML;
    };

    this.val = function () {
        if (arguments.length > 0) {
            var value = arguments[0];
            this.each(function (_, element) {
                element.value = value;
            });
            return this;
        }

        if (this.length == 0) {
            return "";
        }
        return elements[0].value;
    };

    this.length = elements.length;
}

/**
 * 获取元素
 *
 * @param selector 选择器
 * @param parent 父节点
 * @returns {*}
 */
window.Tea.element = function (selector, parent) {
    var elements = [];
    if (typeof(selector) === "object" && /(function|object) \w+Element\b/.test(selector.constructor.toString())) {
        elements = [selector];
    }
    else if (typeof(selector) === "object" && /function TeaElementObjects/.test(TeaElementObjects.constructor.toString())) {
        return selector;
    }
    else if (typeof(selector) === "string") {
        if (typeof(parent) === "object") {
            elements = Array.from(parent.querySelectorAll(selector));
        } else {
            elements = Array.from(document.querySelectorAll(selector));
        }
    }

    return new TeaElementObjects(elements);
};

/**
 * 生成vue for用的key
 * @returns {number}
 */
window.Tea.key = function () {
    return Math.random()
};

// 失败的响应处理
window.Tea.failResponse = function (response) {
    //消息提示
    var hasMessage = false;
    if (response.message != null && response.message.length > 0) {
        hasMessage = true;
        alert(response.message);
    }

    if (typeof(response.errors) === "object" && response.errors != null && response.errors.length > 0) {
        /**
         * errors: [
         *  [field1, [ error1, error2, ....]
         *  ...
         * ]
         * error: [ rule, message ]
         */
        var fieldName = response.errors[0].param;
        var error = response.errors[0].messages[0];
        if (!hasMessage) {
            alert(error);
        }
        var element = Tea.element("*[name='" + fieldName + "']");
        if (element) {
            element.focus();
        }
        else {
            var match = fieldName.match(/^(.+)\[(\d+)]$/);
            if (match != null) {
                var index = parseInt(match[2], 10);
                var fields = Tea.element("*[name='" + match[1].trim() + "[]']");
                if (fields.length > 0 && index < fields.length) {
                    fields[index].focus();
                }
            }
        }
    }
};

if (typeof(window.console) === "undefined") {
    window.console = {
        log: function () {},
        error: function () {},
        group: function () {}
    };
}