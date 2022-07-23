---
sidebar_position: 4
---

# go-restful

[**go-restful**](https://github.com/emicklei/go-restful) is a

> package for building REST-style Web Services using Go.

## Convert `httpin.NewInput` middleware handler to `restful.Filter`

Use [HttpMiddlewareHandlerToFilter](https://pkg.go.dev/github.com/emicklei/go-restful/v3#HttpMiddlewareHandlerToFilter), which is introduced in [v3.9.0](https://github.com/emicklei/go-restful/tree/v3.9.0) by this [PR#505](https://github.com/emicklei/go-restful/pull/505).

```go {8,15}
type ListUsersInput struct {
	Gender  string `in:"query=gender"`
	Page    int    `in:"query=page"`
	PerPage int    `in:"query=per_page,page_size"`
}

func handleListUsers(request *restful.Request, response *restful.Response) {
	input := request.Request.Context().Value(httpin.Input).(*ListUsersInput)
	// ...
}

func main() {
	ws := new(WebService)
	ws.Route(ws.GET("/users").Filter(
		restful.HttpMiddlewareHandlerToFilter(httpin.NewInput(ListUsersInput{})),
	).To(handleListUsers))
}
```
