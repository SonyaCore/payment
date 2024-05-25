package middleware

import (
	"net/http"
	"payment/pkg/errors"
	"payment/pkg/utils"
)

// PhoneValidatorMiddleware is the middleware function to validate phone numbers
func PhoneValidatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		phone := r.URL.Query().Get("phone")
		if !utils.CellphoneValidator(phone) {
			errors.Error(w, http.StatusBadRequest, "Invalid phone number")
			return
		}
		next.ServeHTTP(w, r)
	})
}
