variable "stage" {}

variable "name" {}

variable "project" {}

variable "service_account" {}

variable "region" {}

variable "ingress" {
#  https://cloud.google.com/sdk/gcloud/reference/run/deploy#--ingress
  default = "internal-and-cloud-load-balancing"
}

variable "require_authentication" {
  default = false
}

variable "redis_ip" {
  default = null
}

variable "image" {
  type        = string
  description = "The container image located in Container Registry."

  validation {
    condition     = can(regex("^eu.gcr.io/", var.image))
    error_message = "The image value must be in container registry."
  }
}

variable "max_instances" {
  default = 1
}

variable "min_instances" {
  default = 0
}