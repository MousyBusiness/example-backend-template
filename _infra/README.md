# Infrastructure

### Environment Setup

1. Install Terraform
2. Install gcloud
3. Install Go 1.22+
4. Create gcloud configuration `gcloud config configurations create example-project-dev`
5. Set gcloud project `gcloud config set project example-project-dev`
6. Authenticate gcloud `gcloud auth login`

### Manual Setup Steps

Creating a new environment e.g. example-project-test or updating the project Services or api keys can be done using the initialize script.
1. `./scripts/initialize`
> Creating a new project via the GCP console avoids a mutual dependency issue with setting up gcloud
2. Connect Github repos to CloudBuild
3. Generate a private access key for Github CI (using a read only Github account @YourGithubUserCI), add to Secret Manager under `github-ci-token`.
4. Enable "Secret Manager Accessor" for the  CloudBuild CI (dev only).
5. Enable "Cloud Run" for the CloudBuild CI.
6. Enable "Compute Engine" for the CloudBuild CI.
7. Add "Compute Load Balancer Admin" to the the CloudBuild service account (to allow rolling of compute instances).
8. Add Firebase admin service account key json to Secret Manager under `firebase-admin-sa` (used for verifying admin / test tokens)

##### Dev deploy
`deploy --stage dev`

##### Prod deploy
`deploy --stage prod`
> You will be asked to type 'yes' after confirming resource changes
