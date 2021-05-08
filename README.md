# httpin

[![codecov](https://codecov.io/gh/ggicci/httpin/branch/main/graph/badge.svg?token=RT61L9ngHj)](https://codecov.io/gh/ggicci/httpin)

HTTP Input for Go - Decode an HTTP request into a custom struct

**Define your input struct and then fetch your data!**

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
		// ... invalid page
		return
	}

	perPage, err := strconv.ParseInt(r.FormValue("per_page"), 10, 64)
	if err != nil {
		// ... invalid per_page
		return
	}

	isVip, _ := strconv.ParseBool(r.FormValue("is_vip"))

	// do sth.
}
```

  </td>
  <td>

```go
type ListUsersInput struct {
	Page    int  `in:"form=page"`
	PerPage int  `in:"form=per_page"`
	IsVip   bool `in:"form=is_vip"`
}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	interfaceInput, err := httpin.New(ListUsersInput{}).Decode(r)
	if err != nil {
		// err can be *httpin.InvalidField
		return
	}

	input := interfaceInput.(*ListUsersInput)
	// do sth.
}
```

  </td>
</tr>
</table>

## Features

- [x] Builtin directive `form` to decode a field from HTTP query, i.e. `http.Request.Form`
- [x] Builtin directive `header` to decode a field from HTTP headers, e.g. `http.Request.Header`
- [x] Builtin decoders used by `form` and `header` directives for basic types, e.g. `bool`, `int`, `int64`, `float32`, `time.Time`, ... [full list](./decoders.go)
- [x] Decode a field by inspecting a set of keys from the same source
- [x] Decode a field from multiple sources, e.g. both query and headers
- [ ] Register or replace decoders for both builtin basic types and custom types
- [x] Define input struct with embedded struct fields
- [x] Builtin directive `required` to tag a field as **required**
- [ ] Builtin encoders for basic types
- [ ] Register or replace encoders for both builtin basic types and custom types
- [x] Register custom directive executors to extend the field resolving abilities, see directive [required](./required.go) as an example and think about implementing your own directives like `trim`, `to_lowercase`, `base58_to_int`, etc.

## Sample User Defined Input Structs

```go
type Authorization struct {
	// Decode from multiple sources, the former with higher priority
	Token string `in:"form=access_token,header=x-api-token"`
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

## Advanced - Use Middleware

First, set up the middleware for your handlers. We recommend using [alice](https://github.com/justinas/alice) to chain your HTTP middleware functions.

```go
func init() {
	mux.Handle("/users", alice.New(
		httpin.NewInput(ListUsersInput{}),
	).ThenFunc(ListUsers)).Methods("GET")
}
```

Second, fetch your input with **only one line** of code.

```go
func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*UserQuery)
	// do sth.
}
```
