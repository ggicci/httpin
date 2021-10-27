---
sidebar_position: 1
---

# Concepts

**httpin** is driven by defining Go struct tags.

Let's take the following decleartion of a struct as an example to explain how it works:

```go
type Authorization struct {
	Token string `in:"query=access_token,token;header=x-api-token;required"`
	                  ^----------------------^ ^----------------^ ^------^
	                            d1                    d2            d3
}
```

The key of the struct tag to drive **httpin** is `in`. Which is specific for the **httpin** package, like `json` for **encoding/json**.

We can define multiple [directives](#directive) in the tag, which must be **separated by semicolons (`;`)**. See the example above, `d1`, `d2`, `d3` are three different directives.
**httpin** will run the directives in order (`d1` -> `d2` -> `d3`) for each corresponding struct field.

:::caution

Not every directive will be executed by **httpin**. It's decided by the executors (algorithms) of the directives and the actual input (request data).

The execution of a directive can fail, if a former directive failed, the latter ones won't be executed. If `d1` failed, `d2`, `d3` won't run.

:::

## Directive

`Directive` is a formatted string consisting of two parts, the [directive executor](#directive-executor), and the arguments, **separated by an equal sign (`=`)**, formatted as:

```
name=argv
```

Which works like a function call.

The left to the `=` is the name of the directive. There's a corresponding directive executor (with the same name) working underhood.

The right to the `=` is the arguments, which will be passed to the algorithm at runtime. The way to define arguments can differ across different directives. In general, it will be a comma (`,`) separated strings for multiple arguments. Arguments can be ommited. For more specific usage, you should consult the documentation of the directives.

For the above example, there are three directives:

- d1: `query=access_token,token`
- d2: `header=x-api-token`
- d3: `required`

Let's dissect `d1`, i.e. `query=access_token,token`. The **name** is `query`. The **argv** is `access_token,token`. And after reading the documentation of [**query**](/directives/query), we know the args will be treated as `["access_token", "token"]`.

## Directive Executor

`Directive Executor` is an algorithm with runtime context who's responsible to execute a concrete [**directive**](#directive).

To give a better understanding, we can treat `Directive Executor` as a function in a programming lanaguage, and treat `Directive` as a concrete function call.

e.g.

| Diretive Executor | Directive                          | Execution                                  |
| ----------------- | ---------------------------------- | ------------------------------------------ |
| query             | `query=access_token,token`         | `query(["access_token", "token"])`         |
| header            | `header=x-api-token,Authorization` | `header(["x-api-token", "Authorization"])` |
| required          | `required`                         | `required([])`                             |
