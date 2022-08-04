package azkeyvault

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func GetSecretKV(secretName string) string {

	mySecretName := secretName
	keyVaultName := os.Getenv("KEY_VAULT_NAME")
	keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)

	cred, err := azidentity.NewClientSecretCredential(os.Getenv("TENANT_ID"), os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), nil)

	client, err := azsecrets.NewClient(keyVaultUrl, cred, nil)
	if err != nil {
		log.Fatalf("Failed to connect to client: %v", err)
	}

	resp, err := client.GetSecret(context.TODO(), mySecretName, nil)
	if err != nil {
		log.Fatalf("failed to get the secret: %v", err)
	}

	return *resp.Value
}
