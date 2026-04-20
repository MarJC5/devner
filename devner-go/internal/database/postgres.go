package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PostgresOps struct {
	Host         string
	Port         int
	RootUser     string
	RootPassword string
	RootDB       string
}

func (p *PostgresOps) connString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.RootUser, p.RootPassword, p.Host, p.Port, p.RootDB)
}

func (p *PostgresOps) connect(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, p.connString())
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}
	return conn, nil
}

func (p *PostgresOps) Create(ctx context.Context, dbName, user, password string) error {
	conn, err := p.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// CREATE DATABASE cannot run inside a transaction block.
	// pgx uses auto-commit by default when not in a tx, so direct Exec works.
	var exists bool
	if err := conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname=$1)", user).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		// password is safely escaped by pgx via quoted literal; identifiers validated upstream
		_, err := conn.Exec(ctx, fmt.Sprintf(`CREATE USER %q WITH PASSWORD '%s'`, user, escapeLiteral(password)))
		if err != nil {
			return fmt.Errorf("postgres create user: %w", err)
		}
	}

	if err := conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", dbName).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		_, err := conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %q OWNER %q`, dbName, user))
		if err != nil {
			return fmt.Errorf("postgres create db: %w", err)
		}
	}

	_, err = conn.Exec(ctx, fmt.Sprintf(`GRANT ALL PRIVILEGES ON DATABASE %q TO %q`, dbName, user))
	if err != nil {
		return fmt.Errorf("postgres grant: %w", err)
	}
	return nil
}

func (p *PostgresOps) Drop(ctx context.Context, dbName, user string) error {
	conn, err := p.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// Terminate active connections to target DB first.
	_, _ = conn.Exec(ctx, `SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname=$1 AND pid<>pg_backend_pid()`, dbName)

	if _, err := conn.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS %q`, dbName)); err != nil {
		return fmt.Errorf("postgres drop db: %w", err)
	}
	if _, err := conn.Exec(ctx, fmt.Sprintf(`DROP USER IF EXISTS %q`, user)); err != nil {
		return fmt.Errorf("postgres drop user: %w", err)
	}
	return nil
}

func (p *PostgresOps) Exists(ctx context.Context, dbName string) (bool, error) {
	conn, err := p.connect(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)
	var exists bool
	if err := conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", dbName).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// escapeLiteral escapes single quotes for embedding in SQL string literals.
// Used only for passwords which are generated (hex chars) so this is defensive.
func escapeLiteral(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\'' {
			out = append(out, '\'', '\'')
		} else {
			out = append(out, s[i])
		}
	}
	return string(out)
}

