---
sidebar_position: 3
---

# gin-gonic/gin 🥤

## Integrations

You have to create a [gin middleware](https://github.com/gin-gonic/gin#using-middleware) on your own.
In the following demo code, `BindInput` is a good example to start.

## Run Demo

```go {15,46,54}
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ggicci/httpin"
	"github.com/gin-gonic/gin"
)

// BindInput instances an httpin engine for an input struct as a gin middleware.
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

type ListUsersInput struct {
	Gender   string `in:"query=gender"`
	AgeRange []int  `in:"query=age_range"`
	IsMember bool   `in:"query=is_member"`
}

func ListUsers(c *gin.Context) {
	input := c.Request.Context().Value(httpin.Input).(*ListUsersInput)
	fmt.Printf("input: %#v\n", input)
}

func main() {
	router := gin.New()

	// Bind input struct with handler.
	router.GET("/users", BindInput(ListUsersInput{}), ListUsers)

	r, _ := http.NewRequest("GET", "/users?gender=male&age_range=18&age_range=24&is_member=1", nil)

	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, r)
}
```

Since it will run timeout on the Go Playground. I removed the `Run` button for the above demo code.
You can test the above snippet by using the following command on your local host:

```bash
mkdir /tmp/test && cd $_

touch main.go
# then COPY & PASTE the above code to main.go


go mod init test
go mod tidy

go run main.go
```

The output on my machine looks like this:

```text
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /users                    --> main.ListUsers (2 handlers)
input: &main.ListUsersInput{Gender:"male", AgeRange:[]int{18, 24}, IsMember:true}
```
