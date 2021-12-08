---
sidebar_position: 3
---

# Upload Files

Introduced in v0.7.0.

Use [`httpin.File`](https://pkg.go.dev/github.com/ggicci/httpin#File) to retrieve a file uploaded from the request. Make sure it's a [multipart/form-data](https://stackoverflow.com/q/4526273/1592264) request.

```go {4,5}
type UpdateArticleInput struct {
	Title       string        `in:"form=title"`
	IsPrivate   bool          `in:"form=is_private"`
	Cover       httpin.File   `in:"form=cover"`
	Attachments []httpin.File `in:"form=attachments"`
}
```

**NOTE**: you **MUST check** `httpin.File.Valid` before accessing.

## Access the uploaded file

Access filename, filesize and other information from `httpin.File.Header`, which is of type [`multipart.FileHeader`](https://pkg.go.dev/mime/multipart#FileHeader).

```go
func UpdateArticle(rw http.ResponseWriter, r *http.Request) {
    input := r.Context().Value(httpin.Input).(*UpdateArticleInput)

    // User has uploaded a file for the cover.
    if input.Cover.Valid {
        filename := input.Cover.Header.Filename
        filesize := input.Cover.Header.Size

        // Read content.
        fileBytes, err := ioutil.ReadAll(input.Cover)
    }

    // ...
}
```
