(this.webpackJsonpts=this.webpackJsonpts||[]).push([[0],{17:function(e,t,n){},22:function(e,t,n){"use strict";n.r(t);var a=n(0),c=n.n(a),o=n(6),i=n.n(o),r=(n(17),n(3)),s=n(8),u=n(2);function l(e){return Object(u.jsx)(s.a,{value:e.value,onChange:e.onChange,height:"80vh",options:{formatOnPaste:!0,formatOnType:!0,minimap:{enabled:!1},stablePeek:!0,suggest:{preview:!0}},language:"json",onMount:function(t,n){var a=n.Uri.parse("a://b/foo.json"),c=n.editor.createModel("","json",a);n.languages.json.jsonDefaults.setDiagnosticsOptions({validate:!0,schemaValidation:"error",schemas:[{uri:"a://b/foo.json",schema:JSON.parse(e.schema),fileMatch:[a.toString()]}]}),t.setModel(c),e.onMount(t)},overrideServices:{storageService:{get:function(){},getNumber:function(){},getBoolean:function(e){return"expandSuggestionDocs"===e},remove:function(){},store:function(){},onDidChangeStorage:function(){},onWillSaveState:function(){}}}})}function d(e){return Object(u.jsx)("a",{style:{background:"none",border:"none",margin:0,padding:3,display:"flex",flexDirection:"row",alignItems:"center",justifyContent:"center",cursor:"pointer",fontSize:13},onClick:function(){window.__openLink("https://github.com/wirekang/mouseable")},children:Object(u.jsx)("img",{alt:"github",src:"github.png",width:50,height:20})})}function j(e){return Object(u.jsxs)("div",{style:{fontSize:10},children:["You can close this window safely.",Object(u.jsx)("button",{style:{fontSize:12},onClick:function(){window.__terminate(),window.close()},children:"Terminate Program"})]})}function b(e){return Object(u.jsxs)("div",{style:{display:"flex",flexDirection:"row",alignItems:"center",justifyContent:"space-between",margin:3},children:[Object(u.jsxs)("span",{children:["Version: ",e.version]}),Object(u.jsx)(d,{}),Object(u.jsx)(j,{})]})}function f(e){return Object(u.jsx)("div",{style:{fontSize:12},children:Object(u.jsxs)("ul",{children:[Object(u.jsxs)("li",{children:["Press ",Object(u.jsx)("b",{children:"Ctrl+I"})," to show suggestions."]}),Object(u.jsxs)("li",{children:["Press ",Object(u.jsx)("b",{children:"F1"})," to insert key."]}),Object(u.jsxs)("li",{children:["Press ",Object(u.jsx)("b",{children:"F2"})," to save."]}),Object(u.jsxs)("li",{children:["Press ",Object(u.jsx)("b",{children:"F3"})," to apply."]})]})})}var p=n(23),h=n(11),g=n(24),O=n(25);function v(e){var t=Object(O.a)(e.delay);return(0,Object(r.a)(t,1)[0])()?Object(u.jsx)("h2",{style:{position:"absolute",color:"red",width:"100vw",height:"100vh",backgroundColor:"white"},children:"Mouseable is not running."}):Object(u.jsx)(c.a.Fragment,{})}function x(e){return Object(u.jsx)("div",{style:{display:"flex",flexDirection:"row",justifyContent:"space-between",alignContent:"center",margin:0,padding:0},children:e.children})}function m(e){var t;return Object(u.jsxs)("div",{style:{},children:[Object(u.jsx)("select",{value:e.loadedConfigName,onChange:function(t){e.onLoadConfig(t.target.value)},children:null===(t=e.configNames)||void 0===t?void 0:t.map((function(e){return Object(u.jsx)("option",{children:e},e)}))}),Object(u.jsx)("br",{}),"Applied: ",e.appliedConfigName]})}var w=n(7),y=n(5),C=n.n(y),k=n(4);function _(){return(_=Object(w.a)(C.a.mark((function e(t,n){var a,c,o,i;return C.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,Object(k.a)("Press any keys include double press."),e.next=4,window.__getNextKey();case 4:if(a=e.sent){e.next=7;break}return e.abrupt("return");case 7:if(c=t.getPosition()){e.next=10;break}return e.abrupt("return");case 10:o=new n.Range(c.lineNumber,c.column,c.lineNumber,c.column),i={identifier:{major:1,minor:1},range:o,text:'"'.concat(a,'"')},t.executeEdits("my-source",[i]),e.next=19;break;case 16:e.prev=16,e.t0=e.catch(0),Object(k.a)("".concat(e.t0));case 19:case"end":return e.stop()}}),e,null,[[0,16]])})))).apply(this,arguments)}function S(){return(S=Object(w.a)(C.a.mark((function e(t,n,a){return C.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:if(e.prev=0,a.trigger("anyString","editor.action.formatDocument",null),t){e.next=5;break}return Object(k.a)("Nothing was loaded."),e.abrupt("return");case 5:return e.next=7,window.__saveConfig(t,null!==n&&void 0!==n?n:"");case 7:Object(k.a)("Saved"),e.next=13;break;case 10:e.prev=10,e.t0=e.catch(0),Object(k.a)("".concat(e.t0));case 13:case"end":return e.stop()}}),e,null,[[0,10]])})))).apply(this,arguments)}function N(){return(N=Object(w.a)(C.a.mark((function e(t,n,a){var c;return C.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,window.__loadConfig(t);case 3:c=e.sent,n(t),a(c),e.next=11;break;case 8:e.prev=8,e.t0=e.catch(0),Object(k.a)("".concat(e.t0));case 11:case"end":return e.stop()}}),e,null,[[0,8]])})))).apply(this,arguments)}function F(){return(F=Object(w.a)(C.a.mark((function e(t){return C.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:if(e.prev=0,t){e.next=4;break}return Object(k.a)("Nothing was loaded."),e.abrupt("return");case 4:return e.next=6,window.__applyConfig(t);case 6:Object(k.a)("".concat(t," Applied")),e.next=12;break;case 9:e.prev=9,e.t0=e.catch(0),Object(k.a)("".concat(e.t0));case 12:case"end":return e.stop()}}),e,null,[[0,9]])})))).apply(this,arguments)}var P=function(){var e,t=Object(p.a)(window.__getVersion),n=Object(h.a)(window.__ping),c=Object(r.a)(n,2),o=c[0],i=c[1],d=Object(p.a)(window.__loadSchema),j=Object(h.a)(window.__loadAppliedConfigName),O=Object(r.a)(j,2),w=O[0],y=O[1],C=Object(h.a)(window.__loadConfigNames),k=Object(r.a)(C,2),P=k[0],D=k[1],E=Object(a.useState)(),M=Object(r.a)(E,2),L=M[0],z=M[1],I=Object(a.useState)(),K=Object(r.a)(I,2),A=K[0],J=K[1],V=Object(a.useState)(),B=Object(r.a)(V,2),T=B[0],q=B[1],R=Object(s.b)(),U=function(e){!function(e,t,n){N.apply(this,arguments)}(e,q,z)},W=function(){A&&R?function(e,t){_.apply(this,arguments)}(A,R):console.log(A,R)},Y=function(){A&&function(e,t,n){S.apply(this,arguments)}(T,L,A)},G=function(){!function(e){F.apply(this,arguments)}(T)};return Object(a.useEffect)((function(){y(),D()}),[]),Object(a.useEffect)((function(){w.value&&!T&&U(w.value)}),[w.loading,U,T]),Object(g.a)((function(){i()}),2e3),Object(a.useEffect)((function(){w.value}),[w.value]),Object(a.useEffect)((function(){A&&R&&(A.addCommand(R.KeyCode.F1,W),A.addCommand(R.KeyCode.F2,Y),A.addCommand(R.KeyCode.F3,G))}),[A,R,W,Y,G]),Object(u.jsxs)("div",{style:{height:"100%"},children:[o.loading&&Object(u.jsx)(v,{delay:1e3}),Object(u.jsx)(b,{version:null!==(e=t.value)&&void 0!==e?e:""}),Object(u.jsxs)(x,{children:[Object(u.jsx)(f,{}),Object(u.jsx)(m,{configNames:P.value,loadedConfigName:T,onLoadConfig:U,appliedConfigName:w.value})]}),d.value&&Object(u.jsx)(l,{value:L,onChange:z,schema:d.value,onMount:J})]})};window.addEventListener("keydown",(function(e){"F1"!==e.key&&"F2"!==e.key&&"F3"!==e.key||(e.preventDefault(),Object(k.a)("Press again. (Changed focus to editor)"),document.querySelector("textarea.monaco-mouse-cursor-text").focus())})),window.addEventListener("error",(function(e){Object(k.a)(e.message)})),i.a.render(Object(u.jsx)(c.a.StrictMode,{children:Object(u.jsx)(P,{})}),document.getElementById("root"))}},[[22,1,2]]]);
//# sourceMappingURL=main.c492db1b.chunk.js.map