# httpin

[![codecov](https://codecov.io/gh/ggicci/httpin/branch/main/graph/badge.svg?token=RT61L9ngHj)](https://codecov.io/gh/ggicci/httpin)

HTTP Input for Go

## What is this lib for?

TODO(ggicci): image

## Quick Start

Suppose that we have an RESTful API for querying users. We can define a struct to collect all the input parameters like:

```go
type Authorization struct {
	Token string `in:"query.access_token,header.x-api-token"`
}

type Pagination struct {
	Page    int `in:"query.page"`
	PerPage int `in:"query.per_page"`
}

type UserQuery struct {
	Gender   string `in:"query.gender"`
	AgeRange []int  `in:"query.age_range"`
	IsMember bool   `in:"query.is_member"`
	Pagination
	Authorization
}
```

Use `httpin` to extract the data from `http.Request` for you:

```go
httpinUserQuery, err := httpin.NewEngine(UserQuery{})
if err != nil {
    // You can create your engines at server start up.
    // And reuse them in your http handlers.
}

func QueryUser(rw http.ResponseWriter, r *http.Request) {
    input, err := httpinUserQuery.Read(r)
    if err != nil {
        http.Error(rw, err, http.StatusBadRequest)
        return
    }
    userQuery := input.(*UserQuery)

    if !NewAccess(userQuery.Token).Accessible {
        http.Error(rw, "invalid token", http.StatusUnauthorized)
        return
    }

    // Use `userQuery` to get the parameters you wanted.
    // ...
}

```

## Advanced Usage - Use Middleware

Firstly setup httpin middleware for you APIs.

```go
todo
```

And then fetch your input from the context values of the request.

```go
func QueryUser(rw http.ResponseWriter, r *http.Request) {
    var userQuery = r.Context().Value(httpin.Input).(*UserQuery)
    // Use userQuery here.
    // ...
}
```
