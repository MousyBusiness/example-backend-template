resource "google_vpc_access_connector" "serverless_vpc_connector" {
  name          = "svpc-${var.region}"
  ip_cidr_range = var.vpc_connector_cidr
  network       = "default"
  region        = var.region
}


module "user" {
  source            = "./services/user"
  stage             = var.stage
  project           = var.project
  region            = var.region
  service_account   = var.example_service_account
  min_replicas      = var.stage == "dev" ? 0 : 0
  max_replicas      = var.stage == "dev" ? 1 : 1
  depends_on        = [google_vpc_access_connector.serverless_vpc_connector]
}
