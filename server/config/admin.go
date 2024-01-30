package config

import (
	"context"
	"fmt"
	"encoding/json"

	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitializeFirebase() (*firebase.App, error) {
	serviceAccountKey := map[string]interface{}{
		"type":                        "service_account",
		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 os.Getenv("FIREBASE_PRIVATE_KEY") ,
		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        os.Getenv("FIREBASE_CERT_URL"),
		"universe_domain":             "googleapis.com",
	}


	jsonBytes, err := json.Marshal(serviceAccountKey)
	if err != nil {
		return nil, fmt.Errorf("error marshaling credentials: %v", err)
	}

	opt := option.WithCredentialsJSON([]byte(jsonBytes))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	return app, nil
}
