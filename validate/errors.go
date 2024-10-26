package validate

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

var (
	ruleUndefinedRegex = regexp.MustCompile(
		`^Undefined validation function '(?P<tag>[\w\.\-_]+)'`,
	)
	ruleInvalidRegex = regexp.MustCompile(
		`^Bad field type (?P<field_type>.+)`,
	)
)

// Translate (and clean up) validation errors.
func translateErr(key string, err error) error {
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return err // only care about `ValidationErrors`
	}

	translations := []string{}
	for _, translation := range errs.Translate(translator) {
		if strings.HasPrefix(translation, " ") {
			// All the translations start w/ "{0}" (the field name),
			// but validate.Var() doesn't have access to the field name
			// and thus they end up starting w/ an empty space.
			// Add the field name back in.
			translation = fmt.Sprintf("%s %s", key, strings.TrimSpace(translation))
		} else {
			// If the translation is missing, `fe.Translate()` falls back
			// to the default error message - which has the same field name issue
			// (except this time it's a pair of single quotes rather than a leading space).
			translation = strings.ReplaceAll(translation, "''", fmt.Sprintf("'%s'", key))
		}
		translations = append(translations, strings.TrimSpace(translation))
	}

	sort.Strings(translations)
	err = errors.New(strings.Join(translations, ", "))
	return err
}

// Converts a recovered panic value into an error.
func panicToErr(key string, rules string, panicVal any) error {
	str, ok := panicVal.(string)
	if !ok {
		return fmt.Errorf("%v", panicVal)
	}

	// reword to match the other validation errors
	if ruleUndefinedRegex.MatchString(str) {
		matches := ruleUndefinedRegex.FindStringSubmatch(str)
		tag := matches[ruleUndefinedRegex.SubexpIndex("tag")]
		str = fmt.Sprintf("undefined rule [%s: %s]", key, tag)
	}
	if ruleInvalidRegex.MatchString(str) {
		matches := ruleInvalidRegex.FindStringSubmatch(str)
		fieldType := matches[ruleInvalidRegex.SubexpIndex("field_type")]
		str = fmt.Sprintf("invalid rule for %s [%s: %s]", fieldType, key, rules)
	}

	return errors.New(str)
}
