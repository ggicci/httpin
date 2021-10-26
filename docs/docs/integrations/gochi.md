---
sidebar_position: 1
---

# go-chi/chi

[**go-chi/chi**](https://github.com/go-chi/chi) is

> a lightweight, idiomatic and composable router for building Go HTTP services.

## Chain `httpin.NewInput` middleware with `chi.With` method

```go {21}
type Pagination struct {
	Page    int `in:"query=page"`
	PerPage int `in:"query=per_page,page_size"`
}

type ListUsersInput struct {
	Gender   string `in:"query=gender"`
	AgeRange []int  `in:"query=age_range"`
	IsMember bool   `in:"query=is_member"`
	Pagination
}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*ListUsersInput)
	json.NewEncoder(rw).Encode(input)
}

func main() {
	router := chi.NewRouter()
	// use `With` method to chain the middleware created by `httpin.NewInput`
	router.With(httpin.NewInput(ListUsersInput{})).Get("/users", ListUsers)
}
```

## path Directive by `URLParam` Method

```go {4}
func init() {
	// Register a directive named "path" to retrieve values from `chi.URLParam`,
	// i.e. decode path variables.
	httpin.UseGochiURLParam("path", chi.URLParam)
}

type GetArticleOfUserInput struct {
	Author    string `in:"path=author"` // equivalent to chi.URLParam("author")
	ArticleID int64  `in:"path=article_id"`
}

func GetArticleOfUser(rw http.ResponseWriter, r *http.Request) {
	var input = r.Context().Value(Input).(*GetArticleOfUserInput)
	// ...
}

func main() {
	r := chi.NewRouter()
	r.With(
		httpin.NewInput(GetArticleOfUserInput{}),
	).Get("/{author}/p/{article_id}", GetArticleOfUser)
	// ...
}
```
