package main

import (
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserNotExists     struct { error }
type IncorrectPassword struct { error }
type PasswordsNotSame  struct { error }
type NotAnEmail        struct { error }
type UserExists        struct { error }

type User struct {
	Id				int
	Email           string
	Password        string
	IsAuthenticated bool
	IsActive        bool
	IsAdmin         bool
	IsSuperuser     bool
	LastLogin		string
	CreatedAt		string
}

func AuthenticateUser(formData LoginFormData) (User, error) {
	user, err := GetUserByEmail(formData.Email)
	if err != nil { return User{}, UserNotExists{} }

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), append([]byte(formData.Password), secretKey...))
	if err != nil { return User{}, IncorrectPassword{} }

	return user, nil
}

func RegisterUser(formData RegisterFormData) (User, error) {
	re_email := regexp.MustCompile(`[a-zA-Z0-9]+@(?:[a-zA-Z]+\.)+[a-zA-Z]+`)
	if !re_email.MatchString(formData.Email)                { return User{}, NotAnEmail{} }
	if _, err := GetUserByEmail(formData.Email); err == nil { return User{}, UserExists{} }
	if formData.Password != formData.RPassword              { return User{}, PasswordsNotSame{} }

	hashedPassword, err := bcrypt.GenerateFromPassword(append([]byte(formData.Password), secretKey...), bcrypt.DefaultCost)
	if err != nil { return User{}, err }

	user, err := AddUser(formData.Email, string(hashedPassword))
	if err != nil {
		log.Printf("RegisterUser: %d", err)
		return User{}, err
	}

	return user, nil
}