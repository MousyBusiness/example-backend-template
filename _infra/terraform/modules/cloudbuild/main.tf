locals {
  branch = var.stage == "dev" ? "develop" : "master"
  config = "cloudbuild_${var.stage}.yaml"
}

resource "google_cloudbuild_trigger" "backend_github_trigger" {
  github {
    owner = "ORG_NAME"
    name  = "example-repo-backend"
    push {
      branch = local.branch
    }
  }

  name        = var.stage == "backend-build"
  description = var.stage == "Builds backend services on push to master."

  filename = local.config
}
