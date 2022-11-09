"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[399],{8613:function(t,e,a){a.r(e),a.d(e,{assets:function(){return d},contentTitle:function(){return s},default:function(){return l},frontMatter:function(){return p},metadata:function(){return c},toc:function(){return h}});var i=a(3117),n=a(102),r=(a(7294),a(3905)),o=["components"],p={sidebar_position:5},s="path",c={unversionedId:"directives/path",id:"directives/path",title:"path",description:"path is a directive executor who decodes a field from the path of the request URI, aka. path variables.",source:"@site/docs/directives/path.mdx",sourceDirName:"directives",slug:"/directives/path",permalink:"/httpin/directives/path",draft:!1,editUrl:"https://github.com/ggicci/httpin/edit/documentation/docs/docs/directives/path.mdx",tags:[],version:"current",sidebarPosition:5,frontMatter:{sidebar_position:5},sidebar:"docsSidebar",previous:{title:"body",permalink:"/httpin/directives/body"},next:{title:"required",permalink:"/httpin/directives/required"}},d={},h=[],u={toc:h};function l(t){var e=t.components,a=(0,n.Z)(t,o);return(0,r.kt)("wrapper",(0,i.Z)({},u,a,{components:e,mdxType:"MDXLayout"}),(0,r.kt)("h1",{id:"path"},"path"),(0,r.kt)("p",null,(0,r.kt)("strong",{parentName:"p"},"path")," is a ",(0,r.kt)("a",{parentName:"p",href:"/advanced/concepts#directive-executor"},"directive executor")," who decodes a field from the path of the request URI, aka. path variables."),(0,r.kt)("div",{className:"admonition admonition-danger alert alert--danger"},(0,r.kt)("div",{parentName:"div",className:"admonition-heading"},(0,r.kt)("h5",{parentName:"div"},(0,r.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,r.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"12",height:"16",viewBox:"0 0 12 16"},(0,r.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M5.05.31c.81 2.17.41 3.38-.52 4.31C3.55 5.67 1.98 6.45.9 7.98c-1.45 2.05-1.7 6.53 3.53 7.7-2.2-1.16-2.67-4.52-.3-6.61-.61 2.03.53 3.33 1.94 2.86 1.39-.47 2.3.53 2.27 1.67-.02.78-.31 1.44-1.13 1.81 3.42-.59 4.78-3.42 4.78-5.56 0-2.84-2.53-3.22-1.25-5.61-1.52.13-2.03 1.13-1.89 2.75.09 1.08-1.02 1.8-1.86 1.33-.67-.41-.66-1.19-.06-1.78C8.18 5.31 8.68 2.45 5.05.32L5.03.3l.02.01z"}))),"danger")),(0,r.kt)("div",{parentName:"div",className:"admonition-content"},(0,r.kt)("p",{parentName:"div"},(0,r.kt)("strong",{parentName:"p"},"httpin")," doesn't provide a builtin ",(0,r.kt)("strong",{parentName:"p"},"path")," directive, because ",(0,r.kt)("strong",{parentName:"p"},"httpin")," doesn't provide routing functionality.\nBut ",(0,r.kt)("strong",{parentName:"p"},"httpin")," can be easily integrated with other packages that provide routing functionality, to decode path variables."))),(0,r.kt)("p",null,"You can quickly implement a ",(0,r.kt)("strong",{parentName:"p"},"path")," directive with the following code (routing package specific):"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},"go-chi/chi")),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},'httpin.UseGochiURLParam("path", chi.URLParam)\n')),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},"gorilla/mux")),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},'httpin.UseGorillaMux("path", mux.Vars)\n')),(0,r.kt)("p",null,"You could visit the pages under the ",(0,r.kt)("strong",{parentName:"p"},"Integrations")," section in the sidebar to find more details on how to integrate ",(0,r.kt)("strong",{parentName:"p"},"httpin")," with other packages."),(0,r.kt)("p",null,"If you can't find the package you wanted in the list, you could either open an issue on the ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/ggicci/httpin/issues"},(0,r.kt)("strong",{parentName:"a"},"Github")),"\nor visit ",(0,r.kt)("a",{parentName:"p",href:"/directives/custom"},"custom \ud83d\udd0c"),' to learn to implement a "path" directive of your own.'),(0,r.kt)("p",null,"We also hope that you can make contributions to the ",(0,r.kt)("strong",{parentName:"p"},"httpin")," project to make it great! Thanks in advance \u2764\ufe0f"))}l.isMDXComponent=!0}}]);