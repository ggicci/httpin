"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[520],{2332:(e,t,n)=>{n.d(t,{Z:()=>p});var r=function(e,t,n,r){return new(n||(n=Promise))((function(i,s){function o(e){try{l(r.next(e))}catch(t){s(t)}}function a(e){try{l(r.throw(e))}catch(t){s(t)}}function l(e){var t;e.done?i(e.value):(t=e.value,t instanceof n?t:new n((function(e){e(t)}))).then(o,a)}l((r=r.apply(e,t||[])).next())}))};class i{constructor(e="/goplay"){this.proxyUrl=e}raiseForStatus(e){return r(this,void 0,void 0,(function*(){if(!e.ok){const t=yield e.text();throw new Error(t?e.statusText+": "+t:e.statusText)}}))}compile(e,t){return r(this,void 0,void 0,(function*(){const n=new FormData;n.append("version","2"),n.append("withVet","true"),n.append("body",e);const r=yield fetch(`${this.proxyUrl}/_/compile?backend=${t||""}`,{method:"POST",body:n});return yield this.raiseForStatus(r),yield r.json()}))}renderCompile(e,t,n){return r(this,void 0,void 0,(function*(){e.replaceChildren(this.renderMessage("system","Waiting for remote server..."));const r=yield this.compile(t,n);if(e.replaceChildren(),""!=r.Errors)return e.appendChild(this.renderMessage("error",r.Errors)),void e.appendChild(this.renderMessage("system","\nGo build failed."));for(const t of r.Events||[])e.appendChild(yield this.renderEvent(t));e.appendChild(this.renderMessage("system","\nProgram exited."))}))}renderEvent(e){return r(this,void 0,void 0,(function*(){var t;return e.Delay>=0&&(yield(t=e.Delay/1e6,new Promise((e=>setTimeout(e,t))))),this.renderMessage(e.Kind,e.Message)}))}renderMessage(e,t){const n=document.createElement("span");return n.classList.add(e),n.innerText=t,n}share(e,t){return r(this,void 0,void 0,(function*(){const n=yield fetch(`${this.proxyUrl}/_/share`,{method:"POST",body:e});yield this.raiseForStatus(n);const r="https://go.dev/play/p/"+(yield n.text());return t?`${r}?v=${t}`:r}))}}var s=n(4464),o=n(7294);const a="toolbar_uIxz",l="button_yMrS",u="hidden_X41c",d="https://goplay.ggicci.me",c=e=>{const{children:t,onClick:n}=e;return o.createElement("button",{className:l,onClick:n},t)},p=e=>{const{children:t}=e,n=o.useRef(null),r=o.useRef(null),l=t&&t.props&&"pre"===t.props.mdxType&&t,p=l&&l.props&&l.props.children;if(!p||"code"!==p.props.mdxType)return o.createElement("div",null,"GoPlay: the wrapped data is not a codeblock.");if(!/\blanguage-go\b/.test(p&&p.props.className))return o.createElement("div",null,"GoPlay: only go code supported.");return o.createElement(o.Fragment,null,t,o.createElement("div",{ref:n,className:u},o.createElement(s.Z,{language:"text"},o.createElement("div",{ref:r}))),o.createElement("div",{className:a},o.createElement(c,{onClick:function(){if(r.current){const e=new i(d);n.current.classList.remove(u),e.renderCompile(r.current,p.props.children.trim())}}},"Run"),o.createElement(c,{onClick:async function(){const e=new i(d),t=await e.share(p.props.children.trim());window.open(t,"_blank")}},"Try it yourself \u21e2")))}},1893:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>u,contentTitle:()=>a,default:()=>m,frontMatter:()=>o,metadata:()=>l,toc:()=>d});var r=n(7462),i=(n(7294),n(3905)),s=n(2332);const o={sidebar_position:4},a="go-restful",l={unversionedId:"integrations/go-restful",id:"integrations/go-restful",title:"go-restful",description:"go-restful is a",source:"@site/docs/integrations/go-restful.mdx",sourceDirName:"integrations",slug:"/integrations/go-restful",permalink:"/httpin/integrations/go-restful",draft:!1,editUrl:"https://github.com/ggicci/httpin/edit/documentation/docs/docs/integrations/go-restful.mdx",tags:[],version:"current",sidebarPosition:4,frontMatter:{sidebar_position:4},sidebar:"docsSidebar",previous:{title:"gin-gonic/gin \ud83e\udd64",permalink:"/httpin/integrations/gin"},next:{title:"Concepts",permalink:"/httpin/advanced/concepts"}},u={},d=[{value:"Integrations",id:"integrations",level:2},{value:"Run Demo",id:"run-demo",level:2}],c={toc:d},p="wrapper";function m(e){let{components:t,...n}=e;return(0,i.kt)(p,(0,r.Z)({},c,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"go-restful"},"go-restful"),(0,i.kt)("p",null,(0,i.kt)("a",{parentName:"p",href:"https://github.com/emicklei/go-restful"},(0,i.kt)("strong",{parentName:"a"},"go-restful"))," is a"),(0,i.kt)("blockquote",null,(0,i.kt)("p",{parentName:"blockquote"},"package for building REST-style Web Services using Go.")),(0,i.kt)("h2",{id:"integrations"},"Integrations"),(0,i.kt)("p",null,"Convert ",(0,i.kt)("inlineCode",{parentName:"p"},"httpin.NewInput")," middleware handler to ",(0,i.kt)("inlineCode",{parentName:"p"},"restful.Filter"),"."),(0,i.kt)("p",null,"Use ",(0,i.kt)("a",{parentName:"p",href:"https://pkg.go.dev/github.com/emicklei/go-restful/v3#HttpMiddlewareHandlerToFilter"},"HttpMiddlewareHandlerToFilter"),", which is introduced in ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/emicklei/go-restful/tree/v3.9.0"},"v3.9.0")," by this ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/emicklei/go-restful/pull/505"},"PR#505"),"."),(0,i.kt)("h2",{id:"run-demo"},"Run Demo"),(0,i.kt)(s.Z,{mdxType:"GoPlay"},(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go",metastring:"{20,32}","{20,32}":!0},'package main\n\nimport (\n    "fmt"\n    "net/http"\n    "net/http/httptest"\n\n    restful "github.com/emicklei/go-restful/v3"\n    "github.com/ggicci/httpin"\n)\n\ntype ListUsersInput struct {\n    Gender   string `in:"query=gender"`\n    AgeRange []int  `in:"query=age_range"`\n    IsMember bool   `in:"query=is_member"`\n}\n\nfunc handleListUsers(request *restful.Request, response *restful.Response) {\n    // Retrieve you data in one line of code!\n    input := request.Request.Context().Value(httpin.Input).(*ListUsersInput)\n\n    fmt.Printf("input: %#v\\n", input)\n}\n\nfunc main() {\n    ws := new(restful.WebService)\n\n    wsContainer := restful.NewContainer()\n    wsContainer.Add(ws)\n\n    // Bind input struct with handler.\n    ws.Route(ws.GET("/users").Filter(\n        restful.HttpMiddlewareHandlerToFilter(httpin.NewInput(ListUsersInput{})),\n    ).To(handleListUsers))\n\n    r, _ := http.NewRequest("GET", "/users?gender=male&age_range=18&age_range=24&is_member=1", nil)\n\n    rw := httptest.NewRecorder()\n    wsContainer.ServeHTTP(rw, r)\n}\n'))))}m.isMDXComponent=!0}}]);