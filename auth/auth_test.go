package auth

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type MockUserRepository struct{}

func (MockUserRepository) FindByUsername(username string) (User, error) {

	f, err := os.Open(filepath.Join("/tmp/users/", username))
	if err != nil {
		fmt.Println(fmt.Errorf("failed to open file: %v", err))
		return nil, nil
	}
	defer f.Close()
	io.Copy()
	return User{username: username}, nil
}

func (m MockUserRepository) Save(user User) error {
	err := os.Mkdir("/tmp/users", 0666)
	if err != nil {
		fmt.Println("Error while creating users temp dir")
	}
	fname := filepath.Join("/tmp/users", user.username)
	data := []string{user.username, user.hashedPassword}
	err = os.WriteFile(fname, []byte(strings.Join(data, ",")), 0666)
	if err != nil {
		fmt.Println("Error writing file")
	}
	return nil
}

func TestSuccessfulSignUp(t *testing.T) {
	userRepository := MockUserRepository{}
	as := AuthService{userRepository: userRepository}
	err := as.signup("username", "password")
	if err != nil {
		t.Errorf("failed to signup: %v", err)
	}
}
