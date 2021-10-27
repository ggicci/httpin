---
sidebar_position: 0
---

# net/http

Package [net/http](https://pkg.go.dev/net/http#Handler)

> provides HTTP client and server implementations.

## Chain httpin's Middlware to your http.Handler(s)

We recommend using [justinas/alice](https://github.com/justinas/alice) to chain your middlewares.

```go {4,10}
// Bind input vs. handler.
func init() {
	http.Handle("/users", alice.New(
		httpin.NewInput(ListUsersInput{}),
	).ThenFunc(ListUsers))
}

// Get your input instance with only ONE LINE of code.
func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*ListUsersInput)

	// Do sth.
}
```
