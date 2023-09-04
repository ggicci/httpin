"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[68],{2332:(e,t,n)=>{n.d(t,{Z:()=>u});var r=function(e,t,n,r){return new(n||(n=Promise))((function(i,s){function o(e){try{p(r.next(e))}catch(t){s(t)}}function a(e){try{p(r.throw(e))}catch(t){s(t)}}function p(e){var t;e.done?i(e.value):(t=e.value,t instanceof n?t:new n((function(e){e(t)}))).then(o,a)}p((r=r.apply(e,t||[])).next())}))};class i{constructor(e="/goplay"){this.proxyUrl=e}raiseForStatus(e){return r(this,void 0,void 0,(function*(){if(!e.ok){const t=yield e.text();throw new Error(t?e.statusText+": "+t:e.statusText)}}))}compile(e,t){return r(this,void 0,void 0,(function*(){const n=new FormData;n.append("version","2"),n.append("withVet","true"),n.append("body",e);const r=yield fetch(`${this.proxyUrl}/_/compile?backend=${t||""}`,{method:"POST",body:n});return yield this.raiseForStatus(r),yield r.json()}))}renderCompile(e,t,n){return r(this,void 0,void 0,(function*(){e.replaceChildren(this.renderMessage("system","Waiting for remote server..."));const r=yield this.compile(t,n);if(e.replaceChildren(),""!=r.Errors)return e.appendChild(this.renderMessage("error",r.Errors)),void e.appendChild(this.renderMessage("system","\nGo build failed."));for(const t of r.Events||[])e.appendChild(yield this.renderEvent(t));e.appendChild(this.renderMessage("system","\nProgram exited."))}))}renderEvent(e){return r(this,void 0,void 0,(function*(){var t;return e.Delay>=0&&(yield(t=e.Delay/1e6,new Promise((e=>setTimeout(e,t))))),this.renderMessage(e.Kind,e.Message)}))}renderMessage(e,t){const n=document.createElement("span");return n.classList.add(e),n.innerText=t,n}share(e,t){return r(this,void 0,void 0,(function*(){const n=yield fetch(`${this.proxyUrl}/_/share`,{method:"POST",body:e});yield this.raiseForStatus(n);const r="https://go.dev/play/p/"+(yield n.text());return t?`${r}?v=${t}`:r}))}}var s=n(4464),o=n(7294);const a="toolbar_uIxz",p="button_yMrS",c="hidden_X41c",d="https://goplay.ggicci.me",l=e=>{const{children:t,onClick:n}=e;return o.createElement("button",{className:p,onClick:n},t)},u=e=>{const{children:t}=e,n=o.useRef(null),r=o.useRef(null),p=t&&t.props&&"pre"===t.props.mdxType&&t,u=p&&p.props&&p.props.children;if(!u||"code"!==u.props.mdxType)return o.createElement("div",null,"GoPlay: the wrapped data is not a codeblock.");if(!/\blanguage-go\b/.test(u&&u.props.className))return o.createElement("div",null,"GoPlay: only go code supported.");return o.createElement(o.Fragment,null,t,o.createElement("div",{ref:n,className:c},o.createElement(s.Z,{language:"text"},o.createElement("div",{ref:r}))),o.createElement("div",{className:a},o.createElement(l,{onClick:function(){if(r.current){const e=new i(d);n.current.classList.remove(c),e.renderCompile(r.current,u.props.children.trim())}}},"Run"),o.createElement(l,{onClick:async function(){const e=new i(d),t=await e.share(u.props.children.trim());window.open(t,"_blank")}},"Try it yourself \u21e2")))}},1547:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>a,default:()=>h,frontMatter:()=>o,metadata:()=>p,toc:()=>d});var r=n(7462),i=(n(7294),n(3905)),s=n(2332);const o={sidebar_position:0},a="net/http",p={unversionedId:"integrations/http",id:"integrations/http",title:"net/http",description:"Package net/http",source:"@site/docs/integrations/http.mdx",sourceDirName:"integrations",slug:"/integrations/http",permalink:"/httpin/integrations/http",draft:!1,editUrl:"https://github.com/ggicci/httpin/edit/documentation/docs/docs/integrations/http.mdx",tags:[],version:"current",sidebarPosition:0,frontMatter:{sidebar_position:0},sidebar:"docsSidebar",previous:{title:"Create Your Own \ud83d\udd0c",permalink:"/httpin/directives/custom"},next:{title:"go-chi/chi",permalink:"/httpin/integrations/gochi"}},c={},d=[{value:"Integrations",id:"integrations",level:2},{value:"Run Demo",id:"run-demo",level:2}],l={toc:d},u="wrapper";function h(e){let{components:t,...n}=e;return(0,i.kt)(u,(0,r.Z)({},l,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"nethttp"},"net/http"),(0,i.kt)("p",null,"Package ",(0,i.kt)("a",{parentName:"p",href:"https://pkg.go.dev/net/http#Handler"},"net/http")),(0,i.kt)("blockquote",null,(0,i.kt)("p",{parentName:"blockquote"},"provides HTTP client and server implementations.")),(0,i.kt)("h2",{id:"integrations"},"Integrations"),(0,i.kt)("p",null,"Chain httpin's Middlware to your ",(0,i.kt)("inlineCode",{parentName:"p"},"http.Handler"),"s. We recommend using ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/justinas/alice"},"justinas/alice")," to chain your middlewares."),(0,i.kt)("h2",{id:"run-demo"},"Run Demo"),(0,i.kt)(s.Z,{mdxType:"GoPlay"},(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go",metastring:"{20,27}","{20,27}":!0},'package main\n\nimport (\n    "fmt"\n    "net/http"\n    "net/http/httptest"\n\n    "github.com/ggicci/httpin"\n    "github.com/justinas/alice"\n)\n\ntype ListUsersInput struct {\n    Gender   string `in:"query=gender"`\n    AgeRange []int  `in:"query=age_range"`\n    IsMember bool   `in:"query=is_member"`\n}\n\nfunc ListUsers(rw http.ResponseWriter, r *http.Request) {\n    // Retrieve you data in one line of code!\n    input := r.Context().Value(httpin.Input).(*ListUsersInput)\n\n    fmt.Printf("input: %#v\\n", input)\n}\n\nfunc init() {\n    // Bind input struct with handler.\n    http.Handle("/users", alice.New(\n        httpin.NewInput(ListUsersInput{}),\n    ).ThenFunc(ListUsers))\n}\n\nfunc main() {\n    r, _ := http.NewRequest("GET", "/users?gender=male&age_range=18&age_range=24&is_member=1", nil)\n\n    rw := httptest.NewRecorder()\n    http.DefaultServeMux.ServeHTTP(rw, r)\n}\n'))))}h.isMDXComponent=!0}}]);