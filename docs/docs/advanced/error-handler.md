---
sidebar_position: 2
---

# Error Handler

Introduced since v0.6.0.

While using `httpin.NewInput` to create an HTTP middleware handler, an error handler will be used to handle cases of decoding failures. You can sepcify a custom error handler for **httpin** to use. Which should adhere to the following signature:

```go
func CustomErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
    // ...
}
```

## Use WithErrorHandler option to specify a custom handler

Create an HTTP middleware handler:

```go {5}
router := chi.NewRouter()

func init() {
    router.With(
        httpin.NewInput(ListThingsInput{}, WithErrorHandler(CustomErrorHandler)),
    ).Get("/things/:id", ListThings)
}
```

Create an engine:

```go {1}
engine, err := httpin.New(Thing{}, WithErrorHandler(CustomErrorHandler))
input, err := engine.Decode(req)
```

## Globally replace the default error handler

If you are using `httpin.NewInput`, you will find that it's annoying to add an option to each call in order to use a custom error handler.

So, `httpin.ReplaceDefaultErrorHandler` was introduced to replace the default error handler globally.
