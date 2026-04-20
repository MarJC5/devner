package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLOps struct {
	Host         string
	Port         int
	RootUser     string
	RootPassword string
}

func (m *MySQLOps) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/?multiStatements=false&parseTime=true",
		m.RootUser, m.RootPassword, m.Host, m.Port)
}

func (m *MySQLOps) connect(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("mysql", m.dsn())
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("mysql ping: %w", err)
	}
	return db, nil
}

func (m *MySQLOps) Create(ctx context.Context, dbName, user, password string) error {
	db, err := m.connect(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	// identifiers validated by ValidateName — safe to interpolate
	stmts := []string{
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName),
		fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY ?", user),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'%%'", dbName, user),
		"FLUSH PRIVILEGES",
	}
	for i, s := range stmts {
		var err error
		if i == 1 {
			_, err = db.ExecContext(ctx, s, password)
		} else {
			_, err = db.ExecContext(ctx, s)
		}
		if err != nil {
			return fmt.Errorf("mysql: %s: %w", s, err)
		}
	}
	return nil
}

func (m *MySQLOps) Drop(ctx context.Context, dbName, user string) error {
	db, err := m.connect(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	stmts := []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName),
		fmt.Sprintf("DROP USER IF EXISTS '%s'@'%%'", user),
		"FLUSH PRIVILEGES",
	}
	for _, s := range stmts {
		if _, err := db.ExecContext(ctx, s); err != nil {
			return fmt.Errorf("mysql: %s: %w", s, err)
		}
	}
	return nil
}

func (m *MySQLOps) Exists(ctx context.Context, dbName string) (bool, error) {
	db, err := m.connect(ctx)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var found string
	err = db.QueryRowContext(ctx, "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME=?", dbName).Scan(&found)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
