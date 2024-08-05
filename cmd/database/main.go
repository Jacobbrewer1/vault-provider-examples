package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jacobbrewer1/vault-provider-examples/pkg/logging"
	"github.com/Jacobbrewer1/vault-provider-examples/pkg/repositories"
	"github.com/Jacobbrewer1/vault-provider-examples/pkg/vault"
	"github.com/spf13/viper"
)

var (
	configLocation string
	db             *repositories.Database
)

func main() {
	flag.StringVar(&configLocation, "config", "config.json", "The location of the config file")
	flag.Parse()

	if err := setup(); err != nil {
		slog.Error("Error setting up server", slog.String(logging.KeyError, err.Error()))
		os.Exit(1)
	}

	// Listen for ctrl+c and kill signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		got := <-sig
		slog.Info("Received signal, shutting down", slog.String("signal", got.String()))
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down")
			return
		default:
			// Ping the database to keep the connection alive
			if err := db.PingContext(ctx); err != nil {
				slog.Error("Error pinging database", slog.String(logging.KeyError, err.Error()))
				os.Exit(1)
			}

			slog.Info("Pinged database")

			time.Sleep(5 * time.Second)
		}
	}
}

func setup() (err error) {
	ctx := context.Background()

	v := viper.New()
	v.SetConfigFile(configLocation)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	vaultDb := &repositories.VaultDB{
		Client:         nil,
		Vip:            v,
		Enabled:        false,
		CurrentSecrets: nil,
	}

	// If vault is enabled, create a new vault client and get the secrets
	if v.IsSet("vault") {
		slog.Info("Vault configuration found, attempting to connect")
		vaultDb.Enabled = true

		vc, err := vault.NewClientUserPass(v)
		if err != nil {
			return fmt.Errorf("error creating vault client: %w", err)
		}

		vaultDb.Client = vc

		slog.Debug("Vault client created")

		vs, err := vc.GetSecret(ctx, v.GetString("database.credentials_path"))
		if err != nil {
			return fmt.Errorf("error getting secrets from vault: %w", err)
		}

		slog.Debug("Vault secrets retrieved")
		vaultDb.CurrentSecrets = vs

		dbConnectionString := repositories.GenerateConnectionStr(v, vs)
		v.Set("database.connection_string", dbConnectionString)

		db, err = repositories.ConnectDB(ctx, vaultDb)
		if err != nil {
			return fmt.Errorf("error connecting to database: %w", err)
		}

		slog.Info("Database connection generate from vault secrets")
	} else {
		slog.Warn("Vault configuration not found, using raw values")
		db, err = repositories.ConnectDB(ctx, vaultDb)
		if err != nil {
			return fmt.Errorf("error connecting to database: %w", err)
		}
	}

	return nil
}
