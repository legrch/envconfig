package envconfig

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// MissingEnvError represents a missing environment variable error
type MissingEnvError struct {
	StructField string
	EnvKey      string
}

func (e *MissingEnvError) Error() string {
	return fmt.Sprintf("missing required environment variable %q", e.EnvKey)
}

// Process processes the configuration with detailed error messages.
// It's a drop-in replacement for envconfig.Process that provides more detailed error messages
// for missing environment variables.
func Process(prefix string, spec any) error {
	// First, try standard processing to see if there's an error
	err := envconfig.Process(prefix, spec)
	if err == nil {
		return nil
	}

	// If there's an error about missing required value, enhance the error message
	if strings.Contains(err.Error(), "required key") {
		t := reflect.TypeOf(spec).Elem()
		v := reflect.ValueOf(spec).Elem()

		// Recursively check all fields
		missingVars := findMissingRequiredVars(prefix, t, v)
		if len(missingVars) > 0 {
			var errMsgs []string
			for _, mv := range missingVars {
				errMsgs = append(errMsgs, mv.Error())
			}
			return fmt.Errorf("configuration error:\n\t%s", strings.Join(errMsgs, "\n\t"))
		}
	}

	return err
}

// findMissingRequiredVars recursively checks for missing required environment variables
func findMissingRequiredVars(prefix string, t reflect.Type, v reflect.Value) []*MissingEnvError {
	var missing []*MissingEnvError

	for i := range t.NumField() {
		field := t.Field(i)
		envTag := field.Tag.Get("envconfig")
		reqTag := field.Tag.Get("required")

		if field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) {
			var nextPrefix string
			if envTag != "" {
				nextPrefix = envTag
			}
			if prefix != "" {
				if nextPrefix != "" {
					nextPrefix = prefix + "_" + nextPrefix
				} else {
					nextPrefix = prefix
				}
			}

			var fieldValue reflect.Value
			if field.Type.Kind() == reflect.Ptr {
				if v.Field(i).IsNil() {
					fieldValue = reflect.New(field.Type.Elem()).Elem()
				} else {
					fieldValue = v.Field(i).Elem()
				}
				missing = append(missing, findMissingRequiredVars(nextPrefix, field.Type.Elem(), fieldValue)...)
			} else {
				fieldValue = v.Field(i)
				missing = append(missing, findMissingRequiredVars(nextPrefix, field.Type, fieldValue)...)
			}
			continue
		}

		if reqTag == "true" && envTag != "" {
			envKey := envTag
			if prefix != "" {
				envKey = prefix + "_" + envTag
			}

			if os.Getenv(strings.ToUpper(envKey)) == "" {
				missing = append(missing, &MissingEnvError{
					StructField: field.Name,
					EnvKey:      strings.ToUpper(envKey),
				})
			}
		}
	}

	return missing
}
