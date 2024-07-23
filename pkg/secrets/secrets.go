package secrets

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/api/cloudresourcemanager/v1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"os"
	"strings"
)

type Store interface {
	GetSecret(ctx context.Context, secretID string) (string, error)
	CreateSecret(ctx context.Context, secretID string, secret []byte) error
	DeleteSecret(ctx context.Context, secretID string) error
}

type store struct{}

var Vault Store

func NewSecretManagerStore() *store {
	return &store{}
}

// GetSecret retrieves a secret from GCP Secret Manager
func (s *store) GetSecret(ctx context.Context, secretID string) (string, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to setup client")
	}
	projectNumber, err := getProjectNumber() // TODO store project number in environment variable
	if err != nil {
		return "", err
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%d/secrets/%s/versions/%d", projectNumber, secretID, 1),
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "failed to access secret version")
	}

	if string(result.Payload.Data) == "" {
		return "", errors.New("secret empty")
	}

	return string(result.Payload.Data), nil
}

// CreateSecret will create a secret in GCP Secret Manager
func (s *store) CreateSecret(ctx context.Context, secretID string, secret []byte) error {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to setup client")
	}
	projectNumber, err := getProjectNumber() // TODO store project number in environment variable
	if err != nil {
		return err
	}

	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%d", projectNumber),
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	_, err = client.CreateSecret(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to create secret")
	}

	// after the secret has been created, we still need to create a version
	req2 := &secretmanagerpb.AddSecretVersionRequest{
		Parent: fmt.Sprintf("projects/%d/secrets/%s", projectNumber, secretID),
		Payload: &secretmanagerpb.SecretPayload{
			Data: secret,
		},
	}

	_, err = client.AddSecretVersion(ctx, req2)
	if err != nil {
		return fmt.Errorf("failed to add secret version: %v", err)
	}

	return nil

}

// DeleteSecret will delete a secret in GCP Secret Manager
func (s *store) DeleteSecret(ctx context.Context, secretID string) error {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to setup client")
	}
	projectNumber, err := getProjectNumber() // TODO store project number in environment variable
	if err != nil {
		return err
	}

	req := &secretmanagerpb.DeleteSecretRequest{
		Name: fmt.Sprintf("projects/%d/secrets/%s", projectNumber, secretID),
	}

	if err := client.DeleteSecret(ctx, req); err != nil {
		return errors.Wrap(err, "failed to delete secret")
	}

	return nil

}

// IsNotFoundErr will determine if error is caused because the secret doesn't exist
func IsNotFoundErr(err error) bool {
	return strings.Contains(errors.Cause(err).Error(), "rpc error: code = NotFound")
}

// AlreadyExistsErr will determine if error is caused because the secret already exist
func AlreadyExistsErr(err error) bool {
	return strings.Contains(errors.Cause(err).Error(), "rpc error: code = AlreadyExists")
}

var projectNumberCache int64

func getProjectNumber() (int64, error) {
	if projectNumberCache != 0 {
		return projectNumberCache, nil
	}
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(context.Background())
	if err != nil {
		//log.Fatalf("NewService: %v", err)
		return 0, err
	}

	project, err := cloudresourcemanagerService.Projects.Get(projectID).Do()
	if err != nil {
		//log.Fatalf("Get project: %v", err)
		return 0, err
	}

	projectNumberCache = project.ProjectNumber
	return project.ProjectNumber, nil
}
