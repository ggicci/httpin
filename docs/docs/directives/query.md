---
sidebar_position: 1
---

# query

**query** is a [directive](/advanced/concepts#directives) who decodes a field from URL querystring parameters, i.e. [`http.Request.URL.Query()`](https://pkg.go.dev/net/url#URL.Query).

## Definition

```yaml
Executor: "query"
Args: key1 [,key2 [,key3, ...]]
```

httpin will examine values of the keys one by one (key1 -> key2 -> ...), the **first non-empty** value will be used.

## Usage

```go
type ListUsersInput struct {
	IsMember bool  `in:"query=is_member"`
	AgeRange []int `in:"query=age_range[],age_range`
}
```
