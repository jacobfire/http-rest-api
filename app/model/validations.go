package model

import validation "github.com/go-ozzo/ozzo-validation"

// requiredIf checks validation of a rule by condition
func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}
