package validators

import (
	"net/mail"
)

// func IsValidEmail1(email string) bool {
// 	r, _ := regexp.Compile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

//		return r.MatchString(email)
//	}
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
