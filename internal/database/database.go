// Package database manages MySQL and PostgreSQL database/user lifecycle
// for devner projects. The shared stack exposes MySQL on :3306 and Postgres
// on :5432 with root credentials devner/devner.
package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
)

type Engine string

const (
	MySQL    Engine = "mysql"
	Postgres Engine = "postgres"
)

var validName = regexp.MustCompile(`^[a-z][a-z0-9_]{0,62}$`)

// ValidateName rejects anything that could cause SQL injection via identifier.
// Only lowercase letters, digits, and underscores. Must start with a letter.
// Max 63 chars (Postgres limit).
func ValidateName(name string) error {
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid name %q: must match %s", name, validName.String())
	}
	return nil
}

type Credentials struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

// GeneratePassword returns a 32-char hex password.
func GeneratePassword() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

type Manager struct {
	MySQL    *MySQLOps
	Postgres *PostgresOps
}

func NewManager() *Manager {
	return &Manager{
		MySQL:    &MySQLOps{Host: "127.0.0.1", Port: 3306, RootUser: "root", RootPassword: "devner"},
		Postgres: &PostgresOps{Host: "127.0.0.1", Port: 5432, RootUser: "devner", RootPassword: "devner", RootDB: "devner"},
	}
}

func (m *Manager) Create(ctx context.Context, engine Engine, name string) (Credentials, error) {
	if err := ValidateName(name); err != nil {
		return Credentials{}, err
	}
	password := GeneratePassword()
	switch engine {
	case MySQL:
		if err := m.MySQL.Create(ctx, name, name, password); err != nil {
			return Credentials{}, err
		}
		return Credentials{Host: "mysql", Port: 3306, Database: name, User: name, Password: password}, nil
	case Postgres:
		if err := m.Postgres.Create(ctx, name, name, password); err != nil {
			return Credentials{}, err
		}
		return Credentials{Host: "postgres", Port: 5432, Database: name, User: name, Password: password}, nil
	default:
		return Credentials{}, errors.New("unknown engine")
	}
}

func (m *Manager) Drop(ctx context.Context, engine Engine, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	switch engine {
	case MySQL:
		return m.MySQL.Drop(ctx, name, name)
	case Postgres:
		return m.Postgres.Drop(ctx, name, name)
	default:
		return errors.New("unknown engine")
	}
}

func (m *Manager) Exists(ctx context.Context, engine Engine, name string) (bool, error) {
	if err := ValidateName(name); err != nil {
		return false, err
	}
	switch engine {
	case MySQL:
		return m.MySQL.Exists(ctx, name)
	case Postgres:
		return m.Postgres.Exists(ctx, name)
	default:
		return false, errors.New("unknown engine")
	}
}
