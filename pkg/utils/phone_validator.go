package utils

import "regexp"

func CellphoneValidator(cellphone string) bool {
	if cellphone == "" {
		return false
	}

	//regex := regexp.MustCompile(`^(((98)|(\+98)|(0098)|0)(9){1}((0)|(1)|(3)|(9)|(2)){1}[0-9]{8})+$`)
	regex := regexp.MustCompile("^\\+?[0-9]\\d{1,14}$")
	return regex.MatchString(cellphone)
}
