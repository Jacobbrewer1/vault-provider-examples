package vault

import (
	"context"

	vault "github.com/hashicorp/vault/api"
)

var (
	ErrSecretNotFound = vault.ErrSecretNotFound
)

type renewalFunc func() (*vault.Secret, error)

type Secrets struct {
	*vault.Secret
}

type Client interface {
	// SetKvSecretV2 sets a map of secrets at the given path.
	SetKvSecretV2(ctx context.Context, mount, name string, data map[string]any) error

	// GetKvSecretV2 returns a map of secrets for the given path.
	GetKvSecretV2(ctx context.Context, mount, name string) (*vault.KVSecret, error)

	// GetSecret returns a map of secrets for the given path.
	GetSecret(ctx context.Context, path string) (*Secrets, error)

	// RenewLease renews the lease of the given credentials.
	RenewLease(ctx context.Context, name string, credentials *vault.Secret, renewFunc renewalFunc) error
}
