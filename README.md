<a href="https://ggicci.github.io/httpin/">
  <img src="https://ggicci.github.io//httpin/img/emoji-dango.png" alt="httpin logo" title="httpin Documentation" align="right" height="60" />
</a>

# httpin - HTTP Input for Go

<div align="center"><h3>HTTP Request from/to Go Struct</h3></div>

<div align="center">

[![Go](https://github.com/ggicci/httpin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ggicci/httpin/actions/workflows/go.yml) [![documentation](https://github.com/ggicci/httpin/actions/workflows/documentation.yml/badge.svg?branch=documentation)](https://github.com/ggicci/httpin/actions/workflows/documentation.yml) [![codecov](https://codecov.io/gh/ggicci/httpin/branch/main/graph/badge.svg?token=RT61L9ngHj)](https://codecov.io/gh/ggicci/httpin) [![Go Report Card](https://goreportcard.com/badge/github.com/ggicci/httpin)](https://goreportcard.com/report/github.com/ggicci/httpin) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) [![Go Reference](https://pkg.go.dev/badge/github.com/ggicci/httpin.svg)](https://pkg.go.dev/github.com/ggicci/httpin)

<table>
  <tr>
    <td align="center">
      <a href="https://ggicci.github.io/httpin/">
        <img src="https://docusaurus.io/img/docusaurus.svg" height="48px" />
      </a>
    </td>
  </tr>
  <tr>
    <td>
      <a href="https://ggicci.github.io/httpin/">Documentation</a>
    </td>
  </tr>
</table>

</div>

## Core Features

**httpin** helps you easily decode data from an HTTP request, including:

- **Query parameters**, e.g. `?name=john&is_member=true`
- **Headers**, e.g. `Authorization: xxx`
- **Form data**, e.g. `username=john&password=******`
- **JSON/XML Body**, e.g. `POST {"name":"john"}`
- **Path variables**, e.g. `/users/{username}`
- **File uploads**

You **only** need to define a struct to receive/bind data from an HTTP request, **without** writing any parsing stuff code by yourself.

Since v0.15.0, httpin also supports creating an HTTP request (`http.Request`) from a Go struct instance.

**httpin** is:

- **well documented**: at https://ggicci.github.io/httpin/
- **open integrated**: with [net/http](https://ggicci.github.io/httpin/integrations/http), [go-chi/chi](https://ggicci.github.io/httpin/integrations/gochi), [gorilla/mux](https://ggicci.github.io/httpin/integrations/gorilla), [gin-gonic/gin](https://ggicci.github.io/httpin/integrations/gin), etc.
- **extensible** (advanced feature): by adding your custom directives. Read [httpin - custom directives](https://ggicci.github.io/httpin/directives/custom) for more details.

## How to decode an HTTP request to Go struct?

```go
type ListUsersInput struct {
	Token    string  `in:"query=access_token;header=x-access-token"`
	Page     int     `in:"query=page;default=1"`
	PerPage  int     `in:"query=per_page;default=20"`
	IsMember bool    `in:"query=is_member"`
	Search   *string `in:"query=search,omitempty"`
}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*ListUsersInput)

	if input.IsMember {
		// Do sth.
	}
	// Do sth.
}
```

## How to encode a Go struct to HTTP request?

```go
type ListUsersInput struct {
	Token    string `in:"query=access_token;header=x-access-token"`
	Page     int    `in:"query=page;default=1"`
	PerPage  int    `in:"query=per_page;default=20"`
	IsMember bool   `in:"query=is_member"`
}

func SDKListUsers() {
	payload := &ListUsersInput{
		Token:    os.Getenv("MY_APP_ACCESS_TOKEN"),
		Page:     2,
		IsMember: true,
	}

	// Easy to remember, http.NewRequest -> httpin.NewRequest
	req, err := httpin.NewRequest("GET", "/users", payload)
	// ...
}
```

## Why this package?

#### Compared with using `net/http` package

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

| Benefits                | Before (use net/http package)              | After (use ggicci/httpin package)                                                              |
| ----------------------- | ------------------------------------------ | ---------------------------------------------------------------------------------------------- |
| ⌛️ Developer Time      | 😫 Expensive (too much parsing stuff code) | 🚀 **Faster** (define the struct for receiving input data and leave the parsing job to httpin) |
| ♻️ Code Repetition Rate | 😞 High                                    | 😍 **Lower**                                                                                   |
| 📖 Code Readability     | 😟 Poor                                    | 🤩 **Highly readable**                                                                         |
| 🔨 Maintainability      | 😡 Poor                                    | 🥰 **Highly maintainable**                                                                     |

## Alternatives and Similars

- [gorilla/schema](https://github.com/gorilla/schema): converts structs to and from form values
- [google/go-querystring](https://github.com/google/go-querystring): encoding structs into URL query parameters
