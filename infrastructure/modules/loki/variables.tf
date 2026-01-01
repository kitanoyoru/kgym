variable "namespace" {
  description = "Kubernetes namespace for Loki"
  type        = string
  default     = "monitoring"
}

variable "release_name" {
  description = "Helm release name for Loki"
  type        = string
  default     = "loki"
}

variable "chart_version" {
  description = "Version of the Loki Helm chart (leave empty for latest)"
  type        = string
  default     = ""
}

variable "storage_size" {
  description = "Storage size for Loki (e.g., 100Gi)"
  type        = string
  default     = "100Gi"
}

variable "storage_class" {
  description = "Storage class for persistent volumes"
  type        = string
  default     = ""
}

variable "replicas" {
  description = "Number of Loki replicas"
  type        = number
  default     = 1
}

variable "resources" {
  description = "Resource requests and limits for Loki"
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
