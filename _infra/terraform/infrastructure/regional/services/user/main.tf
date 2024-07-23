module "user" {
  source          = "../../../../modules/cloudrun"
  stage           = var.stage
  name            = "user-${var.region}"
  project         = var.project
  region          = var.region
  redis_ip        = var.redis_ip
  service_account = var.service_account
  min_instances   = var.min_replicas
  max_instances   = var.max_replicas
  image           = "eu.gcr.io/${var.project}/user:latest"
}