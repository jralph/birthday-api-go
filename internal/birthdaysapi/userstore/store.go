package userstore

import (
	"encoding/json"
	"fmt"
	"time"
)

type DateOfBirth time.Time

type User struct {
	Username    string      `json:"-" validate:"alpha,required"`
	DateOfBirth DateOfBirth `json:"dateOfBirth"`
}

type UserStore interface {
	Put(*User) error
	Get(string) (*User, error)
}

func (t *DateOfBirth) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(*t).Format("2006-01-02"))
	return []byte(stamp), nil
}

func (t *DateOfBirth) UnmarshalJSON(b []byte) error {
	str := ""
	err := json.Unmarshal(b, &str)
	if err != nil {
		return fmt.Errorf("error unmarshling DateOfBirth: %w", err)
	}

	dob, err := time.Parse("2006-01-02", str)
	if err != nil {
		return fmt.Errorf("error parsing DateOfBirth: %w", err)
	}

	*t = DateOfBirth(dob)
	return nil
}
