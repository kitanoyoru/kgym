locals {
  monitoring_namespace_shared = var.prometheus_namespace == "monitoring" && var.loki_namespace == "monitoring" && var.grafana_namespace == "monitoring" && var.promtail_namespace == "monitoring"
  monitoring_enabled = var.prometheus_enabled || var.loki_enabled || var.grafana_enabled || var.promtail_enabled
  config_base_path   = "${path.module}/../config/monitoring"
}

resource "kubernetes_namespace_v1" "monitoring" {
  count = local.monitoring_namespace_shared && local.monitoring_enabled ? 1 : 0
  metadata {
    name = "monitoring"
  }
}

module "cockroachdb" {
  count  = var.cockroachdb_enabled ? 1 : 0
  source = "./modules/cockroachdb"

  cluster_name  = var.cockroachdb_cluster_name
  namespace     = var.cockroachdb_namespace
  node_count    = var.cockroachdb_node_count
  storage_size  = var.cockroachdb_storage_size
  storage_class = var.cockroachdb_storage_class
  resources     = var.cockroachdb_resources
  tls_enabled   = var.cockroachdb_tls_enabled

  depends_on = [kubernetes_namespace_v1.cockroachdb]
}

module "redis" {
  count  = var.redis_enabled ? 1 : 0
  source = "./modules/redis"

  cluster_name  = var.redis_cluster_name
  namespace     = var.redis_namespace
  replicas      = var.redis_replicas
  storage_size  = var.redis_storage_size
  storage_class = var.redis_storage_class
  password      = var.redis_password

  depends_on = [kubernetes_namespace_v1.redis]
}

module "prometheus" {
  count  = var.prometheus_enabled ? 1 : 0
  source = "./modules/prometheus"

  namespace       = var.prometheus_namespace
  storage_size     = var.prometheus_storage_size
  storage_class    = var.prometheus_storage_class
  create_namespace = !local.monitoring_namespace_shared

  depends_on = [kubernetes_namespace_v1.monitoring]
}

module "loki" {
  count  = var.loki_enabled ? 1 : 0
  source = "./modules/loki"

  namespace        = var.loki_namespace
  storage_size     = var.loki_storage_size
  storage_class    = var.loki_storage_class
  config_file_path = var.loki_config_file_path != null ? var.loki_config_file_path : "${local.config_base_path}/loki/loki.yaml"
  create_namespace = !local.monitoring_namespace_shared

  depends_on = [kubernetes_namespace_v1.monitoring]
}

module "promtail" {
  count  = var.promtail_enabled ? 1 : 0
  source = "./modules/promtail"

  namespace        = var.promtail_namespace
  config_file_path = var.promtail_config_file_path != null ? var.promtail_config_file_path : "${local.config_base_path}/promtail/promtail.yaml"
  create_namespace = !local.monitoring_namespace_shared

  depends_on = [
    kubernetes_namespace_v1.monitoring,
    module.loki
  ]
}

module "grafana" {
  count  = var.grafana_enabled ? 1 : 0
  source = "./modules/grafana"

  namespace        = var.grafana_namespace
  admin_password   = var.grafana_admin_password
  storage_size     = var.grafana_storage_size
  storage_class    = var.grafana_storage_class
  config_file_path = var.grafana_config_file_path != null ? var.grafana_config_file_path : "${local.config_base_path}/grafana/grafana.ini"
  prometheus_url   = var.prometheus_enabled ? "http://${module.prometheus[0].service_endpoint}" : null
  loki_url         = var.loki_enabled ? "http://${module.loki[0].service_endpoint}" : null
  create_namespace = !local.monitoring_namespace_shared

  depends_on = [kubernetes_namespace_v1.monitoring]
}

module "minio" {
  count  = var.minio_enabled ? 1 : 0
  source = "./modules/minio"

  cluster_name  = var.minio_cluster_name
  namespace     = var.minio_namespace
  replicas      = var.minio_replicas
  storage_size  = var.minio_storage_size
  storage_class = var.minio_storage_class
  access_key    = var.minio_access_key
  secret_key    = var.minio_secret_key
}

module "kafka" {
  count  = var.kafka_enabled ? 1 : 0
  source = "./modules/kafka"

  namespace        = var.kafka_namespace
  replicas         = var.kafka_replicas
  storage_size     = var.kafka_storage_size
  storage_class    = var.kafka_storage_class
  resources        = var.kafka_resources
  image_repository = var.kafka_image_repository
  sasl_users       = "${var.kafka_sasl_user}:${var.kafka_sasl_password}"
}

module "sentry" {
  count  = var.sentry_enabled ? 1 : 0
  source = "./modules/sentry"

  namespace            = var.sentry_namespace
  user_email          = var.sentry_user_email
  user_password       = var.sentry_user_password
  postgresql_host     = var.sentry_postgresql_host
  postgresql_port     = var.sentry_postgresql_port
  postgresql_database = var.sentry_postgresql_database
  postgresql_user     = var.sentry_postgresql_user
  postgresql_password = var.sentry_postgresql_password
  redis_host         = var.sentry_redis_host != "" ? var.sentry_redis_host : (var.redis_enabled ? "${var.redis_cluster_name}.${var.redis_namespace}.svc.cluster.local" : "")
  redis_port         = var.sentry_redis_port
  redis_password     = var.sentry_redis_password != null ? var.sentry_redis_password : var.redis_password
  kafka_host         = var.sentry_kafka_host != "" ? var.sentry_kafka_host : (var.kafka_enabled ? split(":", module.kafka[0].service_endpoint)[0] : "")
  kafka_port         = var.sentry_kafka_port
  kafka_sasl_username = var.kafka_enabled ? var.kafka_sasl_user : null
  kafka_sasl_password = var.kafka_enabled ? var.kafka_sasl_password : null
  storage_size       = var.sentry_storage_size
  storage_class      = var.sentry_storage_class

  depends_on = [
    module.redis,
    module.kafka
  ]
}
