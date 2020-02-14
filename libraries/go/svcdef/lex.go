package svcdef

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// pos represents a byte position in the original input
type pos int

func (p pos) Position() pos {
	return p
}

// token represents a token or text string returned from the scanner
type token struct {
	typ  tokenType // The type of this token
	pos  pos       // The starting position, in bytes, of this token in the input string
	val  string    // The value of this token
	line int       // The line number at the start of this token
}

func (t token) String() string {
	val := t.val

	switch {
	case t == token{}:
		return "zero value token"
	case t.typ == tokEOF:
		val = "EOF"
	case t.typ == tokError:
		return t.val
	case len(t.val) > 10: // truncate long values
		val = fmt.Sprintf("%.10q...", t.val)
	}

	return fmt.Sprintf("%s %q l%d", t.typ, val, t.line)
}

type tokenType string

const (
	tokError        tokenType = "error"         // error occurred; value is the text of the error
	tokSpace        tokenType = "space"         // whitespace
	tokOpenBrace    tokenType = "open_brace"    // {
	tokCloseBrace   tokenType = "close_brace"   // }
	tokOpenParen    tokenType = "open_paren"    // (
	tokCloseParen   tokenType = "close_paren"   // )
	tokOpenBracket  tokenType = "open_bracket"  // [
	tokCloseBracket tokenType = "close_bracket" // ]
	tokAssign       tokenType = "assign"        // =
	tokAsterisk     tokenType = "asterisk"      // *
	tokComma        tokenType = "comma"         // ,
	tokComment      tokenType = "comment"       // a comment in the code
	tokString       tokenType = "string"        // quoted string (includes quotes)
	tokIdentifier   tokenType = "identifier"    // alphanumeric identifier
	tokNumber       tokenType = "number"        // an int, octal, hex or float number
	tokBool         tokenType = "bool"          // a boolean, i.e. true or false
	tokImport       tokenType = "import"        // import keyword
	tokService      tokenType = "service"       // service keyword
	tokRPC          tokenType = "rpc"           // rpc keyword
	tokMessage      tokenType = "message"       // message keyword
	tokEOF          tokenType = "eof"           // end of file
)

var key = map[string]tokenType{
	"import":  tokImport,
	"service": tokService,
	"rpc":     tokRPC,
	"message": tokMessage,
}

var symbol = map[rune]tokenType{
	'{': tokOpenBrace,
	'}': tokCloseBrace,
	'(': tokOpenParen,
	')': tokCloseParen,
	'[': tokOpenBracket,
	']': tokCloseBracket,
	'=': tokAssign,
	'*': tokAsterisk,
	',': tokComma,
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner
type lexer struct {
	input     []byte     // the string being scanned
	pos       pos        // Current position in the input
	start     pos        // Start position of this token
	width     pos        // Width of the last rune read from input
	tokens    chan token // Channel of scanned items
	line      int        // Number of new lines seen
	startLine int        // Start line of this token
}

// next returns the next rune in the input
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRune(l.input[l.pos:])
	l.width = pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

// backup steps back one rune
func (l *lexer) backup() {
	l.pos -= l.width

	// Correct new line count
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// peek returns but does not consume the next rune in the input
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// emit passes an token back to the client
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{
		typ:  t,
		pos:  l.start,
		val:  string(l.input[l.start:l.pos]),
		line: l.startLine,
	}

	l.start = l.pos
	l.startLine = l.line
}

// ignore skips over the pending input before this point
func (l *lexer) ignore() {
	l.line += bytes.Count(l.input[l.start:l.pos], []byte("\n"))
	l.start = l.pos
	l.startLine = l.line
}

// accept consumes the next rune if it's from the valid set
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// accept consumes until the end of the line
func (l *lexer) acceptLine() {
	for {
		if r := l.next(); r == '\n' {
			return
		}
	}
}

// errorf returns an error token and terminates the scan
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{
		typ:  tokError,
		pos:  l.start,
		val:  fmt.Sprintf(format, args...),
		line: l.startLine,
	}
	return nil
}

// nextToken returns the next token from the input
func (l *lexer) nextToken() token {
	return <-l.tokens
}

// nextNonSpaceToken returns the next non-space token from the input
func (l *lexer) nextNonSpaceToken() token {
	for {
		if t := l.nextToken(); t.typ != tokSpace {
			return t
		}
	}
}

// drain drains the output so the lexing goroutine will exit
func (l *lexer) drain() {
	for range l.tokens {
	}
}

func newLexer(input []byte) *lexer {
	l := &lexer{
		input:     input,
		tokens:    make(chan token),
		line:      1,
		startLine: 1,
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer
func (l *lexer) run() {
	for state := lexMain; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func lexMain(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case unicode.IsSpace(r):
			return lexSpace
		case r == '+' || r == '-' || ('0' <= r && r <= '9'):
			l.backup()
			return lexNumber
		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifier
		case symbol[r] != "":
			l.emit(symbol[r])
		case r == '"':
			return lexQuote
		case r == '/':
			if r := l.next(); r == '/' {
				return lexComment
			}
			return l.errorf("found single forward slash")
		case r == eof:
			l.emit(tokEOF)
			return nil
		default:
			return l.errorf("unknown rune %+q", r)
		}
	}
}

// lexSpace will consume a string of whitespace matching any of
//   '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0
func lexSpace(l *lexer) stateFn {
	for unicode.IsSpace(l.next()) {
	}
	l.backup()
	l.emit(tokSpace)
	return lexMain
}

// lexQuote will consume a quoted string, assuming
// the opening quote has already been consumed.
//   "foo bar"
func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(tokString)
	return lexMain
}

// lexComment will consume the remainder of the line,
// assuming we've just seen a "//" pair
func lexComment(l *lexer) stateFn {
	l.acceptLine()
	l.emit(tokComment)
	return lexMain
}

// lexNumber will consume a
func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}

	l.emit(tokNumber)
	return lexMain
}

func (l *lexer) scanNumber() bool {
	// Optional leading sign
	l.accept("+-")

	digits := "0123456789_" // decimal
	if l.accept("0") {
		// Note: Leading 0 does not mean octal in floats.
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_" // hexadecimal
		} else if l.accept("oO") {
			digits = "01234567_" // octal
		} else if l.accept("bB") {
			digits = "01_" // binary
		}
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	if len(digits) == 16+6+1 && l.accept("pP") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

func lexIdentifier(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// Absorb
		default:
			l.backup()
			word := string(l.input[l.start:l.pos])

			switch {
			case key[word] != "":
				l.emit(key[word])
			case word == "true" || word == "false":
				l.emit(tokBool)
			default:
				l.emit(tokIdentifier)
			}

			break Loop
		}
	}
	return lexMain
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
