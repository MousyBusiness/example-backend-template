output "neg" {
  description = "The network endpoint group id"
  value       = google_compute_region_network_endpoint_group.cloudrun_neg.id
}

output "invoke_url" {
  value = google_cloud_run_service.default.status[0].url
}
