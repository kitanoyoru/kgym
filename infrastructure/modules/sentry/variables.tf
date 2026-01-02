variable "namespace" {
  description = "Kubernetes namespace for Sentry deployment"
  type        = string
  default     = "sentry"
}

variable "release_name" {
  description = "Helm release name for Sentry"
  type        = string
  default     = "sentry"
}

variable "chart_version" {
  description = "Version of the Sentry Helm chart (leave empty for latest)"
  type        = string
  default     = ""
}

variable "user_email" {
  description = "Email for the initial Sentry superuser"
  type        = string
}

variable "user_password" {
  description = "Password for the initial Sentry superuser"
  type        = string
  sensitive   = true
}

variable "postgresql_host" {
  description = "PostgreSQL host (can use CockroachDB service endpoint)"
  type        = string
  default     = ""
}

variable "postgresql_port" {
  description = "PostgreSQL port"
  type        = number
  default     = 26257
}

variable "postgresql_database" {
  description = "PostgreSQL database name"
  type        = string
  default     = "sentry"
}

variable "postgresql_user" {
  description = "PostgreSQL username"
  type        = string
  default     = "root"
}

variable "postgresql_password" {
  description = "PostgreSQL password"
  type        = string
  default     = ""
  sensitive   = true
}

variable "redis_host" {
  description = "Redis host (can use Redis service endpoint)"
  type        = string
  default     = ""
}

variable "redis_port" {
  description = "Redis port"
  type        = number
  default     = 6379
}

variable "redis_password" {
  description = "Redis password (if set)"
  type        = string
  default     = null
  nullable    = true
  sensitive   = true
}

variable "kafka_host" {
  description = "Kafka host (can use Kafka service endpoint)"
  type        = string
  default     = ""
}

variable "kafka_port" {
  description = "Kafka port"
  type        = number
  default     = 9092
}

variable "kafka_sasl_username" {
  description = "Kafka SASL username"
  type        = string
  default     = null
  nullable    = true
  sensitive   = true
}

variable "kafka_sasl_password" {
  description = "Kafka SASL password"
  type        = string
  default     = null
  nullable    = true
  sensitive   = true
}

variable "storage_size" {
  description = "Storage size for Sentry data (e.g., 10Gi)"
  type        = string
  default     = "10Gi"
}

variable "storage_class" {
  description = "Storage class for persistent volumes"
  type        = string
  default     = ""
}

variable "resources" {
  description = "Resource requests and limits for Sentry web service"
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

variable "worker_resources" {
  description = "Resource requests and limits for Sentry worker service"
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
      memory = "512Mi"
    }
    limits = {
      cpu    = "1"
      memory = "1Gi"
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
