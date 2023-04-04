package db

import (
	"context"
	"database/sql"
	"embed"
	"github.com/ethereum/go-ethereum/log"
	_ "github.com/lib/pq" //required for sql library
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS
var migrationLock = 1001

// Migrations struct for running Up/Down
type Migrations struct {
	dsn string
}

type gooseFunc func(db *sql.DB, dir string, opts ...goose.OptionsFunc) error

// NewMigrations constructs Migration
func NewMigrations(dsn string) Migrations {
	return Migrations{dsn: dsn}
}

// Up Migrate the DB to the most recent version available
func (a Migrations) Up(ctx context.Context) error {
	return a.executeFunc(ctx, goose.Up)
}

// Down Rollback one migration
func (a Migrations) Down(ctx context.Context) error {
	return a.executeFunc(ctx, goose.Down)
}

func (a Migrations) executeFunc(ctx context.Context, funcToExecute gooseFunc) error {
	//goose.SetLogger(utils.Logger(ctx))
	goose.SetBaseFS(migrations)

	db, err := sql.Open("postgres", a.dsn)
	if err != nil {
		return err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	log.Info("acquiring lock to run migrations")
	if _, err = db.Exec("select pg_advisory_lock($1)", migrationLock); err != nil {
		return err
	}
	log.Info("migration lock acquired")

	defer func() {
		if _, err = db.Exec("select pg_advisory_unlock($1)", migrationLock); err != nil {
			panic(err)
		}
		log.Info("migration lock released")
		//close connection used for applying migrations
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	log.Info("applying migrations...")
	if err := funcToExecute(db, "migrations"); err != nil {
		return err
	}

	return nil
}
