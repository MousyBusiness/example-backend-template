variable "stage" {}

variable "service_account" {}

variable "project" {}

variable "redis_ip" {}

variable "min_replicas" {
}

variable "max_replicas" {
}

variable "app_engine_region" {
  default = "europe-west2"
}

variable "region" {
  default = "europe-west2"
}