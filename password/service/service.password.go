package service

import (
	"encoding/base64"
	"errors"
	"regexp"

	"github.com/Benyam-S/aseri/common"

	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/password"
	"github.com/Benyam-S/aseri/tools"
	"golang.org/x/crypto/bcrypt"
)

// Service is a type that defines password service
type Service struct {
	passwordRepo password.IPasswordRepository
	commonRepo   common.ICommonRepository
}

// NewPasswordService is a function that returns a new password service
func NewPasswordService(passwordRepository password.IPasswordRepository,
	commonRepository common.ICommonRepository) password.IService {
	return &Service{passwordRepo: passwordRepository, commonRepo: commonRepository}
}

// FindPassword is a method that find and return a user's password that matchs the identifier value
func (service *Service) FindPassword(identifier string) (*entity.Password, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("password not found")
	}

	memberPassword, err := service.passwordRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("password not found")
	}
	return memberPassword, nil
}

// VerifyPassword is a method that verify a user has provided a valid password with a matching verifypassword entry
func (service *Service) VerifyPassword(memberPassword *entity.Password, verifyPassword string) error {
	matchPassword, _ := regexp.MatchString(`^[a-zA-Z0-9\._\-&!?=#]{8}[a-zA-Z0-9\._\-&!?=#]*$`, memberPassword.Password)

	if len(memberPassword.Password) < 8 {
		return errors.New("password should contain at least 8 characters")
	}

	if !matchPassword {
		return errors.New("invalid characters used in password")
	}

	if memberPassword.Password != verifyPassword {
		return errors.New("password does not match")
	}

	memberPassword.Salt = tools.RandomStringGN(30)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(memberPassword.Password+memberPassword.Salt), 12)
	memberPassword.Password = base64.StdEncoding.EncodeToString(hashedPassword)

	return nil
}

// UpdatePassword is a method that updates a certain user password
func (service *Service) UpdatePassword(memberPassword *entity.Password) error {
	err := service.passwordRepo.Update(memberPassword)
	if err != nil {
		return errors.New("unable to update password")
	}
	return nil
}

// DeletePassword is a method that deletes a certain user password
func (service *Service) DeletePassword(identifier string) (*entity.Password, error) {
	password, err := service.passwordRepo.Delete(identifier)
	if err != nil {
		return nil, errors.New("unable to delete password")
	}

	return password, nil
}
