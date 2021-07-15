# httpin

![Go Workflow](https://github.com/ggicci/httpin/actions/workflows/go.yml/badge.svg) [![codecov](https://codecov.io/gh/ggicci/httpin/branch/main/graph/badge.svg?token=RT61L9ngHj)](https://codecov.io/gh/ggicci/httpin) [![Go Reference](https://pkg.go.dev/badge/github.com/ggicci/httpin.svg)](https://pkg.go.dev/github.com/ggicci/httpin)

HTTP Input for Go - Decode an HTTP request into a custom struct

**Define the struct for your input and then fetch your data!**

## Quick View

<table>
<tr>
  <th>BEFORE (use net/http)</th>
  <th>AFTER (use httpin)</th>
</tr>
<tr>
  <td>

```go
func ListUsers(rw http.ResponseWriter, r *http.Request) {
	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		// Invalid parameter: page.
		return
	}
	perPage, err := strconv.ParseInt(r.FormValue("per_page"), 10, 64)
	if err != nil {
		// Invalid parameter: per_page.
		return
	}
	isMember, err := strconv.ParseBool(r.FormValue("is_member"))
	if err != nil {
		// Invalid parameter: is_member.
		return
	}

	// Do sth.
}
```

  </td>
  <td>

```go
type ListUsersInput struct {
	Page     int  `in:"form=page"`
	PerPage  int  `in:"form=per_page"`
	IsMember bool `in:"form=is_member"`
}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*ListUsersInput)
	// Do sth.
}
```

  </td>
</tr>
</table>

## Features

- [x] Builtin directive `form` to decode a field from HTTP query (URL params), i.e. `http.Request.Form`
- [x] Builtin directive `header` to decode a field from HTTP headers, i.e. `http.Request.Header`
- [x] Builtin decoders used by `form` and `header` directives for basic types, e.g. `bool`, `int`, `int64`, `float32`, `time.Time`, ... [full list](./internal/decoders.go)
- [x] Decode a field by inspecting a set of keys from the same source, e.g. `in:"form=per_page,page_size"`
- [x] Decode a field from multiple sources, e.g. both query and headers, `in:"form=access_token;header=x-api-token"`
- [x] Register custom type decoders by implementing `httpin.Decoder` interface
- [x] Compose an input struct by embedding struct fields
- [x] Builtin directive `required` to tag a field as **required**
- [x] Register custom directive executors to extend the ability of field resolving, see directive [required](./required.go) as an example and think about implementing your own directives like `trim`, `to_lowercase`, `base58_to_int`, etc.
- [x] Easily integrating with popular Go web frameworks and packages

## Sample User Defined Input Structs

```go
type Authorization struct {
	// Decode from multiple sources, the former with higher priority
	Token string `in:"form=access_token;header=x-api-token;required"`
}

type Pagination struct {
	Page int `in:"form=page"`

	// Decode from multiple keys in the same source, the former with higher priority
	PerPage int `in:"form=per_page,page_size"`
}

type ListUsersInput struct {
	Gender   string `in:"form=gender"`
	AgeRange []int  `in:"form=age_range"`
	IsMember bool   `in:"form=is_member"`

	Pagination    // Embedded field works
	Authorization // Embedded field works
}
```

## Integrate with Go Native http.Handler (Use Middleware)

First, set up the middleware for your handlers (**bind Input vs. Handler**). We recommend using [alice](https://github.com/justinas/alice) to chain your HTTP middleware functions.

```go
func init() {
	http.Handle("/users", alice.New(
		httpin.NewInput(ListUsersInput{}),
	).ThenFunc(ListUsers))
}
```

Second, fetch your input with only **ONE LINE** of code.

```go
func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*ListUsersInput)

	// Do sth.
}
```

## Integrate with Popular Go Web Frameworks and Components

### Frameworks

- [go-chi/chi](https://github.com/ggicci/httpin/wiki/Integrate-with-gochi)
- [gin-gonic/gin](https://github.com/ggicci/httpin/wiki/Integrate-with-gin)
- ...

### Components

- [HTTP Router: gorilla/mux](https://github.com/ggicci/httpin/wiki/Integrate-with-gorilla-mux)

## Advanced

### 🔥 Extend `httpin` by adding custom directives

Know the concept of a `Directive`:

```go
type Authorization struct {
	Token string `in:"form=access_token,token;header=x-api-token;required"`
	                  ^---------------------^ ^----------------^ ^------^
	                            d1                    d2            d3
}
```

There are three directives above, separated by semicolons (`;`):

- d1: `form=access_token,token`
- d2: `header=x-api-token`
- d3: `required`

A directive consists of two parts separated by an equal sign (`=`). The left part is the name of an executor who was designed to run this directive. The right part is a list of arguments separated by commas (`,`) which will be passed to the corresponding executor at runtime.

For instance, `form=access_token,token`, here `form` is the name of the executor, and `access_token,token` is the argument list which will be parsed as `[]string{"access_token", "token"}`.

There are several builtin directive executors, e.g. `form`, `header`, `required`, ... [full list](./directives.go). And you can define your own by:

First, create a **directive executor** by implementing the `httpin.DirectiveExecutor` interface:

```go
func toLower(ctx *DirectiveContext) error {
	if ctx.ValueType.Kind() != reflect.String {
		return errors.New("not a string")
	}

	currentValue := ctx.Value.Elem().String()
	newValue := strings.ToLower(currentValue)
	ctx.Value.Elem().SetString(newValue)
	return nil
}

// Adapt toLower to httpin.DirectiveExecutor.
var MyLowercaseDirectiveExecutor = httpin.DirectiveExecutorFunc(toLower)
```

Seconds, register it to `httpin`:

```go
httpin.RegisterDirectiveExecutor("to_lowercase", MyLowercaseDirectiveExecutor)
```

Finally, you can use your own directives in the struct tags:

```go
type Authorization struct {
	Token string `in:"form=token;required;to_lowercase"`
}
```

The directives will run in the order as they defined in the struct tag.
