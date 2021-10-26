---
sidebar_position: 1
---

# Concepts

## Directive

`Directive` is a formatted string consisting of two parts, the name of the directive, and the arguments for the directive, **separated by an equal sign (`=`)**, formatted as:

```
directive=arg1,arg2,...,argN
```

The left to the `=` is the **name of an directive**. There's a corresponding directive executor (algorithm) working underhood. The right to the `=` is a **list of arguments separated by commas (`,`)** which will be passed to the algorithm at runtime. Arguments can be ommited.

Let's take a look at the following example:

```
type Authorization struct {
	Token string `in:"query=access_token,token;header=x-api-token;required"`
	                  ^----------------------^ ^----------------^ ^------^
	                            d1                    d2            d3
}
```

There are three directives above, **separated by semicolons (`;`)**:

- d1: `query=access_token,token`
- d2: `header=x-api-token`
- d3: `required`

For instance, `query=access_token,token`, here `query` is the name of the directive, and `access_token,token` is the arguments, which will be parsed as `["access_token", "token"]` at runtime.
