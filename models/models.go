package models

import (
	"errors"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           primitive.ObjectID `json:"-" bson:"-"`
	Email        string             `json:"email" bson:"email"`
	EmailVerfied bool               `json:"emailVerified,omitempty" bson:"emailVerfied,omitempt"`
	Name         string             `json:"name" bson:"name"`
	Password     string             `json:"password" bson:"password"`
	Gender       string             `json:"gender" bson:"gender"`
	DateOfBirth  primitive.DateTime `json:"dob" bson:"dob"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
}

func CheckRequiredFeilds(user *User) error {
	if user.Name == "" || user.Password == "" || user.Gender == "" || user.Email == "" {
		return errors.New("required feilds not present")
	}

	// now we will check if email is okay or not
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(user.Email) {
		return errors.New("malformed email")
	}

	if user.Gender != "Male" && user.Gender != "Female" {
		return errors.New("incompatible gender found")
	}

	if len(user.Password) < 8 {
		return errors.New("password insufficient size")
	}

	return nil
}

func HashUserPassword(user *User) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if err != nil {
		return err
	}

	user.Password = string(bytes)

	return nil
}
