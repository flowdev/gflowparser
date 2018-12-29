# Flow Documentation For File: parse_utils.go

## Flow: ParseNameIdent
ParseNameIdent parses a name identifier.
* Regexp: [a-z][a-zA-Z0-9]*
* Semantic result: The parsed text.

![Flow: ParseNameIdent](./ParseNameIdent.svg)

Components | Data
---------- | -----
[TextSemantic](parse_utils.go#L225L228) | [gparselib.ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L105L109)
[gparselib.ParseRegexp](https://github.com/flowdev/gparselib/blob/master/simple_parser.go#L157L179) | 

## Flow: ParsePackageIdent
ParsePackageIdent parses a package identifier.
* Regexp: [a-z][a-z0-9]*\.
* Semantic result: The parsed text (without the dot).

![Flow: ParsePackageIdent](./ParsePackageIdent.svg)

Components | Data
---------- | -----
[TextSemantic](parse_utils.go#L225L228) | [gparselib.ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L105L109)
[gparselib.ParseRegexp](https://github.com/flowdev/gparselib/blob/master/simple_parser.go#L157L179) | 

## Flow: ParseLocalTypeIdent
ParseLocalTypeIdent parses a local (without package) type identifier.
* Regexp: [A-Za-z][a-zA-Z0-9]*
* Semantic result: The parsed text.

![Flow: ParseLocalTypeIdent](./ParseLocalTypeIdent.svg)

Components | Data
---------- | -----
[TextSemantic](parse_utils.go#L225L228) | [gparselib.ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L105L109)
[gparselib.ParseRegexp](https://github.com/flowdev/gparselib/blob/master/simple_parser.go#L157L179) | 

## Flow: ParseOptSpc
ParseOptSpc parses optional space but no newline.
* Semantic result: The parsed text.

![Flow: ParseOptSpc](./ParseOptSpc.svg)

Components | Data
---------- | -----
[ParseASpc](#flow-parseaspc) | [gparselib.ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L105L109)
[TextSemantic](parse_utils.go#L225L228) | 
[gparselib.ParseOptional](https://github.com/flowdev/gparselib/blob/master/complex_parser.go#L84L102) | 

## Flow: ParseASpc
ParseASpc parses space but no newline.
* Semantic result: The parsed text.

![Flow: ParseASpc](./ParseASpc.svg)

Components | Data
---------- | -----
[TextSemantic](parse_utils.go#L225L228) | [gparselib.ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L105L109)
[gparselib.ParseSpace](https://github.com/flowdev/gparselib/blob/master/simple_parser.go#L114L137) | 

## Flow: ParseSpaceComment
ParseSpaceComment parses any amount of space (including newline) and line
(`//` ... <NL>) and block (`/*` ... `*/`) comments.
* Semantic result: The parsed text plus a signal whether a newline was
  parsed.

![Flow: ParseSpaceComment](./ParseSpaceComment.svg)

Components | Data
---------- | -----
[TextSemantic](parse_utils.go#L225L228) | [gparselib.ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L105L109)
pAny | 
pBlkCmnt | 
pLnCmnt | 
pSpc | 
[spaceCommentSemantic](parse_utils.go#L113L118) | 
[gparselib.ParseAny](https://github.com/flowdev/gparselib/blob/master/complex_parser.go#L135L169) | 
[gparselib.ParseBlockComment](https://github.com/flowdev/gparselib/blob/master/simple_parser.go#L223L312) | 
[gparselib.ParseLineComment](https://github.com/flowdev/gparselib/blob/master/simple_parser.go#L184L216) | 
[gparselib.ParseMulti0](https://github.com/flowdev/gparselib/blob/master/complex_parser.go#L62L69) | 
[gparselib.ParseSpace](https://github.com/flowdev/gparselib/blob/master/simple_parser.go#L114L137) | 

## Flow: ParseStatementEnd
ParseStatementEnd parses optional space and comments as defined by
`ParseSpaceComment` followed by a semicolon (`;`) and more optional space
and comments.
The semicolon can be omited if the space or comments contain a new line or
at the end of the input.
* Semantic result: The parsed text.

![Flow: ParseStatementEnd](./ParseStatementEnd.svg)

Components | Data
---------- | -----
BooleanSemantic | [gparselib.ParseData](https://github.com/flowdev/gparselib/blob/master/base.go#L105L109)
[ParseSpaceComment](#flow-parsespacecomment) | 
[TextSemantic](parse_utils.go#L225L228) | 
checkSemicolonOrNewLineOrEOF | 
nil | 
pEOF | 
pOptEOF | 
pOptSemi | 
pSemicolon | 
[gparselib.ParseAll](https://github.com/flowdev/gparselib/blob/master/complex_parser.go#L105L131) | 
[gparselib.ParseEOF](https://github.com/flowdev/gparselib/blob/master/simple_parser.go#L89L109) | 
[gparselib.ParseLiteral](https://github.com/flowdev/gparselib/blob/master/simple_parser.go#L15L35) | 
[gparselib.ParseOptional](https://github.com/flowdev/gparselib/blob/master/complex_parser.go#L84L102) | 

