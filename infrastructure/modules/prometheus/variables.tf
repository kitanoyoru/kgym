variable "namespace" {
  description = "Kubernetes namespace for Prometheus"
  type        = string
  default     = "monitoring"
}

variable "release_name" {
  description = "Helm release name for Prometheus"
  type        = string
  default     = "prometheus"
}

variable "chart_version" {
  description = "Version of the Prometheus Helm chart (leave empty for latest)"
  type        = string
  default     = ""
}

variable "storage_size" {
  description = "Storage size for Prometheus (e.g., 50Gi)"
  type        = string
  default     = "50Gi"
}

variable "storage_class" {
  description = "Storage class for persistent volumes"
  type        = string
  default     = ""
}

variable "retention" {
  description = "How long to retain metrics (e.g., 15d)"
  type        = string
  default     = "15d"
}

variable "resources" {
  description = "Resource requests and limits for Prometheus"
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
      cpu    = "1"
      memory = "2Gi"
    }
    limits = {
      cpu    = "2"
      memory = "4Gi"
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

variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = null
  nullable    = true
}

variable "kubeconfig_context" {
  description = "Kubernetes context to use"
  type        = string
  default     = null
  nullable    = true
}
