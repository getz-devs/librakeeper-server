package services

import (
	"context"

	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var client *auth.Client

func InitializeFirebase(credentialPath string) error {
	opt := option.WithCredentialsFile(credentialPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}
	client, err = app.Auth(context.Background())
	return err
}

func VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return client.VerifyIDToken(ctx, idToken)
}
