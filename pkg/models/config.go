// Package models contains representations for different parts of the service
package models

// Config ...
type Config struct {
	Port     int
	Path     string
	Redirect bool

	TLS TLSConfig

	DB DBConfig

	Log LoggerConfig
}

// TLSConfig represents configuration options for TLS/HTTPS
type TLSConfig struct {
	Port    int
	Enabled bool
	Cert    string
	Key     string
}

// DBConfig represents configuration options for the database connection
type DBConfig struct {
	Engine string
	User   string
	Pass   string
	Host   string
	Port   int
	Name   string
	Conn   string
}

// LoggerConfig respresents configuration options for logrus
type LoggerConfig struct {
	Level string
	File  string
}
