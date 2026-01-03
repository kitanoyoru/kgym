locals {
  monitoring_namespace_shared = var.prometheus_namespace == "monitoring" && var.loki_namespace == "monitoring" && var.grafana_namespace == "monitoring" && var.promtail_namespace == "monitoring"
  monitoring_enabled          = var.prometheus_enabled || var.loki_enabled || var.grafana_enabled || var.promtail_enabled
  config_base_path            = "${path.module}/../config/monitoring"
}

resource "kubernetes_namespace_v1" "monitoring" {
  count = local.monitoring_namespace_shared && local.monitoring_enabled ? 1 : 0
  metadata {
    name = "monitoring"
  }
}
