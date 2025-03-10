package vault

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var vaultToken []byte = []byte(os.Getenv("Vault_pass")) // Fetching the token stored as env variable at os level

func initiateVaultAuth() (context.Context, *vault.Client) {

	ctx := context.Background()

	// prepare a client with the given base address
	client, err := vault.New(
		// vault.WithAddress("http://127.0.0.1:8200"),
		vault.WithAddress("http://host.docker.internal:8200"),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Println(err)
	}

	// authenticate with a root token (insecure)
	if err := client.SetToken(string(vaultToken)); err != nil {
		log.Println(err)
	}

	return ctx, client
}

func ReadSecrets(secretName, path string) map[string]interface{} {

	ctx, client := initiateVaultAuth()
	// read the secret
	s, err := client.Secrets.KvV2Read(ctx, secretName, vault.WithMountPath(path))
	if err != nil {
		log.Println(err)
	}

	return s.Data.Data
}

func WriteSNOWRefreshToken(refreshToken string) {
	ctx, client := initiateVaultAuth()
	currTime := time.Now().Unix()

	// write a secret
	_, err := client.Secrets.KvV2Write(ctx, "SNOW_refresh", schema.KvV2WriteRequest{
		Data: map[string]any{
			"refresh_token":      refreshToken,
			"refresh_epoch_time": currTime,
		}},
		vault.WithMountPath("secret"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New Refresh token successfully updated")
}
