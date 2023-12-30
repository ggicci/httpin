---
sidebar_position: 1
---

# Concepts

**httpin** is driven by [owl](https://github.com/ggicci/owl) - a framework that drives particular algorithms by defining Go struct tags.

Let's take the following declaration of a struct as an example to explain how it works:

```go
type Authorization struct {
	Token string `in:"query=access_token,token;header=x-api-token;required"`
	                  ^----------------------^ ^----------------^ ^------^
	                            d1                    d2            d3
}
```

The struct tag key for **httpin** is `in`. This key is specific to **httpin**, just in the same way `json` is specifically used by the **encoding/json** package.

We can define multiple [directives](#directive) in the tag, which must be **separated by semicolons (`;`)**. See the example above, where `d1`, `d2` and `d3` are three different directives.
**httpin** will run the directives in order (`d1` -> `d2` -> `d3`) for each corresponding struct field.

:::caution

Not every directive will be executed by **httpin**. It's decided by the executors (algorithms) of the directives and the actual input (request data).

The execution of a directive can fail. If a directive fails, none of the directives listed after it will execute. i.e. If `d1` fails, `d2` and `d3` will not run.

:::

## Directive

`Directive` is a formatted string consisting of two parts, the [directive executor](#directive-executor), and the arguments, **separated by an equal sign (`=`)**, formatted as:

```
name=argv
```

Which works like a function call.

To the left of the `=` is the name of the directive. There's a corresponding directive executor (with the same name) working under the hood.

To the right of the `=` are the arguments, which will be passed to the algorithm at runtime. The way to define arguments can differ across different directives. In general, comma (`,`) separated strings are used for multiple arguments. Arguments can be ommited. For more specific usage, you should consult the documentation of the directives.

For the above example, there are three directives:

- d1: `query=access_token,token`
- d2: `header=x-api-token`
- d3: `required`

Let's dissect `d1`, i.e. `query=access_token,token`.

The **name** is `query`.

The **argv** is `access_token,token`.

After reading the documentation of the [**query**](/directives/query) directive, we know the args will be treated as `["access_token", "token"]`.

## Directive Executor

A `Directive Executor` is an algorithm with runtime context, that's responsible for executing a concrete [**directive**](#directive).

For better understanding, we can think of a `Directive Executor` as a function in a programming language, and a `Directive` as a concrete function call.

e.g.

| Directive Executor | Directive                          | Execution                                  |
| ------------------ | ---------------------------------- | ------------------------------------------ |
| query              | `query=access_token,token`         | `query(["access_token", "token"])`         |
| header             | `header=x-api-token,Authorization` | `header(["x-api-token", "Authorization"])` |
| required           | `required`                         | `required([])`                             |
