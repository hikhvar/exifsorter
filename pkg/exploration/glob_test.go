package exploration

import (
	"testing"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGobwasMatcherFromPatterns(t *testing.T) {
	tests := []struct {
		name            string
		patterns        []string
		expectedMatcher []Matcher
		expectedError   error
	}{
		{
			name:            "normal patterns",
			patterns:        []string{"*", "**foo**bar"},
			expectedMatcher: []Matcher{glob.MustCompile("*"), glob.MustCompile("**foo**bar")},
		},
		{
			name:          "wrong pattern",
			patterns:      []string{"[\\&"},
			expectedError: errors.New("can not instantiate matcher: unexpected end of input"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := GobwasMatcherFromPatterns(test.patterns)
			assert.Equal(t, test.expectedMatcher, m)
			if test.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expectedError.Error())
			}
		})
	}
}
