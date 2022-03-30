package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
	"go.uber.org/zap"
)

func newVault() (*api.Client, error) {
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
	return client, err
}

// readVault - Secret import
func readVault(client *api.Client) error {
	secretPath := os.Getenv("VAULT_PATH")
	secret, err := client.Logical().Read(secretPath)
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

func cronVault(client *api.Client) func() {
	return func() {
		zap.L().Info("vault checking new secret...")
		oldVersion := secretVersion
		err := readVault(client)
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
