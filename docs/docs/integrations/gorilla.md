---
sidebar_position: 2
---

# gorilla/mux 🦍

[gorilla/mux](https://github.com/gorilla/mux) is

> A powerful HTTP router and URL matcher for building Go web servers

## path Directive by `mux.Vars`

```go {4}
func init() {
	// Register a directive named "path" to retrieve values from `mux.Vars`,
	// i.e. decode path variables.
	httpin.UseGorillaMux("path", mux.Vars)
}

type GetPostOfUserInput struct {
	Username string `in:"path=username"` // equivalent to `mux.Vars(r)["username"]`
	PostID   int64  `in:"path=pid"`
}
```
