// Package models contains representations for different parts of the service
package models

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

const dateISO = "2006-01-02"

// Person represents a database table containing information about a person
type Person struct {
	gorm.Model
	Name      string    `gorm:"UNIQUE_INDEX"`
	BirthDate BirthDate `json:"dateOfBirth"`
}

// BirthDate is a type alias to allow for custom (un)marshalling
type BirthDate time.Time

// UnmarshalJSON is a custom unmarshaller
func (bd *BirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(dateISO, s)
	if err != nil {
		return err
	}
	*bd = BirthDate(t)
	return nil
}

// MarshalJSON is a custom marshaller
func (bd BirthDate) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(bd).Format(dateISO) + "\""), nil
}

// Format is a custom formatter
func (bd BirthDate) Format(s string) string {
	t := time.Time(bd)
	return t.Format(s)
}

func (bd BirthDate) String() string {
	return bd.Format(dateISO)
}

// Scan is a custom SQL scanner
func (bd *BirthDate) Scan(v interface{}) error {
	t, err := time.Parse(dateISO, strings.Split(v.(string), " ")[0])
	if err != nil {
		return err
	}

	*bd = BirthDate(t)
	return nil
}

// Value is a custom SQL valuer
func (bd BirthDate) Value() (driver.Value, error) {
	return time.Time(bd), nil
}
