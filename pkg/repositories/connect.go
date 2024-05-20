package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Jacobbrewer1/vault-provider-examples/pkg/vault"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	EnvDbConnStr = "DB_CONN_STR"
	EnvRedisHost = "REDIS_HOST"
)

// ConnectDB connects to the database
func ConnectDB(v *viper.Viper) (*sqlx.DB, error) {
	connectionString := v.GetString("database.connection_string")
	if connectionString == "" {
		return nil, errors.New("no database connection string provided")
	}

	db, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Test the connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	slog.Info("Connected to database")

	return db, nil
}

func GenerateConnectionStr(v *viper.Viper, vs *vault.Secrets) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=90s&multiStatements=true&parseTime=true",
		vs.Data["username"],
		vs.Data["password"],
		v.GetString("database.host"),
		v.GetString("database.schema"),
	)
}
