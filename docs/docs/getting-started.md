---
sidebar_position: 0
slug: /
---

# Getting Started

httpin is a **Go** package for **Decoding an HTTP request into a custom struct**. We can decode

- [Query parameters](/directives/query), e.g. `?name=john&is_member=true`
- [Headers](/directives/header), e.g. `Authorization: xxx`
- [Form data](/directives/form), e.g. `username=john&password=******`
- [JSON/XML Body](/directives/body), e.g. `POST {"name":"john"}`
- [Path variables](/directives/path), e.g. `/users/{username}`

and [more](/directives/custom) of an HTTP request into a struct in Go.

## Install

```shell
go get github.com/ggicci/httpin
```

## Quick View

### Before using httpin

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

### Using httpin

```go
type ListUsersInput struct {
	Page     int  `in:"query=page"`
	PerPage  int  `in:"query=per_page"`
	IsMember bool `in:"query=is_member"`
}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*ListUsersInput)
	// Do sth.
}
```

### Comparison

| Items                | Before (use net/http package)              | After (use ggicci/httpin package)                                                              |
| -------------------- | ------------------------------------------ | ---------------------------------------------------------------------------------------------- |
| Developer Time       | 😫 Expensive (too much parsing stuff code) | 🚀 **Faster** (define the struct for receiving input data and leave the parsing job to httpin) |
| Code Repetition Rate | 😞 High                                    | **Lower**                                                                                      |
| Code Readability     | 😟 Poor                                    | **Highly readable**                                                                            |
| Maintainability      | 😡 Poor                                    | 😍 **Highly maintainable**                                                                     |

## ⭕ Example Project

You could visit https://github.com/ggicci/httpin-example/blob/main/main.go for a more detailed example.
