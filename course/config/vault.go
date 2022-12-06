package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

var VaultClient *vault.Client
var VaultTransitKey string

var (
	vaultAddress          string
	vaultRoleIdFilePath   string
	vaultSecretIdFilePath string
	vaultSecretIdWrapped  bool
)

func GetVaultClient() *vault.Client {
	return VaultClient
}

func VaultInit() {
	// get configuration
	vaultAddress = Config.Vault.Address
	vaultRoleIdFilePath = Config.Vault.Auth.RoleIdFilePath
	vaultSecretIdFilePath = Config.Vault.Auth.SecretIdFilePath
	vaultSecretIdWrapped = Config.Vault.Auth.Wrapped
	VaultTransitKey = Config.Vault.Transit.Key

	// create vault client
	config := vault.DefaultConfig()
	config.Address = vaultAddress
	VaultClient, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create vault client, err: %v", err)
	}

	// loop: login and renew token
	// https://github.com/hashicorp/vault-examples
	go func() {
		for {
			vaultLoginResp, err := login(VaultClient)
			if err != nil {
				log.Fatalf("Unable to authenticate to Vault: %v", err)
			}
			tokenErr := renew(VaultClient, vaultLoginResp)
			if tokenErr != nil {
				log.Fatalf("Unable to start managing token lifecycle: %v", tokenErr)
			}
		}
	}()
}

func login(client *vault.Client) (*vault.Secret, error) {
	// read role id from file
	bytes, err := ioutil.ReadFile(vaultRoleIdFilePath)
	if err != nil {
		log.Fatalf("Error reading role ID file: %v", err)
	}
	roleID := strings.TrimSpace(string(bytes))
	if len(roleID) == 0 {
		log.Fatalln("Error: role ID file exists but read empty value")
	}

	// prepare secret id
	secretID := &auth.SecretID{FromFile: vaultSecretIdFilePath}

	// login option: if secret id wrapped
	opts := []auth.LoginOption{}
	if vaultSecretIdWrapped == true {
		opts = append(opts, auth.WithWrappingToken())
	}

	// initialize approle auth method
	appRoleAuth, err := auth.NewAppRoleAuth(roleID, secretID, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize AppRole auth method: %w", err)
	}

	// login to vault AppRole auth method
	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to AppRole auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	return authInfo, nil
}

func renew(client *vault.Client, token *vault.Secret) error {
	//renewable := token.Auth.Renewable
	//if !renewable {
	//	log.Printf("Token is not configured to be renewable. Re-attempting login.")
	//	return nil
	//}

	watcher, err := client.NewLifetimeWatcher(&vault.LifetimeWatcherInput{
		Secret: token,
		//Increment: 3600,
	})
	if err != nil {
		return fmt.Errorf("unable to initialize new lifetime watcher for renewing auth token: %w", err)
	}

	go watcher.Start()
	defer watcher.Stop()

	for {
		select {
		// `DoneCh` will return if renewal fails, or if the remaining lease
		// duration is under a built-in threshold and either renewing is not
		// extending it or renewing is disabled. In any case, the caller
		// needs to attempt to log in again.
		case err := <-watcher.DoneCh():
			if err != nil {
				log.Printf("Failed to renew token: %v. Re-attempting login.", err)
				return nil
			}
			// This occurs once the token has reached max TTL.
			log.Printf("Token can no longer be renewed. Re-attempting login.")
			return nil

		// Successfully completed renewal
		case renewal := <-watcher.RenewCh():
			log.Printf("Successfully renewed: %#v", renewal)
		}
	}
}
