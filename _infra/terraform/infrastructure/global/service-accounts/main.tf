locals {
  datastore           = "${var.project}=>roles/datastore.user"
  functions           = "${var.project}=>roles/cloudfunctions.invoker"
  secrets             = "${var.project}=>roles/secretmanager.secretAccessor"
  secrets_admin       = "${var.project}=>roles/secretmanager.admin"
  compute_log         = "${var.project}=>roles/logging.logWriter"
  compute_metric      = "${var.project}=>roles/monitoring.metricWriter"
  compute_agent       = "${var.project}=>roles/cloudtrace.agent"
  compute_viewer      = "${var.project}=>roles/compute.viewer"
  storage_obj_viewer  = "${var.project}=>roles/storage.objectViewer"
  storage_obj_creator = "${var.project}=>roles/storage.objectCreator"
  storage_obj_admin   = "${var.project}=>roles/storage.objectAdmin"
  container           = "${var.project}=>roles/containerregistry.ServiceAgent"
  cloud_build         = "${var.project}=>roles/cloudbuild.builds.builder"
  pubsub_publisher    = "${var.project}=>roles/pubsub.publisher"
  pubsub_subscriber   = "${var.project}=>roles/pubsub.subscriber"
  memorystore_viewer  = "${var.project}=>roles/redis.viewer"
  run_invoker         = "${var.project}=>roles/run.invoker"
}


module "example_service_account" {
  source        = "terraform-google-modules/service-accounts/google"
  version       = "~> 3.0"
  project_id    = var.project
  names         = ["user-sa"]
  project_roles = [
    local.datastore,
    local.secrets,
    local.pubsub_publisher,
  ]
}
