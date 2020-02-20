package svcdef

import (
	"fmt"
	"strings"
)

func messagesByQualifiedName(file *File, prefix string) (map[string]*Message, error) {
	byQualifiedName, err := flattenMessages(file.Messages, prefix)
	if err != nil {
		return nil, err
	}

	for alias, i := range file.Imports {
		imported, err := messagesByQualifiedName(i.File, alias)
		if err != nil {
			return nil, err
		}

		if err := mergeMaps(byQualifiedName, imported); err != nil {
			return nil, err
		}
	}

	return byQualifiedName, nil
}

// flattenMessages returns a map of all fully-qualified
// message names to messages including nested messages
func flattenMessages(messages []*Message, prefix string) (map[string]*Message, error) {
	byQualifiedName := make(map[string]*Message)
	for _, m := range messages {
		if _, ok := byQualifiedName[m.QualifiedName]; ok {
			return nil, fmt.Errorf("duplicate message found: %s", m.QualifiedName)
		}

		m.QualifiedName = prefix + "." + m.Name
		byQualifiedName[m.QualifiedName] = m

		// Recurse for nested messages
		nested, err := flattenMessages(m.Nested, m.QualifiedName)
		if err != nil {
			return nil, err
		}

		// Merge the maps
		if err := mergeMaps(byQualifiedName, nested); err != nil {
			return nil, err
		}
	}

	return byQualifiedName, nil
}

func mergeMaps(m, n map[string]*Message) error {
	for k, v := range n {
		if _, ok := m[k]; ok {
			return fmt.Errorf("duplicate message found: %s", v.QualifiedName)
		}
		m[k] = v
	}
	return nil
}

func qualifyMessageTypes(messages []*Message, byQualifiedName map[string]*Message) error {
	for _, m := range messages {
		for _, f := range m.Fields {
			var err error

			if f.Type.Map {
				f.Type.MapKey.Qualified, err = qualifyType(f.Type.MapKey.Name, m.QualifiedName, byQualifiedName)
				if err != nil {
					return fmt.Errorf("failed to qualify key type %s on field %s in message %s", f.Type.Name, f.Name, m.QualifiedName)
				}

				f.Type.MapValue.Qualified, err = qualifyType(f.Type.MapValue.Name, m.QualifiedName, byQualifiedName)
				if err != nil {
					return fmt.Errorf("failed to qualify value type %s on field %s in message %s", f.Type.Name, f.Name, m.QualifiedName)
				}
			} else {
				f.Type.Qualified, err = qualifyType(f.Type.Name, m.QualifiedName, byQualifiedName)
				if err != nil {
					return fmt.Errorf("failed to qualify type %s on field %s in message %s", f.Type.Name, f.Name, m.QualifiedName)
				}
			}
		}

		if err := qualifyMessageTypes(m.Nested, byQualifiedName); err != nil {
			return fmt.Errorf("failed to qualify nested message field types: %v", err)
		}
	}
	return nil
}

// qualifyType returns the fully-qualified type by looking
//   - for an imported type (top-level from imported file only)
//   - for a scoped type
//   - for a top-level local type
func qualifyType(typ, scope string, messagesByQualifiedName map[string]*Message) (string, error) {
	// If it contains a dot, assume imported type but verify.
	if parts := strings.SplitN(typ, ".", 2); len(parts) == 2 {
		// We should be able to look this up directly
		if _, ok := messagesByQualifiedName[typ]; ok {
			return typ, nil
		}
		return "", fmt.Errorf("failed to find imported message matching %s", typ)
	}

	// If this type has a scope, i.e. it is inside a message
	if scope != "" {
		// Look for a message in the scope
		if _, ok := messagesByQualifiedName[scope+"."+typ]; ok {
			return scope + "." + typ, nil
		}
	}

	// Look for a message defined at the top level
	if _, ok := messagesByQualifiedName["."+typ]; ok {
		return "." + typ, nil
	}

	// This must be a simple type. Return the empty string
	// because it doesn't make sense to have a qualified
	// type for anything but messages.
	return "", nil
}
