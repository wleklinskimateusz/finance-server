package auth

import (
	"fmt"

	"github.com/yourusername/wleklinskimateusz/option"
)

type UserDB struct {
	id             option.Option[string]
	username       string
	hashedPassword string
}

type UserReposiotry interface {
	FindByUsername(username string) (option.Option[UserDB], error)
	Save(user UserDB) error
}

type AuthService struct {
	userRepository UserReposiotry
}

func (as AuthService) signup(username string, password string) error {

	hash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	user := UserDB{username: username, hashedPassword: hash}

	err = as.userRepository.Save(user)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	return nil
}

func handleError[T any](err error, message string) (option.Option[T], error) {
	fmt.Println(fmt.Errorf(message+": %v", err))
	return option.None[T](), err
}

func (as AuthService) login(username string, password string) (option.Option[string], error) {
	optionalUser, err := as.userRepository.FindByUsername(username)
	if err != nil {
		return handleError[string](err, "failed to fetch user from the database")
	}
	if optionalUser.IsNone() {
		return handleError[string](fmt.Errorf("user is none"), "no such user exists")
	}
	user := optionalUser.Unwrap()
	result := verifyPassword(password, user.hashedPassword)
	if !result {
		return handleError[string](fmt.Errorf("password does not match"), "")
	}
	if user.id.IsNone() {
		return handleError[string](fmt.Errorf("id is not set"), "")
	}
	return user.id, nil

}
