variable "namespace" {
  description = "Kubernetes namespace for Grafana"
  type        = string
  default     = "monitoring"
}

variable "release_name" {
  description = "Helm release name for Grafana"
  type        = string
  default     = "grafana"
}

variable "chart_version" {
  description = "Version of the Grafana Helm chart"
  type        = string
  default     = "7.3.7"
}

variable "admin_user" {
  description = "Grafana admin username"
  type        = string
  default     = "admin"
}

variable "admin_password" {
  description = "Grafana admin password"
  type        = string
  default     = null
  nullable    = true
  sensitive   = true
}

variable "persistence_enabled" {
  description = "Enable persistent storage for Grafana"
  type        = bool
  default     = true
}

variable "storage_size" {
  description = "Storage size for Grafana (e.g., 10Gi)"
  type        = string
  default     = "10Gi"
}

variable "storage_class" {
  description = "Storage class for persistent volumes"
  type        = string
  default     = ""
}

variable "resources" {
  description = "Resource requests and limits for Grafana"
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
      cpu    = "200m"
      memory = "256Mi"
    }
    limits = {
      cpu    = "500m"
      memory = "512Mi"
    }
  }
}

variable "prometheus_url" {
  description = "URL of Prometheus service (for data source configuration)"
  type        = string
  default     = null
  nullable    = true
}

variable "loki_url" {
  description = "URL of Loki service (for data source configuration)"
  type        = string
  default     = null
  nullable    = true
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

variable "config_file_path" {
  description = "Path to Grafana configuration INI file (optional, if not provided uses default Helm values)"
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
