package models

import (
	"github.com/asaskevich/govalidator"
	"regexp"
)

func init() {
	govalidator.CustomTypeTagMap.Set(
		"stringArray",
		govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
			subject, ok := i.([]string)
			if !ok {
				return false
			}

			if len(subject) == 0 {
				return false
			}

			return true
		}),
	)

	govalidator.CustomTypeTagMap.Set(
		"phoneNumber",
		govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
			subject, ok := i.(string)
			if !ok {
				return false
			}

			if len(subject) != 12 {
				return false
			}

			const reg = `^((\+?[0-9]{1,3})|(\+?\([0-9]{1,3}\)))[\s-]?(?:\(0?[0-9]{1,5}\)|[0-9]{1,5})[-\s]?[0-9][\d\s-]{5,7}\s?(?:x[\d-]{0,4})?$`
			ok, _ = regexp.MatchString(reg, subject)

			return ok
		}),
	)

	govalidator.CustomTypeTagMap.Set(
		"password",
		govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
			subject, ok := i.(string)
			if !ok {
				return false
			}

			if len(subject) < 6 || len(subject) > 30 {
				return false
			}

			return ok
		}),
	)
}
