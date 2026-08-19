// Package keywords groups the lexer's keywords for editor syntax
// highlighting. Both editor generators consume the same groups, so colouring
// stays consistent across editors and cannot drift from the lexer.
package keywords

import (
	"fmt"

	"github.com/Open-MBEE/OpenSysML/internal/core/lexer"
)

// Group is a set of keywords that highlight alike; generators map Name to
// their own colour vocabulary.
type Group struct {
	Name     string
	Keywords []string
}

// groups only refines the colour a keyword gets; every keyword the lexer
// knows is highlighted regardless (see Partition).
var groups = []Group{
	{
		Name: "declaration",
		Keywords: []string{
			"action", "allocation", "analysis", "assoc", "attribute", "behavior",
			"binding", "calc", "case", "class", "classifier", "comment", "concern",
			"connection", "connector", "constraint", "datatype", "def", "dependency",
			"doc", "enum", "event", "expr", "feature", "flow", "function", "interaction",
			"interface", "item", "language", "metaclass", "metadata", "multiplicity",
			"namespace", "occurrence", "package", "part", "port", "predicate",
			"rendering", "rep", "requirement", "state", "step", "struct", "type",
			"verification", "view", "viewpoint",
		},
	},
	{
		Name: "control",
		Keywords: []string{
			// The state-notation words (`initial`, `done`) are not here: the lexer
			// does not reserve them, so generators read lexer.ContextualWords.
			"accept", "after", "assert", "assign", "at", "decide",
			"do", "else", "entry", "exhibit", "exit", "first",
			"for", "fork", "if", "include", "join",
			"loop", "merge", "parallel", "perform", "render", "require",
			"return", "satisfy", "send", "succession", "terminate", "then",
			"to", "transition", "until", "via", "when", "while",
		},
	},
	{
		Name: "modifier",
		Keywords: []string{
			"abstract", "all", "composite", "const", "constant", "default", "derived",
			"end", "in", "individual", "inout", "library", "member", "new",
			"nonunique", "ordered", "out", "portion", "private", "protected", "public",
			"ref", "snapshot", "standard", "timeslice", "variant", "variation",
		},
	},
	{
		Name: "relationship",
		Keywords: []string{
			"alias", "conjugate", "conjugates", "conjugation", "crosses", "differences",
			"disjoining", "disjoint", "featured", "featuring", "import", "inverse",
			"inverting", "intersects", "redefines", "redefinition", "references",
			"specialization", "specializes", "subclassifier", "subset", "subsets",
			"subtype", "typed", "typing", "unions",
		},
	},
	{
		Name: "operator",
		Keywords: []string{
			"and", "as", "chains", "defined", "hastype", "implies", "istype", "meta",
			"not", "or", "xor",
		},
	},
	{
		Name:     "constant",
		Keywords: []string{"false", "null", "true"},
	},
}

// Partition returns the keyword groups plus the lexer keywords no group
// claims. It fails on a keyword the lexer no longer knows or one claimed by
// two groups.
func Partition() ([]Group, []string, error) {
	return partition(groups)
}

func partition(groups []Group) ([]Group, []string, error) {
	known := map[string]bool{}
	for _, kw := range lexer.Keywords() {
		known[kw] = true
	}
	grouped := map[string]bool{}
	for _, g := range groups {
		for _, kw := range g.Keywords {
			if !known[kw] {
				return nil, nil, fmt.Errorf("keyword group %q names %q, which the lexer no longer knows", g.Name, kw)
			}
			if grouped[kw] {
				return nil, nil, fmt.Errorf("keyword %q appears in more than one group", kw)
			}
			grouped[kw] = true
		}
	}

	var rest []string
	for _, kw := range lexer.Keywords() {
		if !grouped[kw] {
			rest = append(rest, kw)
		}
	}
	return groups, rest, nil
}
