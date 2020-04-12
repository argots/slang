# Slang

Slang is a simple declarative meta-programming language for data and code.

## Goals

- A clean human-readable simple syntax that can represent data as well
as code. 
- Type checks and performance optimizations are layered on top of the
fixed syntax instead of requiring additional syntax.
- A canonical format (so all code/data can be formatted cleanly).
- Ability to transform cleanly from/to JSON.
- Ability to represent patches of the syntax within the syntax itself
(meta programming).

## Specifics

The language is an expresion-based language: there are no statements.

The **literals** in the langauge are strings and numbers.  Strings can
use single quotes, double quotes, backquotes or any unicode quote
character  and can all be multiline.

Identifiers are letters (including unicode) followed by any letter +
number combinations. Identifiers can include quoted strings if no
space separates the identifier and the quoted string.  This allows
arbitrary characters in identifiers.

Expressions can use standard binary arithmetic operators: `+, -, *,
/`.  Minus can also be used as a unary prefix operator.

Logical operations are expressed with `&` and `|`.  The unary prefix
operator is absent and a function `not` is used instead.

Inequality and equality are expressed with `<, >, =, <=, >=, !=`

Expressions can be grouped with paranetheses `()`.

The standard set of collections can be ordered (sequences) or
unordered (sets).  Ordered collections use `[a, b, c]` syntax while
unordered collections use the `{a, b, c}` syntax.

Collections can have an identifier before them: `hello[ a, b, c]` or
`hello{a, b, c}`.  When used in the context of data, this represents a
named collection (with the name `hello`).  When used in the context of
code, this represents a function call.  There is no explicit function
calls using the `f(x)` syntax.

Tuples are a special composite type: `a:b` represents a pair.  `a:b:c`
is allowed. The meaning of tuple is context dependent. Within
collections, they can represent a key for the collection entry.  So,
`map{x: 5, y: 20}` can be used to specify a map data type. Within
function calls, they can represent named parameters: `lineTo{x: 5, y:
10}`.  They can also just represent tuples as such.

Note `map{[1, 2]: 42}` is syntactically valid but may be invalid
depending on the context.

