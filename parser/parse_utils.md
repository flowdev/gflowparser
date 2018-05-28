# Utilities for the Flow Parser

This file contains some utilities that help building the flow parser.
Most of them are themself simple parsers.

## ParseSmallIdent [flow]
ParseSmallIdent parses an identifier that starts with a lower case character
(a - z). Potentially followed by more valid identifier characters
(A - Z, a - z or 0 - 9).  The semantic result is the parsed text.

### Flow
     (MainData)-> p:gparselib.ParseRegexp[semantics: TextSemantic] ->

### Details
- [MainData](../data/data.md#maindata)
- [gparselib.ParseRegexp](https://github.com/flowdev/gparselib/blob/master/simpleParser.go#L163)
- [TextSemantic](./parseUtils.md#textsemantic)
