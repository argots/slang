# Slang

Slang is a very simple language:

## Value types

Slang has three basic value types: **Int64**, **String** and **Float64**:  `19`, `"Hello"` and `1.5` respectively.

Floating point numbers require a leading digit: `.5` is not  allowed.

Strings support using slash to escape the next character but sequences like `\r` or `\n` are not interpreted.  Slash is mainly present for escaping the double-quote: `"... named \"Goo\" ..."`

## Identifiers

Slang allows identifiers to be any unicode letter followed by any number of non-whitespace non-special charaters.  The only special character defined at this point is `(` though `{` is also in contention

A unicode letter followed by a sequence of non-whitespace characters followed by a quote is treated specially as an `dialect` which is described later.

The value of identifiers is defined by lexical scope with the exception of `it` which is dynamically scoped.

## Functions

A function is an expression of the form `hello(1, "boo")`.  

An alternate form of functions could be `hello{ count: 1, val: "boo"}` which is syntactic sugar for
`hello(Pair("count", 1), Pair("val", "boo"))` but this is very likely to end up in an extension instead.

## Dialects

Dialects (such as embedding XML) can be invoked like so: `Something( html"<div>42</div>" )`.  Once the quote is seen, slang offloads the parsing of the rest to the underlying extension (which would control the syntax and meaning of the embedded text). The extension defines how the block ends though this is typically with a double-quote.

## Operators

It would be useful to define the Dot operators and simple binary, comparison, logical and arithmetic operators though all of these can be invoked via special functions: `Dot(x, "y")` instead of `x.y`.  

At this point, the  plan is to define operators in a `plus` dialect.

## Scoping

Creating a scope is via the Scope function:

```
Scope(
   Pair("x", ....),
   Pair("y", ....)
)
```

The `plus` dialect will make this more readable:

```
  Scope{x = 42, y = x + 3}.y
```

## Markdown

Slang works with mardown files: it simply concatenates all fenced code blocks and treating that as the code.
