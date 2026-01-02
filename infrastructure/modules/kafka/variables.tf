variable "namespace" {
  description = "Kubernetes namespace for Kafka deployment"
  type        = string
  default     = "kafka"
}

variable "release_name" {
  description = "Helm release name for Kafka"
  type        = string
  default     = "kafka"
}

variable "chart_version" {
  description = "Version of the Kafka Helm chart (uses latest if empty)"
  type        = string
  default     = ""
}

variable "replicas" {
  description = "Number of Kafka broker replicas"
  type        = number
  default     = 1
  validation {
    condition     = var.replicas >= 1
    error_message = "Replicas must be at least 1"
  }
}

variable "storage_size" {
  description = "Storage size for each Kafka broker (e.g., 10Gi)"
  type        = string
  default     = "10Gi"
}

variable "storage_class" {
  description = "Storage class for persistent volumes"
  type        = string
  default     = ""
}

variable "resources" {
  description = "Resource requests and limits for Kafka brokers"
  type = object({
    requests = object({
      cpu    = string
      memory = string
    })
    limits = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    requests = {
      cpu    = "500m"
      memory = "1Gi"
    }
    limits = {
      cpu    = "2"
      memory = "2Gi"
    }
  }
}


variable "labels" {
  description = "Labels to apply to resources"
  type        = map(string)
  default     = {}
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "timeout" {
  description = "Timeout for Helm release installation (in seconds)"
  type        = number
  default     = 600
}

variable "image_repository" {
  description = "Docker image repository for Kafka (use 'bitnami/kafka' for secure images or 'bitnamilegacy/kafka' for legacy images)"
  type        = string
  default     = "bitnamilegacy/kafka"
}

variable "image_tag" {
  description = "Docker image tag for Kafka (leave empty to use chart default)"
  type        = string
  default     = ""
}

variable "sasl_users" {
  description = "SASL users for Kafka authentication (format: username:password)"
  type        = string
  default     = "admin:admin-secret"
  sensitive   = true
}
