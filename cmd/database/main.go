package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Jacobbrewer1/vault-provider-examples/pkg/logging"
	"github.com/Jacobbrewer1/vault-provider-examples/pkg/repositories"
	"github.com/Jacobbrewer1/vault-provider-examples/pkg/vault"
	vault2 "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

var (
	configLocation string
)

func main() {
	flag.StringVar(&configLocation, "config", "config.json", "The location of the config file")
	flag.Parse()

	if err := setup(); err != nil {
		slog.Error("Error setting up server", slog.String(logging.KeyError, err.Error()))
		os.Exit(1)
	}

	// Just wait and logs will appear as vault renews and reconnects to the database.
	<-make(chan any)
}

func setup() error {
	ctx := context.Background()

	v := viper.New()
	v.SetConfigFile(configLocation)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Check if the vault category is in the config file.
	var vaultClient vault.Client
	dbSecrets := new(vault.Secrets)
	if v.IsSet("vault") {
		slog.Info("Vault configuration found, attempting to connect")

		vc, err := vault.NewClient(v)
		if err != nil {
			return fmt.Errorf("error creating vault client: %w", err)
		}
		vaultClient = vc

		slog.Debug("Vault client created")

		vs, err := vc.GetSecrets(v.GetString("database.credentials_path"))
		if err != nil {
			return fmt.Errorf("error getting secrets from vault: %w", err)
		}
		dbSecrets = vs

		slog.Debug("Vault secrets retrieved")

		dbConnectionString := repositories.GenerateConnectionStr(v, vs)
		v.Set("database.connection_string", dbConnectionString)

		slog.Info("Database connection generate from vault secrets")
	} else {
		slog.Warn("Vault configuration not found, using raw values")
	}

	sqlxDb, err := repositories.ConnectDB(v)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	db := repositories.NewDatabase(sqlxDb)

	if v.IsSet("vault") {
		go func() {
			err := vaultClient.RenewLease(ctx, v.GetString("database.credentials_path"), dbSecrets.Secret, func() (*vault2.Secret, error) {
				slog.Warn("Vault lease expired, reconnecting to database")

				vs, err := vaultClient.GetSecrets(v.GetString("database.credentials_path"))
				if err != nil {
					return nil, fmt.Errorf("error getting secrets from vault: %w", err)
				}

				dbConnectionString := repositories.GenerateConnectionStr(v, vs)
				v.Set("database.connection_string", dbConnectionString)

				newDb, err := repositories.ConnectDB(v)
				if err != nil {
					return nil, fmt.Errorf("error connecting to database: %w", err)
				}

				if err := db.Reconnect(ctx, newDb); err != nil {
					return nil, fmt.Errorf("error reconnecting to database: %w", err)
				}

				slog.Info("Database reconnected")

				return vs.Secret, nil
			})
			if err != nil {
				slog.Error("Error renewing vault lease", slog.String(logging.KeyError, err.Error()))
				os.Exit(1) // Forces new credentials to be fetched
			}
		}()
	}

	return nil
}
