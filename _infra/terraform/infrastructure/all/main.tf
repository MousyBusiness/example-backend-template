resource "google_compute_global_address" "default" {
  name = "example-ip"
}

resource "google_compute_health_check" "http_health_check" {
  name = "http-health-check"

  http_health_check {
    port         = 80
    request_path = "/healthcheck"
  }
}

module "ci_cd" {
  source = "../../modules/cloudbuild"
  stage  = var.stage
}

module "service_accounts" {
  source  = "../global/service-accounts"
  stage   = var.stage
  project = var.project
}

module "buckets" {
  source  = "../global/buckets"
  stage   = var.stage
  project = var.project
}

module "network" {
  source  = "../../modules/network"
  stage   = var.stage
  project = var.project
}

module "services_north_america" {
  source                 = "../regional"
  stage                  = var.stage
  project                = var.project
  region                 = "us-west2"
  example_service_account   = module.service_accounts.example_service_account
  vpc_connector_cidr     = "10.10.0.0/28"
  health_check           = google_compute_health_check.http_health_check.self_link
}


module "lb" {
  source      = "../../modules/load-balancer"
  stage       = var.stage
  project     = var.project
  external_ip = google_compute_global_address.default.address
  example_negs   = [/*module.services_europe.example_neg, */module.services_north_america.example_neg]
}