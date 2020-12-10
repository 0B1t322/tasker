package db

// Databaser have method to connect DB
type Databaser interface {
	DB() string
}