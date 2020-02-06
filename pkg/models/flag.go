// Package models contains representations for different parts of the service
package models

// Flag represents a CLI flag
type Flag struct {
	Type    string
	Name    string
	Short   string
	Default interface{}
	Desc    string
}
