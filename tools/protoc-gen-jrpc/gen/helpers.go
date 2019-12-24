package gen

import (
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jakewright/home-automation/tools/protoc-gen-jrpc/gen/typemap"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func (g *Generator) goTypeName(def *typemap.MessageDefinition) string {
	var prefix string
	if pkg := goPackageName(def.File); pkg != g.pkg {
		prefix = pkg + "."
	}

	var name string
	for _, parent := range def.Lineage() {
		name += camelCase(parent.Descriptor.GetName()) + "_"
	}
	name += camelCase(def.Descriptor.GetName())
	return prefix + name
}

// goFileName returns the output name for the generated Go file
func goFileName(d *descriptor.FileDescriptorProto) string {
	name := *d.Name
	if ext := path.Ext(name); ext == ".proto" {
		name = name[:len(name)-len(ext)]
	}
	name += ".rpc.go"
	return name
}

func goPackageName(d *descriptor.FileDescriptorProto) string {
	_, pkg, ok := goPackageOption(d)
	if !ok {
		panic(fmt.Sprintf("File %s does not have option go_package set", d.GetName()))
	}

	return pkg
}

// goPackageOption interprets the file's go_package option.
// If there is no go_package, it returns ("", "", false).
// If there's a simple name, it returns ("", pkg, true).
// If the option implies an import path, it returns (impPath, pkg, true).
func goPackageOption(f *descriptor.FileDescriptorProto) (string, string, bool) {
	opt := f.GetOptions().GetGoPackage()
	if opt == "" {
		return "", "", false
	}

	// A semicolon-delimited suffix delimits the import path and package name.
	sc := strings.Index(opt, ";")
	if sc >= 0 {
		return opt[:sc], cleanPackageName(opt[sc+1:]), true
	}

	// The presence of a slash implies there's an import path
	slash := strings.LastIndex(opt, "/")
	if slash < 0 {
		return "", opt, true
	}

	return opt, opt[slash+1:], true
}

func cleanPackageName(name string) string {
	name = strings.Map(badToUnderscore, name)
	// Identifier must not be keyword: insert _.
	if isGoKeyword[name] {
		name = "_" + name
	}
	// Identifier must not begin with digit: insert _.
	if r, _ := utf8.DecodeRuneInString(name); unicode.IsDigit(r) {
		name = "_" + name
	}
	return name
}

// badToUnderscore is the mapping function used to generate Go names from package names,
// which can be dotted in the input .proto file.  It replaces non-identifier characters such as
// dot or dash with underscore.
func badToUnderscore(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return r
	}
	return '_'
}

var isGoKeyword = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"else":        true,
	"defer":       true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

// Given a protobuf name for a Message, return the Go name we will use for that
// type, including its package prefix.
//func (t *twirp) goTypeName(protoName string) string {
//	def := t.reg.MessageDefinition(protoName)
//	if def == nil {
//		gen.Fail("could not find message for", protoName)
//	}
//
//	var prefix string
//	if pkg := t.goPackageName(def.File); pkg != t.genPkgName {
//		prefix = pkg + "."
//	}
//
//	var name string
//	for _, parent := range def.Lineage() {
//		name += stringutils.CamelCase(parent.Descriptor.GetName()) + "_"
//	}
//	name += stringutils.CamelCase(def.Descriptor.GetName())
//	return prefix + name
//}

func stderr(format string, args ...interface{}) {
	format += "\n"
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}

// CamelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
