package build

import (
	"errors"
	"strings"

	"github.com/Masterminds/semver/v3"
)

var (
	// ErrInvalidRule is returned when the rule is invalid.
	ErrInvalidRule = errors.New("invalid rule")
)

type HuloCompilerRule struct {
	name    string
	version *semver.Version
	raw     string
}

// ParseRule parses the rule and returns the HuloCompilerRule.
//
// The rule is in the format of "name&version".
// If the version is not provided, the rule is parsed as the latest version.
func ParseRule(fullRule string) (hcr *HuloCompilerRule, err error) {
	rules := strings.Split(fullRule, "&")
	if len(rules) == 0 {
		return nil, ErrInvalidRule
	}
	hcr = &HuloCompilerRule{name: rules[0], raw: fullRule}

	if len(rules) > 1 {
		hcr.version, err = semver.NewVersion(rules[1])
		if err != nil {
			return nil, err
		}
	}

	return
}

func (r *HuloCompilerRule) Name() string {
	return r.name
}

func (r *HuloCompilerRule) Version() *semver.Version {
	return r.version
}

func (r *HuloCompilerRule) Raw() string {
	return r.raw
}
