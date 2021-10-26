---
sidebar_position: 3
---

# gin-gonic/gin 🥤

## Create a gin middleware - `BindInput`

About [gin middleware](https://pkg.go.dev/github.com/gin-gonic/gin#section-readme).

```go {2,39,46}
// BindInput instances an httpin engine for a input struct as a gin middleware.
func BindInput(inputStruct interface{}) gin.HandlerFunc {
	engine, err := httpin.New(inputStruct)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		input, err := engine.Decode(c.Request)
		if err != nil {
			var invalidFieldError *httpin.InvalidFieldError
			if errors.As(err, &invalidFieldError) {
				c.AbortWithStatusJSON(http.StatusBadRequest, invalidFieldError)
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(c.Request.Context(), httpin.Input, input)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

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

func ListUsers(c *gin.Context) {
	input := c.Request.Context().Value(httpin.Input).(*ListUsersInput)
	c.JSON(http.StatusOK, input)
}

func main() {
	r := gin.New()
	// Bind the input struct with your API handler.
	r.GET("/users", BindInput(ListUsersInput{}), ListUsers)
	r.Run()
}
```
