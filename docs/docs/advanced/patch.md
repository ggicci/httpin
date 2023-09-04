---
sidebar_position: 4
---

# Patch Field

Introduced in v0.12.0.

```go
import "github.com/ggicci/httpin/patch"
```

[`patch.Field`](https://pkg.go.dev/github.com/ggicci/httpin/patch#Field) is a generic struct:

```go
type Field[T any] struct {
	Value T
	Valid bool
}
```

It takes in a type parameter as the type of the actual value it holds, and wraps in a `Valid` field as a sentinel, which is used to **tell "the field is missing" from "the field is empty"**.

Use `patch.Field[T]` as the type of a field.

## JSON Payload Request

```go
import "github.com/ggicci/httpin/patch"

type AccountPatchPayload struct {
	Username patch.Field[string]
	Gender   patch.Field[string]
	Age      patch.Field[int]
}

func PatchAccount(rw http.ResponseWriter, r *http.Request) {
	var payload AccountPatchPayload
	json.NewDecoder(r.Body).Decode(&payload)

	if !payload.Username.Valid {
		// field "Username" is missing (not found or null)
	}
}
```

:::caution

In this package, and in a JSON object, a field is defined as a missing field when:

- if the name/key of the field is not found in the JSON object
- or the name/key of the field is present but its value is null, null is interpreted as having no value

:::

For example, in the following two JSON objects, field Name is missing:

```json
{ "Age": 18 }
{ "Name": null, "Age": 18 }
```

## Form Request

```go
import "github.com/ggicci/httpin/patch"

type AccountPatchForm struct {
	Username patch.Field[string] `in:"form=username"`
	Gender   patch.Field[string] `in:"form=gender"`
	Age      patch.Field[int]    `in:"form=age"`
}

func PatchAccount(rw http.ResponseWriter, r *http.Request) {
    payload := r.Context().Value(httpin.Input).(*AccountPatchForm)

	if !payload.Username.Valid {
		// field "Username" is missing (not found or null)
	}
}
```

:::caution

In this package, in a querystring, form or multipart-form request, a field is defined as a missing field when:

- the name/key of the field doesn't appear in the request

:::

For example, in the following request, `gender` field is empty but not missing (i.e. `Gender.Valid == true`), while `age` field is missing (i.e. `Age.Valid == false`):

```text
GET /tasks?username=ggicci&gender=
```
