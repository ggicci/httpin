---
sidebar_position: 2
---

# Error Handler

Introduced in v0.6.0.

While using `httpin.NewInput` to create an HTTP middleware handler, an error handler will be used to handle cases of decoding failures. You can sepcify a custom error handler for **httpin** to use. Which should adhere to the following signature:

```go
func CustomErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
    // ...
}
```

## The WithErrorHandler Option

Using with an HTTP middleware handler:

```go {5}
router := chi.NewRouter()

func init() {
    router.With(
        httpin.NewInput(ListThingsInput{}, WithErrorHandler(CustomErrorHandler)),
    ).Get("/things/:id", ListThings)
}
```

Using with a core:

```go {1}
co, err := httpin.New(Thing{}, WithErrorHandler(CustomErrorHandler))
input, err := co.Decode(req)
```

## Global Error Handler

If you are using `httpin.NewInput` to create middlewares, you will find that it's annoying to add an option to each call in order to use a custom error handler.

Replace the default error handler globally:

```go {8}
import httpin_core "github.com/ggicci/httpin/core"

func myCustomErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
    // ...
}

func init() {
    httpin_core.RegisterErrorHandler(myCustomErrorHandler)
}
```
