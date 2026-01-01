variable "namespace" {
  description = "Kubernetes namespace for MinIO deployment"
  type        = string
  default     = "minio"
}

variable "cluster_name" {
  description = "Name of the MinIO cluster"
  type        = string
}

variable "replicas" {
  description = "Number of MinIO replicas"
  type        = number
  default     = 1
  validation {
    condition     = var.replicas >= 1
    error_message = "Replicas must be at least 1"
  }
}

variable "storage_size" {
  description = "Storage size for each MinIO pod (e.g., 100Gi)"
  type        = string
  default     = "100Gi"
}

variable "storage_class" {
  description = "Storage class for persistent volumes"
  type        = string
  default     = ""
}

variable "access_key" {
  description = "MinIO root access key"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "MinIO root secret key"
  type        = string
  sensitive   = true
}

variable "resources" {
  description = "Resource requests and limits for MinIO pods"
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
      memory = "2Gi"
    }
    limits = {
      cpu    = "2"
      memory = "4Gi"
    }
  }
}

variable "image" {
  description = "MinIO container image"
  type        = string
  default     = "minio/minio:latest"
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
