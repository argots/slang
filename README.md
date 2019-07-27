# Slang

Slang is a simple declarative meta-programming language for data and code.

A core principle of slang is allow intent to declared declaratively and to separate representation, validations (such as type checks), performance and other choices from the intent.

## Example

The following example charts a table:

```slang
Chart(series: Series(data), type: 'BarChart')
  .where(
     data: table.group(group(it)).map(size(it)),
     group(x): math '⌊x / bucket ⌋ * bucket',
     bucket: time '1 second',
  )
  .rewrite(
    note 'Chart is actually [ZChartV2_3](zchart.com), so do the renaming here'
    Any / Chart: "ZChartV2_3
  )
```

### Markdown

Slang starts with [Markdown](https://en.wikipedia.org/wiki/Markdown)
as the container for code.  This whole markdown file, for example, is
a valid slang program -- the example above is actually picked up
because it is in a code fence (tagged `slang`)

In addition, slang allows custom DSLs inline.  The `note '....'`
line is an example of inline markdown documentation

### Syntactic simplicity

Slang treats single-quotes, double-quotes and back-quotes the same.
Square brackets, curly brackets and paranthesis are the same.

Slang has very limited syntax: things that would require a special
syntax in other languages are just function calls.  The `where`
function in the example is an example: it defines names used on the
first line.  

Similarly, `if(condition, then, else)` is function.

### Functional

Slang works with immutable data and functions for control structures.
The `where` function allows declaring variables used in the previous
expression (the scope is limited to this).  In addition, even creating
functions uses the simplified syntax (see `group(x): ...`).

Inline functions in slang are common.  `.group` and `.map` both expect
this. This is implemented by a dynamically scoped variable `it`. Any
expression within the function which uses this will automatically be
treated as a function.

Note that the following wont work: `table.map(z).where(z: it.field)`

Slang also allows chaining in cases where the dot notation won't work:
`delta.then(base - it)`.  The `then()` call simply pass its receiver
to the function expression via `it`.

### Inline DSLs

Special syntax can be added to Slang relatively easily (say, to
support JSX).  The example shows the use of Math `floor` syntax (via
`⌊x ⌋`) by invoking the `math` parser.

Any ID followed immediately by a quote is treated as invoking a DSL
with the actual parsing left to the domain extension.

Inline markdowns (`note`) are implemented using this.

The actual parsing of everthying past the quote is to be handled by the `jsx` extension in 
this case which is expected to emit valid slang code (and so can refer names declared in
the where clause)

### Meta programming

Slang provides a built-in mechanism to rewrite parts of the code via
the `rewrite` function (as shown in the example).  All rewrite calls
work on the AST of the code and so these are executed first to create
term tranformations.

The output of rewrite is expected to be a valid slang program.

Similar to `.rewrite`, other functions exist for declaring types or
suggesting use of mutable types for performance etc.

All of these have a standard way of referring to the parts of the
code -- by using a path formed out of function names, arg names,
fields, etc.

### Semantic forking

The ability to rewrite parts of the code allows creating hooks where
consumers can modify that. This allows thinking of `.rewrite()` as a
patch - except, it can be semantically connected unlike git patches
which are connected by line-numbers.
