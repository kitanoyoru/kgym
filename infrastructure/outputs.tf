output "cockroachdb" {
  description = "CockroachDB cluster outputs"
  value = var.cockroachdb_enabled ? {
    cluster_name   = module.cockroachdb[0].cluster_name
    namespace      = module.cockroachdb[0].namespace
    service_name   = module.cockroachdb[0].service_name
    service_endpoint = module.cockroachdb[0].service_endpoint
    http_endpoint  = module.cockroachdb[0].http_endpoint
    node_count     = module.cockroachdb[0].node_count
  } : null
}

output "redis" {
  description = "Redis cluster outputs"
  value = var.redis_enabled ? {
    cluster_name     = module.redis[0].cluster_name
    namespace        = module.redis[0].namespace
    service_name     = module.redis[0].service_name
    service_endpoint = module.redis[0].service_endpoint
    replicas         = module.redis[0].replicas
    password_secret_name = module.redis[0].password_secret_name
  } : null
}

output "prometheus" {
  description = "Prometheus outputs"
  value = var.prometheus_enabled ? {
    namespace      = module.prometheus[0].namespace
    service_name   = module.prometheus[0].service_name
    service_endpoint = module.prometheus[0].service_endpoint
  } : null
}

output "loki" {
  description = "Loki outputs"
  value = var.loki_enabled ? {
    namespace      = module.loki[0].namespace
    service_name   = module.loki[0].service_name
    service_endpoint = module.loki[0].service_endpoint
  } : null
}

output "grafana" {
  description = "Grafana outputs"
  value = var.grafana_enabled ? {
    namespace      = module.grafana[0].namespace
    service_name   = module.grafana[0].service_name
    service_endpoint = module.grafana[0].service_endpoint
    admin_user     = module.grafana[0].admin_user
  } : null
}

output "minio" {
  description = "MinIO cluster outputs"
  value = var.minio_enabled ? {
    cluster_name     = module.minio[0].cluster_name
    namespace        = module.minio[0].namespace
    service_name     = module.minio[0].service_name
    service_endpoint = module.minio[0].service_endpoint
    console_endpoint = module.minio[0].console_endpoint
    replicas         = module.minio[0].replicas
    access_key_secret_name = module.minio[0].access_key_secret_name
  } : null
}

output "kafka" {
  description = "Kafka cluster outputs"
  value = var.kafka_enabled ? {
    namespace        = module.kafka[0].namespace
    service_name     = module.kafka[0].service_name
    service_endpoint = module.kafka[0].service_endpoint
    replicas         = module.kafka[0].replicas
  } : null
}

output "sentry" {
  description = "Sentry outputs"
  value = var.sentry_enabled ? {
    namespace      = module.sentry[0].namespace
    service_name   = module.sentry[0].service_name
    service_endpoint = module.sentry[0].service_endpoint
    user_email     = module.sentry[0].user_email
  } : null
  sensitive = true
}
