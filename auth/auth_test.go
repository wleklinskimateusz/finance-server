package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/wleklinskimateusz/option"
)

type MockUserRepository struct{}

func (MockUserRepository) FindByUsername(username string) (option.Option[UserDB], error) {

	f, err := os.Open(filepath.Join("/tmp/users/", username))
	if err != nil {
		return handleError[UserDB](err, "failed to open a file "+username)
	}
	defer f.Close()

	filestat, err := f.Stat()
	if err != nil {
		return handleError[UserDB](err, "failed to read file details")
	}

	buffor := make([]byte, filestat.Size())
	_, err = f.Read(buffor)
	if err != nil {
		return handleError[UserDB](err, "failed to read from file")
	}
	hashedPassword := string(buffor)

	return option.Some(UserDB{username: username, hashedPassword: hashedPassword}), nil
}

func (m MockUserRepository) Save(user UserDB) error {
	err := os.Mkdir("/tmp/users", 0666)
	if err != nil {
		fmt.Println("Error while creating users temp dir")
	}
	fname := filepath.Join("/tmp/users", user.username)
	err = os.WriteFile(fname, []byte(user.hashedPassword), 0666)
	if err != nil {
		fmt.Println("Error writing file")
	}
	return nil
}

func removeFile(username string) error {
	err := os.Remove("/tmp/users/" + username)
	return err
}

func TestSuccessfulSignUp(t *testing.T) {
	userRepository := MockUserRepository{}
	as := AuthService{userRepository: userRepository}
	err := as.signup("username", "password")
	if err != nil {
		t.Errorf("failed to signup: %v", err)
	}
	err = removeFile("username")
	if err != nil {
		t.Errorf("failed to remove a file: %v", err)
	}
}

func TestLogIn(t *testing.T) {
	userRepository := MockUserRepository{}
	as := AuthService{userRepository: userRepository}
	as.signup("username", "password")

}
