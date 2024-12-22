package infra

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq" // Register pq driver
	"go.uber.org/dig"

	"github.com/rs/zerolog/log"
)

type (
	// Databases Database package implementation
	Databases struct {
		dig.Out
		Pg *sql.DB
	}

	// DatabaseCfgs Database configs
	DatabaseCfgs struct {
		dig.In
		Pg *DatabaseCfg
	}

	// DatabaseCfg DatabaseCfg using envconfig with
	DatabaseCfg struct {
		DBName string `envconfig:"DBNAME" required:"true" default:"dbname"`
		DBUser string `envconfig:"DBUSER" required:"true" default:"dbuser"`
		DBPass string `envconfig:"DBPASS" required:"true" default:"dbpass"`
		Host   string `envconfig:"HOST" required:"true" default:"localhost"`
		Port   string `envconfig:"PORT" required:"true" default:"9999"`

		MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" default:"20" required:"true"`
		MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" default:"5" required:"true"`
		ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"15m" required:"true"`
	}
)

// NewDatabases - creating new Database object / instance
func NewDatabases(cfgs DatabaseCfgs) Databases {
	return Databases{
		Pg: OpenPostgres(cfgs.Pg),
	}
}

// Open connection to postgre
func OpenPostgres(p *DatabaseCfg) *sql.DB {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.DBUser, url.QueryEscape(p.DBPass), p.Host, p.Port, p.DBName,
	)

	// Open a new connection to PostgreSQL using the standard database/sql interface.
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal().Err(err).Msg("postgres: open")
	}

	// Configure the connection pool settings.
	db.SetConnMaxLifetime(p.ConnMaxLifetime)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetMaxOpenConns(p.MaxOpenConns)

	// Test the connection to ensure it is valid.
	if err = db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("postgres: ping")
	}

	return db
}
