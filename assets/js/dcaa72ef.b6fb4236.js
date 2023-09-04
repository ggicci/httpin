"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[749],{4155:(e,t,a)=>{a.r(t),a.d(t,{assets:()=>d,contentTitle:()=>l,default:()=>u,frontMatter:()=>r,metadata:()=>o,toc:()=>p});var n=a(7462),i=(a(7294),a(3905));const r={sidebar_position:3},l="Upload Files",o={unversionedId:"advanced/upload-files",id:"advanced/upload-files",title:"Upload Files",description:"Introduced in v0.7.0.",source:"@site/docs/advanced/upload-files.md",sourceDirName:"advanced",slug:"/advanced/upload-files",permalink:"/httpin/advanced/upload-files",draft:!1,editUrl:"https://github.com/ggicci/httpin/edit/documentation/docs/docs/advanced/upload-files.md",tags:[],version:"current",sidebarPosition:3,frontMatter:{sidebar_position:3},sidebar:"docsSidebar",previous:{title:"Error Handler",permalink:"/httpin/advanced/error-handler"},next:{title:"Patch Field",permalink:"/httpin/advanced/patch"}},d={},p=[{value:"Access the uploaded file",id:"access-the-uploaded-file",level:2}],s={toc:p},c="wrapper";function u(e){let{components:t,...a}=e;return(0,i.kt)(c,(0,n.Z)({},s,a,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"upload-files"},"Upload Files"),(0,i.kt)("p",null,"Introduced in v0.7.0."),(0,i.kt)("p",null,"Use ",(0,i.kt)("a",{parentName:"p",href:"https://pkg.go.dev/github.com/ggicci/httpin#File"},(0,i.kt)("inlineCode",{parentName:"a"},"httpin.File"))," to retrieve a file uploaded from the request. Make sure it's a ",(0,i.kt)("a",{parentName:"p",href:"https://stackoverflow.com/q/4526273/1592264"},"multipart/form-data")," request."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go",metastring:"{4,5}","{4,5}":!0},'type UpdateArticleInput struct {\n    Title       string        `in:"form=title"`\n    IsPrivate   bool          `in:"form=is_private"`\n    Cover       httpin.File   `in:"form=cover"`\n    Attachments []httpin.File `in:"form=attachments"`\n}\n')),(0,i.kt)("p",null,(0,i.kt)("strong",{parentName:"p"},"NOTE"),": you ",(0,i.kt)("strong",{parentName:"p"},"MUST check")," ",(0,i.kt)("inlineCode",{parentName:"p"},"httpin.File.Valid")," before accessing."),(0,i.kt)("h2",{id:"access-the-uploaded-file"},"Access the uploaded file"),(0,i.kt)("p",null,"Access filename, filesize and other information from ",(0,i.kt)("inlineCode",{parentName:"p"},"httpin.File.Header"),", which is of type ",(0,i.kt)("a",{parentName:"p",href:"https://pkg.go.dev/mime/multipart#FileHeader"},(0,i.kt)("inlineCode",{parentName:"a"},"multipart.FileHeader")),"."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func UpdateArticle(rw http.ResponseWriter, r *http.Request) {\n    input := r.Context().Value(httpin.Input).(*UpdateArticleInput)\n\n    // User has uploaded a file for the cover.\n    if input.Cover.Valid {\n        filename := input.Cover.Header.Filename\n        filesize := input.Cover.Header.Size\n\n        // Read content.\n        fileBytes, err := ioutil.ReadAll(input.Cover)\n    }\n\n    // ...\n}\n")))}u.isMDXComponent=!0}}]);