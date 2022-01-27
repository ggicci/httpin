# httpin

[![Go](https://github.com/ggicci/httpin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ggicci/httpin/actions/workflows/go.yml) [![documentation](https://github.com/ggicci/httpin/actions/workflows/documentation.yml/badge.svg?branch=documentation)](https://github.com/ggicci/httpin/actions/workflows/documentation.yml) [![codecov](https://codecov.io/gh/ggicci/httpin/branch/main/graph/badge.svg?token=RT61L9ngHj)](https://codecov.io/gh/ggicci/httpin) [![Go Reference](https://pkg.go.dev/badge/github.com/ggicci/httpin.svg)](https://pkg.go.dev/github.com/ggicci/httpin) [![Go Report Card](https://goreportcard.com/badge/github.com/ggicci/httpin)](https://goreportcard.com/report/github.com/ggicci/httpin)

HTTP Input for Go - <b>Decode an HTTP request into a custom struct</b>

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

## Quick View

<table>
<tr>
  <th>Before (use net/http)</th>
  <th>After (use httpin)</th>
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
	Page     int  `in:"query=page"`
	PerPage  int  `in:"query=per_page"`
	IsMember bool `in:"query=is_member"`
}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*ListUsersInput)
	// Do sth.
}
```

  </td>
</tr>
</table>

## Why this package?

| Items                | Before (use net/http package)              | After (use ggicci/httpin package)                                                              |
| -------------------- | ------------------------------------------ | ---------------------------------------------------------------------------------------------- |
| Developer Time       | 😫 Expensive (too much parsing stuff code) | 🚀 **Faster** (define the struct for receiving input data and leave the parsing job to httpin) |
| Code Repetition Rate | 😞 High                                    | **Lower**                                                                                      |
| Code Readability     | 😟 Poor                                    | **Highly readable**                                                                            |
| Maintainability      | 😡 Poor                                    | 😍 **Highly maintainable**                                                                     |
