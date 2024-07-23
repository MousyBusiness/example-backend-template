//https://cloud.google.com/load-balancing/docs/https/ext-http-lb-tf-module-examples#cloud-run
module "lb_http" {
  source                          = "GoogleCloudPlatform/lb-http/google"
  version                         = "6.1.0"
  name                            = "lb"
  project                         = var.project
  address                         = var.external_ip
  ssl                             = true
  use_ssl_certificates            = false
  http_forward                    = false
  create_url_map                  = false
  url_map                         = google_compute_url_map.default.self_link
  managed_ssl_certificate_domains = [
    "example.${var.stage}.yourdomain.com",
  ]

  backends = {
    example-bes = {
      protocol                        = "HTTP"
      timeout_sec                     = 30
      enable_cdn                      = false
      health_check                    = null
      affinity_cookie_ttl_sec         = null
      connection_draining_timeout_sec = null
      custom_request_headers          = null
      custom_response_headers         = null
      description                     = null
      port                            = 80
      port_name                       = null
      security_policy                 = null
      session_affinity                = null
      log_config                      = {
        enable      = false
        sample_rate = null
      }
      groups                          = [
      for neg in var.example_negs :
      {
          group                        = neg
          balancing_mode               = null
          capacity_scaler              = null
          description                  = null
          max_connections              = null
          max_connections_per_instance = null
          max_connections_per_endpoint = null
          max_rate                     = null
          max_rate_per_instance        = null
          max_rate_per_endpoint        = null
          max_utilization              = null
        }
      ]

      iap_config = {
        enable               = false
        oauth2_client_id     = null
        oauth2_client_secret = null
      }
    }
  }
}

resource "google_compute_url_map" "default" {
  name            = "example-lb"
  default_service = module.lb_http.backend_services["example-bes"].self_link

  host_rule {
    path_matcher = "example-pm"
    hosts        = ["example.${var.stage}.yourdomain.com"]
  }

  path_matcher {
    name            = "example-pm"
    default_service = module.lb_http.backend_services["example-bes"].self_link
    path_rule {
      paths   = [
        "/"
      ]
      service = module.lb_http.backend_services["example-bes"].self_link
    }
  }
}
