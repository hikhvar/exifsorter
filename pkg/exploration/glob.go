package exploration

import (
	"github.com/gobwas/glob"
	"github.com/pkg/errors"
)

// Matcher is a interface to abstract from gobwas/glob
type Matcher interface {
	Match(string) bool
}

// MatcherFromPatterns returns the matcher according to the given patterns or an error.
// It uses github.com/gobwas/glob to generate Matcher. Thus for the supported syntax have a look there.
func GobwasMatcherFromPatterns(patterns []string) ([]Matcher, error) {
	ret := make([]Matcher, 0, len(patterns))
	for _, p := range patterns {
		matcher, err := glob.Compile(p)
		if err != nil {
			return nil, errors.Wrap(err, "can not instantiate matcher")
		}
		ret = append(ret, matcher)
	}
	return ret, nil
}

func isIgnored(ignores []Matcher, dir string) bool {
	for _, g := range ignores {
		if g.Match(dir) {
			return true
		}
	}
	return false
}
