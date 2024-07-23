# 1. gcloud config set project example-project-dev
# 2. gcloud auth application-default login
# 3. terraform init
# 4. terraform get
# 5. terraform plan
# 6. terraform apply

provider "google" {
  project = var.project
  region  = "europe-west2"
}

terraform {
  required_version = "~>1.5.6"
    backend "gcs" {
      bucket = "example-dev-terraform"
    }
}

module "infrastructure" {
  source  = "../../infrastructure/all"
  stage   = var.stage
  project = var.project
}