resource "google_compute_firewall" "healthcheck_fw" {
  project     = var.project
  name        = "healthcheck-fw"
  network     = "default"
  description = "Allow http traffic for managed instance and load balancer health checks"

  allow {
    protocol  = "tcp"
    ports     = ["80"]
  }
  source_ranges = ["35.191.0.0/16", "130.211.0.0/22"]
  target_tags = ["healthcheck"]
}
