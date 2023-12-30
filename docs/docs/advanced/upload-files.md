---
sidebar_position: 3
---

# File Uploads

Introduced in v0.7.0.

Use [`httpin.File`](https://pkg.go.dev/github.com/ggicci/httpin#File) to manipulate files, including uploading files from the client side, and retrieving files from the HTTP request.

**NOTE**: make sure the HTTP request is of [multipart/form-data](https://stackoverflow.com/q/4526273/1592264).

```go {4,5}
type UpdateArticleInput struct {
	Title       string         `in:"form=title"`
	IsPrivate   bool           `in:"form=is_private"`
	Cover       *httpin.File   `in:"form=cover"`
	Attachments []*httpin.File `in:"form=attachments"`
}
```

## Upload Files (Client)

```go
updateArticleRequest := &UpdateArticleInput{
    Title:     "About Me",
    IsPrivate: false,
    Cover:     httpin.UploadFile("/path/to/my/album/travel-selfie-no1.jpg"),
    Attachments: []*httpin.File{
        httpin.UploadFile("/path/to/my/videos/vlog-sunset.mp4"),
        httpin.UploadFile("/path/to/my/videos/vlog-sea.mp4"),
    },
}

req, err := httpin.NewRequest("POST", "/posts/about-me", updateArticleRequest)
```

## Retrieve Files (Server)

`httpin.File` implemented the [`httpin_core.FileHeader`](https://pkg.go.dev/github.com/ggicci/httpin/core#FileHeader) interface, where you can access the filename, filesize, MIME info, as well as the file content.

```go
func UpdateArticle(rw http.ResponseWriter, r *http.Request) {
    input := r.Context().Value(httpin.Input).(*UpdateArticleInput)

    filename := input.Cover.Filename()
    filesize := input.Cover.Size()

    // Read content.
    fileBytes, err := input.Cover.ReadAll()

    // ...
}
```
