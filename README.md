# Slang

![Test](https://github.com/argots/slang/workflows/Test/badge.svg)
![Lint](https://github.com/argots/slang/workflows/Lint/badge.svg)

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

## Language details

The language is an expression-based language: there are no statements.

### Literals

The basic literals in the language are strings and numbers.  Strings
can use single quotes, double quotes, back quotes or any Unicode quote
character (though the closing character must match the opening)  and
can all be multi-line. There is no escape sequence available but these
can be provided with functions.

### Identifiers

Identifiers are letters (including Unicode) followed by any letter +
number combinations. Identifiers can include quoted strings if no
space separates the identifier and the quoted string.  This allows
arbitrary characters in identifiers. For example, `x"The vector's
average"` is a valid identifier.

### Operators

| Operators      | Description                                      |
| -------------- | ------------------------------------------------ |
| + - * /        | Standard arithmetic. Minus is also unary prefix. |
| = != < > <= >= | Equality, inequality operators.                  |
| & \|           | Logical operators. `not` is a function           |
| ()             | Grouping.  Not used for functions                |
| [] {}          | Ordered sequences or Unordered sets              |
| :              | Tuple operator                                   |
| ,              | Comma separator for sequences and sets           |


### Sequences, sets and function calls

The meaning of sequences, sets and tuples depend on the context of
their usage.  The context is outside the scope of the AST definition.

Sequences and sets can have an identifier before them: `hello[ a, b,
c]` or `hello{a, b, c}`.  When used in the context of data, this
represents a named collection (with the name `hello`).  When used in
the context of code, this represents a function call.  There is no
explicit function call syntax (i.e using the typical `f(x)` form).

Tuples are a special composite type: `a:b` represents a pair.  `a:b:c`
is allowed. The meaning of tuple is context dependent. Within
collections, they can represent a key for the collection entry.  So,
`map{x: 5, y: 20}` can be used to specify a map data type. Within
function calls, they can represent named parameters: `lineTo{x: 5, y:
10}`.  They can also just represent tuples as such.

Note `map{[1, 2]: 42}` is syntactically valid but again, the meaning
may depend on the context and might even be invalid.
