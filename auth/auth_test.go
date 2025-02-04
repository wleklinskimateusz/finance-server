package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
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
	userDetails := strings.Split(string(buffor), ",")
	uuid := userDetails[0]
	hashedPassword := userDetails[1]

	return option.Some(UserDB{username: username, hashedPassword: hashedPassword, id: option.Some(uuid)}), nil
}

func (m MockUserRepository) Save(user UserDB) error {
	err := os.MkdirAll("/tmp/users", 0777)
	if err != nil {
		return fmt.Errorf("failed creating /tmp/users: %v", err)
	}
	fname := filepath.Join("/tmp/users", user.username)
	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("failed opening the file: %v", err)
	}
	defer f.Close()

	uuid := uuid.NewString()

	_, err = f.Write([]byte(uuid + "," + user.hashedPassword))
	if err != nil {
		return fmt.Errorf("failed writing to file /tmp/users/{username}: %v", err)
	}
	return nil
}

func removeFile(username string) error {
	err := os.RemoveAll("/tmp/users/" + username)
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

	id, err := as.login("username", "password")
	if err != nil {
		t.Errorf("error while logging in: %v", err)
	}
	if id.IsNone() {
		t.Errorf("id is not set: %v", err)
	}

	removeFile("username")
}
