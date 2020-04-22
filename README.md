# Slang

[![Test](https://github.com/argots/slang/workflows/Test/badge.svg)](https://github.com/argots/slang/actions?query=workflow%3ATest)
[![Lint](https://github.com/argots/slang/workflows/Lint/badge.svg)](https://github.com/argots/slang/actions?query=workflow%3ALint)
[![Go Report Card](https://goreportcard.com/badge/github.com/argots/slang)](https://goreportcard.com/report/github.com/argots/slang)

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

## Packages

Slang comes with a default set of packages to help manipulate slang
code.  In particular, the AST package has JSON helper that converts
ASTs to JSON and back without loss.

Slang also comes with an AST builder that allows creating ASTs easily
and an AST pattern matcher/replacer.  Between these, slang should
allow programmatically working with code a lot easier.

| Package | Descripton |
| ------- | ---------- |
| [ast](https://github.com/argots/slang/tree/master/pkg/ast) | implements a parser and formatter |
| [cast](https://github.com/argots/slang/tree/master/pkg/cast) | create and build AST nodes }
| [mast](https://github.com/argots/slang/tree/master/pkg/mast) | pattern match AST nodes }
| [eval](https://github.com/argots/slang/tree/master/pkg/eval) | interpreter |


## Slang AST

The slang AST parser is a very permissive expression parser which
produces an AST node.  In particular, the parser allows colon and
commas in all contexts even if they don't actively make sense.

### Literals

The basic literals in the language are strings and numbers.  Strings
can use single quotes, double quotes, back quotes or any Unicode quote
character (though the closing character must match the opening)  and
can all be multi-line.  Unlike most languages, strings have only one
escape sequence: a slash followed by a rune is treated as the rune. 

### Identifiers

Identifiers are letters (including Unicode) followed by any letter +
number combinations. Strings that immediately follow an identifer
(with no space in between) are considered part of the identifier
thereby allowing any character as part of the identifier.

### Operators

| Operators      | Description                                      |
| -------------- | ------------------------------------------------ |
| + - * /        | Standard arithmetic. Minus is also unary prefix. |
| = != < > <= >= | Equality, inequality operators.                  |
| & \|           | Logical operators. `not` is a function           |
| ()             | Grouping as well function-like                   |
| [] {}          | Ordered sequences or Unordered sets              |
| :              | Tuple operator                                   |
| .              | Field/property access                            |
| ,              | Comma separator for sequences and sets           |


### Sequences, sets and function calls

The meaning of sequences, sets and tuples depend on the context of
their usage.  The context is outside the scope of the AST definition.

Sequences and sets can have an identifier before them: `hello[ a, b,
c]` or `hello{a, b, c}`.  When used in the context of data, this
represents a named collection (with the name `hello`).  When used in
the context of code, this represents a function call.

Slang syntax also allows function-call like usage: `a.b(5, 2)`.  The
meaning of this would depend on the context as well (and this might
very well be disallowed for a pure data context).

Tuples are a special composite type: `a:b` represents a pair.  `a:b:c`
is allowed. The meaning of tuple is context dependent. Within
collections, they can represent a key for the collection entry.  So,
`map{x: 5, y: 20}` can be used to specify a map data type. Within
function calls, they can represent named parameters: `lineTo{x: 5, y:
10}`.  They can also just represent tuples as such.

Note `map{[1, 2]: 42}` is syntactically valid but again, the meaning
may depend on the context and might even be invalid.
