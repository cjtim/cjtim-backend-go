package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
	"go.uber.org/zap"
)

type Vault struct {
	client *api.Client
}

func newVault() (*Vault, error) {
	vaultToken := os.Getenv("VAULT_TOKEN")
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultPath := os.Getenv("VAULT_PATH")
	if vaultToken == "" || vaultPath == "" || vaultAddr == "" {
		return nil, fmt.Errorf("VAULT_TOKEN and VAULT_PATH and VAULT_ADDR is required.")
	}
	client, err := api.NewClient(&api.Config{
		Address: vaultAddr,
	})
	client.SetToken(vaultToken)
	client.Auth().Token().RenewSelf(768 * 3600) // renew 768hr
	return &Vault{
		client: client,
	}, err
}

// readVault - Secret import
func (v *Vault) readVault() error {
	secretPath := os.Getenv("VAULT_PATH")
	secret, err := v.client.Logical().Read(secretPath)
	if err != nil {
		return err
	}
	if secret == nil {
		return nil
	}
	data := secret.Data["data"].(map[string]interface{})
	for k, v := range data {
		if os.Getenv(k) == "" {
			os.Setenv(k, v.(string))
		}
	}
	version, err := secret.Data["metadata"].(map[string]interface{})["version"].(json.Number).Int64()
	if err != nil {
		return err
	}
	secretVersion = version
	return nil
}

func (v *Vault) cronVault() func() {
	return func() {
		zap.L().Info("vault checking new secret...")
		oldVersion := secretVersion
		err := v.readVault()
		if err != nil {
			log.Default().Println("Vault secret error:", err.Error())
			return
		}
		if secretVersion != oldVersion {
			zap.L().Info("vault got new secret. restarting...")
			os.Exit(0) // respawn new pods
		}
		zap.L().Info("vault running latest secret")
	}
}
