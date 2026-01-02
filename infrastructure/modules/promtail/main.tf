resource "kubernetes_namespace_v1" "promtail" {
  count = var.create_namespace ? 1 : 0
  metadata {
    name = var.namespace
    labels = merge(
      var.labels,
      {
        name = var.namespace
      }
    )
  }
}

resource "kubernetes_config_map_v1" "promtail_config" {
  metadata {
    name      = "${var.release_name}-config"
    namespace = var.namespace
    labels    = var.labels
  }

  data = {
    "promtail.yaml" = file(var.config_file_path)
  }

  depends_on = [kubernetes_namespace_v1.promtail]
}

resource "helm_release" "promtail" {
  name       = var.release_name
  repository = "https://grafana.github.io/helm-charts"
  chart      = "promtail"
  version    = var.chart_version != "" ? var.chart_version : null
  namespace  = var.namespace
  timeout    = var.timeout
  wait       = true

  values = [
    yamlencode({
      config = {
        existingSecret = false
      }
      configMap = {
        name = kubernetes_config_map_v1.promtail_config.metadata[0].name
      }
      serviceAccount = {
        create = true
        name   = var.release_name
      }
      rbac = {
        create = true
      }
      resources = {
        requests = {
          cpu    = var.resources.requests.cpu
          memory = var.resources.requests.memory
        }
        limits = {
          cpu    = var.resources.limits.cpu
          memory = var.resources.limits.memory
        }
      }
      tolerations = var.tolerations
      nodeSelector = var.node_selector
      affinity = var.affinity
    })
  ]

  depends_on = [
    kubernetes_namespace_v1.promtail,
    kubernetes_config_map_v1.promtail_config
  ]
}
