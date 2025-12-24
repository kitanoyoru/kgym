variable "namespace" {
  description = "Kubernetes namespace for CockroachDB deployment"
  type        = string
  default     = "cockroachdb"
}

variable "cluster_name" {
  description = "Name of the CockroachDB cluster"
  type        = string
}

variable "node_count" {
  description = "Number of CockroachDB nodes"
  type        = number
  default     = 3
  validation {
    condition     = var.node_count >= 3 && var.node_count % 2 == 1
    error_message = "Node count must be at least 3 and odd (for quorum)"
  }
}

variable "storage_size" {
  description = "Storage size for each node (e.g., 100Gi)"
  type        = string
  default     = "100Gi"
}

variable "storage_class" {
  description = "Storage class for persistent volumes"
  type        = string
  default     = ""
}

variable "resources" {
  description = "Resource requests and limits for CockroachDB pods"
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
      cpu    = "2"
      memory = "8Gi"
    }
    limits = {
      cpu    = "2"
      memory = "8Gi"
    }
  }
}

variable "image" {
  description = "CockroachDB container image"
  type        = string
  default     = "cockroachdb/cockroach:v23.1.0"
}

variable "tls_enabled" {
  description = "Enable TLS for CockroachDB"
  type        = bool
  default     = true
}

variable "labels" {
  description = "Labels to apply to resources"
  type        = map(string)
  default     = {}
}

variable "operator_namespace" {
  description = "Namespace for the CockroachDB operator"
  type        = string
  default     = "cockroach-operator-system"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}
