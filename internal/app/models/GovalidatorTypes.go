package models

import "github.com/asaskevich/govalidator"

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
}
