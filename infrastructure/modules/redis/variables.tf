variable "namespace" {
  description = "Kubernetes namespace for Redis deployment"
  type        = string
  default     = "redis"
}

variable "cluster_name" {
  description = "Name of the Redis cluster"
  type        = string
}

variable "replicas" {
  description = "Number of Redis replicas"
  type        = number
  default     = 3
  validation {
    condition     = var.replicas >= 1
    error_message = "Replicas must be at least 1"
  }
}

variable "storage_size" {
  description = "Storage size for each Redis pod (e.g., 10Gi)"
  type        = string
  default     = "10Gi"
}

variable "storage_class" {
  description = "Storage class for persistent volumes"
  type        = string
  default     = ""
}

variable "resources" {
  description = "Resource requests and limits for Redis pods"
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
      cpu    = "1"
      memory = "2Gi"
    }
  }
}

variable "image" {
  description = "Redis container image"
  type        = string
  default     = "redis:7-alpine"
}

variable "password" {
  description = "Redis password (optional, leave null for no password)"
  type        = string
  default     = null
  nullable    = true
  sensitive   = true
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
