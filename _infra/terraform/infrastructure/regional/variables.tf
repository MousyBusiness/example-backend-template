variable "stage" {}

variable "project" {}

variable "region" {}

variable "example_service_account" {}

variable "health_check" {}

variable "vpc_connector_cidr" {}

variable "relay_min_replicas" {
  default = 1
}

variable "relay_max_replicas" {
  default = 2
}
