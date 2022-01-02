package tools

import (
	"errors"
	"regexp"

	"github.com/Benyam-S/aseri/entity"
)

// ValidateProfile is a function that validate given entities
// [0] - firstName, [1] - lastName, [2] - phoneNumber, [3] - email
func ValidateProfile(role string, entries ...string) entity.ErrMap {

	var matchFirstName, matchLastName, matchEmail, matchPhoneNumber bool

	if role == entity.RoleAdmin || role == entity.RoleStaff {
		matchEmail, _ = regexp.MatchString(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, entries[3])
	}

	errMap := make(map[string]error)
	matchFirstName, _ = regexp.MatchString(`^[a-zA-Z]\w*$`, entries[0])
	matchLastName, _ = regexp.MatchString(`^\w*$`, entries[1])
	matchPhoneNumber, _ = regexp.MatchString(`^(\+\d{11,12})|(0\d{9})$`, entries[2])

	if !matchFirstName {
		errMap["first_name"] = errors.New("firstname should only contain alpha numerical values and have at least one character")
	}
	if !matchLastName {
		errMap["last_name"] = errors.New("lastname should only contain alpha numerical values")
	}
	if !matchEmail {
		errMap["email"] = errors.New("invalid email address used")
	}

	if !matchPhoneNumber {
		errMap["phone_number"] = errors.New("phonenumber should be +XXXXXXXXXXXX or 0XXXXXXXXX formate, also use url escaping if country code was used")
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// IsValidURL is a function that check whether the given url is valid or not
func IsValidURL(url string) bool {

	pattern, _ := regexp.Compile(`^(?i)\b((?:[a-z][\w-]+:(?:/{1,3}|[a-z0-9%])|
	www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|
	(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:'\".,<>?«»“”‘’]))$`)

	return pattern.MatchString(url)
}
