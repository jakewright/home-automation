package svcdef

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
)

// ErrCircularImport is returned if a circular import exists
var ErrCircularImport = errors.New("circular import")

// Parse returns a structured representation of the def
// file at the given path.
func Parse(filename string) (*File, error) {
	return NewParser(&osFileReader{}).Parse(filename)
}

// Parser parses a def file
type Parser struct {
	fr  FileReader
	lex *lexer
	buf []token
	l   token
	f   *File
	err error
}

// NewParser returns a parser initialised with the FileReader
func NewParser(fr FileReader) *Parser {
	return &Parser{
		fr: fr,
	}
}

// Parse returns a structured representation of the def
// file at the given path. Imports are recursively parsed.
func (p *Parser) Parse(filename string) (file *File, err error) {
	//defer func() {
	//	if e := recover(); e != nil {
	//		switch v := e.(type) {
	//		case error:
	//			err = v
	//		default:
	//			err = fmt.Errorf("%s", e)
	//		}
	//	}
	//}()

	// Read the file
	if p.fr.SeenFile(filename) {
		return nil, fmt.Errorf("%s: %w", filename, ErrCircularImport)
	}
	b, err := p.fr.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", filename, err)
	}

	// Initialise the lexer
	p.lex = newLexer(b)

	// Initialise the file
	p.f = &File{
		Path: filename,
	}

	// Parse tokens until there are none left
	for state := parse; state != nil; {
		state = state(p)
	}

	// For each import
	for alias, imp := range p.f.Imports {
		// Expect import paths to be relative to file we're parsing
		importedFilename := filepath.Join(path.Dir(filename), imp.Path)

		// Recursively parse the imported file
		g, err := NewParser(p.fr).Parse(importedFilename)
		if err != nil {
			if errors.Is(err, ErrCircularImport) {
				return nil, fmt.Errorf("%s -> %w", filename, err)
			}

			return nil, err
		}

		// Update the import with the parsed file
		p.f.Imports[alias].File = g
	}

	// Now that we've parsed all of the tokens
	// from the lexer, we can do a second pass
	// and generate fully-qualified type names.

	byQualifiedName, byQualifiedSlice, err := messagesByQualifiedName(p.f, "")
	if err != nil {
		p.error("failed to get messages by qualified name: %v", err)
	}
	p.f.FlatMessages = byQualifiedSlice

	if p.f.Service != nil {
		for _, r := range p.f.Service.RPCs {
			r.InputType.Qualified, err = qualifyType(r.InputType.Name, "", byQualifiedName)
			if err != nil {
				p.error("failed to qualify %s input type %s: %v", r.Name, r.InputType.Name, err)
			}

			r.OutputType.Qualified, err = qualifyType(r.OutputType.Name, "", byQualifiedName)
			if err != nil {
				p.error("failed to qualify %s output type %s: %v", r.Name, r.OutputType.Name, err)
			}
		}
	}

	if err := qualifyMessageTypes(p.f.Messages, byQualifiedName); err != nil {
		p.error("failed to qualify message field types: %v", err)
	}

	return p.f, nil
}

// next consumes and returns the next token
func (p *Parser) next() token {
	// If there's a token in the buffer, return it
	if len(p.buf) > 0 {
		var n token
		n, p.buf = p.buf[0], p.buf[1:]
		p.l = n
		return n
	}

	p.l = p.lex.nextToken()
	return p.l
}

// nextNonSpace consumes and returns the next non-space token
func (p *Parser) nextNonSpace() token {
	for {
		if t := p.next(); t.typ != tokSpace {
			return t
		}
	}
}

// nextToSpace consumes and returns all of the tokens up to the next space
// note that the space is consumed but not returned
func (p *Parser) nextToSpace() (ts []token) {
	ts = append(ts, p.nextNonSpace())
	for t := p.next(); t.typ != tokSpace; t = p.next() {
		ts = append(ts, t)
	}
	return
}

func (p *Parser) peek() token {
	return p.peekn(1)[0]
}

// peekn returns, but does not consume, the next n tokens
func (p *Parser) peekn(n int) []token {
	if n > len(p.buf) {
		// Fill up the buffer
		l := n - len(p.buf)
		for i := 0; i < l; i++ {
			p.buf = append(p.buf, p.lex.nextNonSpaceToken())
		}
	}

	// Return the first n elements
	return p.buf[0:n]
}

// expect consumes and returns the next token if it matches
// the given token type, otherwise an error is reported.
func (p *Parser) expect(e tokenType) token {
	t := p.nextNonSpace()
	if t.typ != e {
		p.error("unexpected token %s, wanted %s", t, e)
	}
	return t
}

// expectOneOf is like expect but takes a slice of allowed token types
func (p *Parser) expectOneOf(es ...tokenType) token {
	t := p.nextNonSpace()
	for _, e := range es {
		if t.typ == e {
			return t
		}
	}
	p.error("unexpected token: %s", t)
	return token{} // unreachable
}

// expectn expects the given token types in the given order
func (p *Parser) expectn(es ...tokenType) []token {
	ts := make([]token, len(es))
	for i, e := range es {
		ts[i] = p.expect(e)
	}
	return ts
}

// error reports a parsing error
func (p *Parser) error(format string, args ...interface{}) {
	fmt.Printf("Line: %d:%d\n", p.l.line, p.l.pos)
	panic(fmt.Errorf(format, args...))
}

type parseFn func(*Parser) parseFn

// parse looks for the next top level element which can be one of
//   - import statement
//   - service definition
//   - message definition
//   - option e.g. foo = "bar"
// or and end-of-file token.
//
// A valid identifier in the language must begin with a letter,
// and contain only letters, digits, underscores and/or periods.
func parse(p *Parser) parseFn {
	for {
		switch t := p.peekn(1); t[0].typ {
		case tokImport:
			return parseImport
		case tokService:
			return parseService
		case tokMessage:
			p.f.addMessage(parseMessage(p, nil))
		case tokIdentifier:
			key, val := parseOption(p)
			p.f.addOption(key, val)
		case tokComment:
			p.expect(tokComment) // comments are ignored
		case tokEOF:
			return nil
		default:
			p.error("unexpected token: %s", t)
		}
	}
}

// parseImport parses imports which have the form
//   import foo "../service.foo/foo.def"
//          ⬑ an alias must be included
//              ⬑ the path must be a quoted string
func parseImport(p *Parser) parseFn {
	ts := p.expectn(tokImport, tokIdentifier, tokString)
	pth, err := strconv.Unquote(ts[2].val)
	if err != nil {
		p.error("failed to unquote import path: %v", err)
	}
	p.f.addImport(ts[1].val, &Import{
		Alias: ts[1].val,
		Path:  pth,
	})
	return parse
}

// parseService parses a service definition (one per file)
// service Foo {
//         ⬑ the service name must be a valid identifier
func parseService(p *Parser) parseFn {
	p.expect(tokService)

	if p.f.Service != nil {
		p.error("found multiple service definitions")
	}

	p.f.Service = &Service{
		Name: p.expect(tokIdentifier).val,
	}

	p.expect(tokOpenBrace)
	return parseInsideService
}

// parseInsideService looks for service options and RPC definitions
// An option must have the form
//   foo = "bar"
//   ⬑ must be a valid identifier
//         ⬑ the value can be a quoted string, a number, or a boolean (true or false)
func parseInsideService(p *Parser) parseFn {
	for {
		switch t := p.peekn(1); t[0].typ {
		case tokIdentifier:
			key, val := parseOption(p)
			p.f.Service.addOption(key, val)
		case tokRPC:
			return parseRPC
		case tokCloseBrace: // end of the service definition
			p.nextNonSpace()
			return parse
		default:
			p.error("unexpected token inside service: %s", t)
		}
	}
}

// parseRPC looks for an RPC definition of the form
// rpc Foo(FooRequest) FooResponse {
//         ⬑ the request and response types are arbitrary identifiers
//           but would typically be a message defined either in this
//           file or in another (denoted by prefixing the type with
//           the import alias and a dot, e.g. user.Address)
//
// inside the braces, arbitrary options can be defined, e.g.
//   method = "GET"
//   ⬑ must be a valid identifier
//            ⬑ the value can be a quoted string, a number, or a boolean (true or false)
func parseRPC(p *Parser) parseFn {
	ts := p.expectn(tokRPC, tokIdentifier, tokOpenParen, tokIdentifier, tokCloseParen, tokIdentifier, tokOpenBrace)

	inType := ts[3].val
	outType := ts[5].val

	rpc := &RPC{
		Name: ts[1].val,
		InputType: &Type{ // RPC types cannot be optional, repeated or map types
			Name:     inType,
			Original: inType,
		},
		OutputType: &Type{ // RPC types cannot be optional, repeated or map types
			Name:     outType,
			Original: outType,
		},
		// We can't fill in the fully-qualified type names yet
		// because we probably haven't parsed all of the messages
	}

Loop:
	for {
		switch t := p.peekn(1); t[0].typ {
		case tokIdentifier:
			key, val := parseOption(p)
			rpc.addOption(key, val)
		case tokCloseBrace:
			p.nextNonSpace()
			break Loop
		default:
			p.error("unexpected token in RPC: %s", t)
		}
	}

	p.f.Service.RPCs = append(p.f.Service.RPCs, rpc)
	return parseInsideService
}

// parseMessage parses a message definition that begins with
//   message Foo {
//           ⬑ must be a valid identifier
//     foo = "bar"
//         ⬑ options can be defined as in other parts of the file
//     string name
//         ⬑ type names are arbitrary unless they include a period
//     user.Address address
//         ⬑ type names with a period are typically
//           references to a type from an imported file
//     []int numbers
//         ⬑ prefixing a type with [] will mark it as repeated
//     *bool marketing_emails
//         ⬑ prefixing a type with a * will mark it as optional
//     *[]string children
//         ⬑ in this case, it is the list that is optional
//           []* is not valid syntax
//
// It is up to the code generator to decide whether the
// type names are valid (including imported names) and
// to decide what to do with repeated and optional fields.
func parseMessage(p *Parser, parent *Message) *Message {
	ts := p.expectn(tokMessage, tokIdentifier, tokOpenBrace)

	qualifiedMessageName := ts[1].val
	if parent != nil {
		qualifiedMessageName = parent.Name + "." + qualifiedMessageName
	}

	message := &Message{
		Name:          ts[1].val,
		QualifiedName: qualifiedMessageName,
	}

Loop:
	for {
		// If this is the end of the message }
		if p.peek().typ == tokCloseBrace {
			p.nextNonSpace()
			break Loop
		}

		// If this is a comment
		if p.peek().typ == tokComment {
			p.expect(tokComment)
			continue // Ignore it
		}

		// If this is an option
		ps := p.peekn(2)
		if ps[0].typ == tokIdentifier && ps[1].typ == tokAssign {
			key, val := parseOption(p)
			message.addOption(key, val)
			continue
		}

		// If this is a nested message declaration
		if p.peek().typ == tokMessage {
			message.addMessage(parseMessage(p, message))
			continue
		}

		typ := parseType(p)
		f := p.expect(tokIdentifier) // field name

		// If a field option is defined
		var opts map[string]interface{}
		if p.peek().typ == tokOpenParen {
			opts = parseFieldOptions(p)
			if len(opts) == 0 {
				opts = nil // This makes tests cleaner
			}
		}

		if p.peek().typ == tokComment {
			p.expect(tokComment) // Ignore end-of-line comments
		}

		message.Fields = append(message.Fields, &Field{
			Name:    f.val,
			Type:    typ,
			Options: opts,
		})
	}

	return message
}

func parseType(p *Parser) *Type {
	t := p.expectOneOf(tokAsterisk, tokOpenBracket, tokIdentifier)

	var original string

	// If this is an optional field *
	optional := false
	if t.typ == tokAsterisk {
		optional = true
		original += "*"
		t = p.expectOneOf(tokOpenBracket, tokIdentifier)
	}

	// If this is a repeated field []
	repeated := false
	if t.typ == tokOpenBracket {
		repeated = true
		p.expect(tokCloseBracket) // ]
		original += "[]"
		t = p.expectOneOf(tokIdentifier)
	}

	original += t.val
	name := t.val

	// If this is a map type
	mapType := false
	var key, val *Type
	if t.val == "map" {
		mapType = true
		p.expect(tokOpenBracket)
		key = parseType(p)
		p.expect(tokCloseBracket)
		val = parseType(p)

		name = "map"
		original = fmt.Sprintf("%s[%s]%s", original, key.Original, val.Original)
	}

	return &Type{
		Name:     name,
		Original: original,
		// We can't fill in the qualified type yet because
		// we probably haven't parsed all of the messages yet
		Repeated: repeated,
		Optional: optional,
		Map:      mapType,
		MapKey:   key,
		MapValue: val,
	}
}

// parseOption returns the identifier and value from an assignment expression
//   foo = "bar"
//   foo = 500
//   foo = true
func parseOption(p *Parser) (string, interface{}) {
	ts := p.expectn(tokIdentifier, tokAssign)

	switch t := p.expectOneOf(tokBool, tokNumber, tokString); t.typ {
	case tokBool:
		b, err := strconv.ParseBool(t.val)
		if err != nil {
			p.error("failed to parse bool token: %s", t)
		}

		return ts[0].val, b

	case tokNumber:
		// ParseInt with a base of 0 will handle the
		// binary, octal, and hexadecimal cases.
		i, err := strconv.ParseInt(t.val, 0, 64)
		if err == nil {
			return ts[0].val, i
		}

		f, err := strconv.ParseFloat(t.val, 64)
		if err == nil {
			return ts[0].val, f
		}

	case tokString:
		s, err := strconv.Unquote(t.val)
		if err != nil {
			p.error("failed to unquote token %s", t)
		}
		return ts[0].val, s
	}

	// unreachable
	p.error("unexpected token parsing option")
	return "", nil
}

func parseFieldOptions(p *Parser) map[string]interface{} {
	p.expect(tokOpenParen)

	opts := make(map[string]interface{})

	for {
		ts := p.peekn(2)

		// The next token should be an identifier
		if ts[0].typ != tokIdentifier {
			p.error("unexpected token in field options: %s", ts[0])
		}

		// What type of option is it
		switch ts[1].typ {
		case tokComma, tokCloseParen: // short-hand option
			t := p.nextNonSpace() // consume the identifier
			opts[t.val] = true
		case tokAssign: // normal option
			id, val := parseOption(p)
			opts[id] = val
		}

		switch t := p.expectOneOf(tokComma, tokCloseParen); t.typ {
		case tokComma: // expect another option
			continue
		case tokCloseParen: // end of options
			return opts
		}

	}
}
