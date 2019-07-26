# Slang

Slang is a simple declarative meta-programming language for data and code.

A core principle of slang is allow intent to declared declaratively and to separate representation, validations (such as type checks), performance and other choices from the intent.

## Readability first

### Markdown

Slang starts with [Markdown](https://en.wikipedia.org/wiki/Markdown) as the container.  Actual code snippets are embedded with code-fences:

```
  This is a *slang*  hello world program.
  
  ## Hello World
  ```slang
     'hello world'
  `` `
  
```

### Syntactic simplicity

Slang treats single-quotes, double-quotes and back-quotes the same.  Square brackets, curly brackets and paranthesis are the same.

Slang has very limited syntax: things that would require a special syntax in other languages  are just function calls in slang:

Consider the following JS snippet:

```js
function averagePositive(items) {
  let result = 0;
  let count = 0;
  for(let val in items) {
    if (val > 0) {
      result += val;
      count ++;
    }
  }
  return result/count;
}

```

The equivalent version in slang would be like so:

```
Object(
  AveragePositive(items): div(reduce(filter(items, it > 0), initial, sum)),
  where(
     initial = Pair(Total: 0, Count: 0),
     div(pair) = pair.Total/pair.Count,
     sum(it) = Pair(Total: it.last.Total + it.current, Count: it.last.Count + 1)
  )
)
```

The `where` function is the way to declare local scopes (which can be declared before or 
after the actual term is used).  Similarly, `if()` is a function as well rather than a 
syntactic term.

This choice may seem counter to better readability through keywords and special syntactic 
forms but having a large syntax has a high onboarding cost as well as a higher cost for 
parsers and readability tools.


### Inline DSLs

Slang offers a special form for embedding XML or Math or even documentation: any identifier 
followed by a quoted string (without an operator inbetween) would be considered a language
DSL and the corresponding extension invoked:

```
  do(
    RootComponent = jsx "
       <MyReactJSXCode>
        ....
       </MyReactJSXCode>
    ",
    where(
       React = import("react"),
       ...
    )
  )
```

The actual parsing of evertying past the quote is to be handled by the `jsx` extension in 
this case which is expected to emit valid slang code (and so can refer names declared in
the where clause)


In 
Slang is a very simple language. The core language can be thought of a data serialization format or a config or an actual program.

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

Creating a scope is via the `do` function:

```
do(
   Pair("x", ....),
   Pair("y", ....)
)
```

The `plus` dialect will make this more readable:

```
  do{x = 42, y = x + 3}.y
```

## Markdown

Slang works with mardown files: it simply concatenates all fenced code blocks and treating that as the code.

## Plus dialect

The plus dialect adds a few features:

1. Curly braces for function calls with named parameters: `f{x = 23, y = 24}`
2. Standard binary operators: `a + 2` or `x.y` 
3. Scopes via `do{ x  = 23, y = x + 22 }.y`
4. `that` behavior: `(x + 1).(that * log(that))` is equivalent to `do{y0 = x + 1, y1 = y0 * log(y0)}.y1`
   - this allows long chained expressions even in cases where the value is not an object 
