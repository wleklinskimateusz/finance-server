package auth

import "fmt"

type User struct {
	username       string
	hashedPassword string
}

type UserReposiotry interface {
	FindByUsername(username string) (User, error)
	Save(user User) error
}

type AuthService struct {
	userRepository UserReposiotry
}

func (as AuthService) signup(username string, password string) error {

	hash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	user := User{username: username, hashedPassword: hash}

	err = as.userRepository.Save(user)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	return nil
}
