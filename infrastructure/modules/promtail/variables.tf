variable "namespace" {
  description = "Kubernetes namespace for Promtail"
  type        = string
  default     = "monitoring"
}

variable "release_name" {
  description = "Helm release name for Promtail"
  type        = string
  default     = "promtail"
}

variable "chart_version" {
  description = "Version of the Promtail Helm chart (leave empty for latest)"
  type        = string
  default     = ""
}

variable "config_file_path" {
  description = "Path to Promtail configuration YAML file"
  type        = string
}

variable "resources" {
  description = "Resource requests and limits for Promtail"
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
      cpu    = "100m"
      memory = "128Mi"
    }
    limits = {
      cpu    = "500m"
      memory = "512Mi"
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

variable "tolerations" {
  description = "Tolerations for Promtail pods"
  type        = list(any)
  default     = []
}

variable "node_selector" {
  description = "Node selector for Promtail pods"
  type        = map(string)
  default     = {}
}

variable "affinity" {
  description = "Affinity rules for Promtail pods"
  type        = any
  default     = {}
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
