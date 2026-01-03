variable "kubeconfig_path" {
  description = "Path to kubeconfig file. If not set, uses default kubeconfig location (~/.kube/config or KUBECONFIG env var)"
  type        = string
  default     = null
  nullable    = true
}

variable "kubeconfig_context" {
  description = "Kubernetes context to use. If not set, uses current context from kubeconfig"
  type        = string
  default     = null
  nullable    = true
}

variable "cockroachdb_enabled" {
  description = "Enable CockroachDB deployment"
  type        = bool
  default     = true
}

variable "cockroachdb_cluster_name" {
  description = "Name of the CockroachDB cluster"
  type        = string
  default     = "cockroachdb"
}

variable "cockroachdb_namespace" {
  description = "Kubernetes namespace for CockroachDB"
  type        = string
  default     = "cockroachdb"
}

variable "cockroachdb_node_count" {
  description = "Number of CockroachDB nodes"
  type        = number
  default     = 3
}

variable "cockroachdb_storage_size" {
  description = "Storage size for each CockroachDB node"
  type        = string
  default     = "10Gi"
}

variable "cockroachdb_storage_class" {
  description = "Storage class for CockroachDB persistent volumes"
  type        = string
  default     = ""
}

variable "cockroachdb_tls_enabled" {
  description = "Enable TLS for CockroachDB (disabled by default to allow passwordless connections)"
  type        = bool
  default     = false
}

variable "cockroachdb_resources" {
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
      cpu    = "500m"
      memory = "1Gi"
    }
    limits = {
      cpu    = "2"
      memory = "2Gi"
    }
  }
}

variable "redis_enabled" {
  description = "Enable Redis deployment"
  type        = bool
  default     = true
}

variable "redis_cluster_name" {
  description = "Name of the Redis cluster"
  type        = string
  default     = "redis"
}

variable "redis_namespace" {
  description = "Kubernetes namespace for Redis"
  type        = string
  default     = "redis"
}

variable "redis_replicas" {
  description = "Number of Redis replicas"
  type        = number
  default     = 3
}

variable "redis_storage_size" {
  description = "Storage size for each Redis pod"
  type        = string
  default     = "1Gi"
}

variable "redis_storage_class" {
  description = "Storage class for Redis persistent volumes"
  type        = string
  default     = ""
}

variable "redis_password" {
  description = "Redis password (optional, leave null for no password)"
  type        = string
  default     = null
  nullable    = true
  sensitive   = true
}

variable "prometheus_enabled" {
  description = "Enable Prometheus deployment"
  type        = bool
  default     = true
}

variable "prometheus_namespace" {
  description = "Kubernetes namespace for Prometheus"
  type        = string
  default     = "monitoring"
}

variable "prometheus_storage_size" {
  description = "Storage size for Prometheus"
  type        = string
  default     = "5Gi"
}

variable "prometheus_storage_class" {
  description = "Storage class for Prometheus persistent volumes"
  type        = string
  default     = ""
}

variable "loki_enabled" {
  description = "Enable Loki deployment"
  type        = bool
  default     = true
}

variable "loki_namespace" {
  description = "Kubernetes namespace for Loki"
  type        = string
  default     = "monitoring"
}

variable "loki_storage_size" {
  description = "Storage size for Loki"
  type        = string
  default     = "1Gi"
}

variable "loki_storage_class" {
  description = "Storage class for Loki persistent volumes"
  type        = string
  default     = ""
}

variable "grafana_enabled" {
  description = "Enable Grafana deployment"
  type        = bool
  default     = true
}

variable "grafana_namespace" {
  description = "Kubernetes namespace for Grafana"
  type        = string
  default     = "monitoring"
}

variable "grafana_admin_password" {
  description = "Grafana admin password"
  type        = string
  default     = null
  nullable    = true
  sensitive   = true
}

variable "grafana_storage_size" {
  description = "Storage size for Grafana"
  type        = string
  default     = "1Gi"
}

variable "grafana_storage_class" {
  description = "Storage class for Grafana persistent volumes"
  type        = string
  default     = ""
}

variable "grafana_config_file_path" {
  description = "Path to Grafana configuration file (optional, defaults to config/monitoring/grafana/grafana.ini)"
  type        = string
  default     = null
  nullable    = true
}

variable "loki_config_file_path" {
  description = "Path to Loki configuration file (optional, defaults to config/monitoring/loki/loki.yaml)"
  type        = string
  default     = null
  nullable    = true
}

variable "promtail_enabled" {
  description = "Enable Promtail deployment"
  type        = bool
  default     = true
}

variable "promtail_namespace" {
  description = "Kubernetes namespace for Promtail"
  type        = string
  default     = "monitoring"
}

variable "promtail_config_file_path" {
  description = "Path to Promtail configuration file (optional, defaults to config/monitoring/promtail/promtail.yaml)"
  type        = string
  default     = null
  nullable    = true
}

variable "minio_enabled" {
  description = "Enable MinIO deployment"
  type        = bool
  default     = true
}

variable "minio_cluster_name" {
  description = "Name of the MinIO cluster"
  type        = string
  default     = "minio"
}

variable "minio_namespace" {
  description = "Kubernetes namespace for MinIO"
  type        = string
  default     = "minio"
}

variable "minio_replicas" {
  description = "Number of MinIO replicas"
  type        = number
  default     = 1
}

variable "minio_storage_size" {
  description = "Storage size for each MinIO pod"
  type        = string
  default     = "1Gi"
}

variable "minio_storage_class" {
  description = "Storage class for MinIO persistent volumes"
  type        = string
  default     = ""
}

variable "minio_access_key" {
  description = "MinIO root access key"
  type        = string
  default     = "minioadmin"
  sensitive   = true
}

variable "minio_secret_key" {
  description = "MinIO root secret key"
  type        = string
  default     = "minioadmin"
  sensitive   = true
}

variable "sentry_enabled" {
  description = "Enable Sentry deployment"
  type        = bool
  default     = false
}

variable "sentry_namespace" {
  description = "Kubernetes namespace for Sentry"
  type        = string
  default     = "sentry"
}

variable "sentry_user_email" {
  description = "Email for the initial Sentry superuser"
  type        = string
  default     = "admin@sentry.local"
}

variable "sentry_user_password" {
  description = "Password for the initial Sentry superuser"
  type        = string
  default     = "sentry"
  sensitive   = true
}

variable "sentry_postgresql_host" {
  description = "PostgreSQL host (leave empty to use embedded PostgreSQL from Sentry Helm chart, or set to external PostgreSQL host)"
  type        = string
  default     = ""
}

variable "sentry_postgresql_port" {
  description = "PostgreSQL port"
  type        = number
  default     = 26257
}

variable "sentry_postgresql_database" {
  description = "PostgreSQL database name"
  type        = string
  default     = "sentry"
}

variable "sentry_postgresql_user" {
  description = "PostgreSQL username"
  type        = string
  default     = "root"
}

variable "sentry_postgresql_password" {
  description = "PostgreSQL password"
  type        = string
  default     = ""
  sensitive   = true
}

variable "sentry_redis_host" {
  description = "Redis host (leave empty to use embedded Redis, or set to Redis service endpoint)"
  type        = string
  default     = ""
}

variable "sentry_redis_port" {
  description = "Redis port"
  type        = number
  default     = 6379
}

variable "sentry_redis_password" {
  description = "Redis password (if Redis requires authentication)"
  type        = string
  default     = null
  nullable    = true
  sensitive   = true
}

variable "sentry_storage_size" {
  description = "Storage size for Sentry data"
  type        = string
  default     = "10Gi"
}

variable "sentry_storage_class" {
  description = "Storage class for Sentry persistent volumes"
  type        = string
  default     = ""
}
