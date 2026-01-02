resource "kubernetes_namespace_v1" "grafana" {
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

resource "kubernetes_secret_v1" "grafana_admin" {
  count = var.admin_password != null ? 1 : 0
  metadata {
    name      = "${var.release_name}-admin"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    admin-user     = base64encode(var.admin_user)
    admin-password = base64encode(var.admin_password)
  }
  depends_on = [kubernetes_namespace_v1.grafana]
}

locals {
  datasources_list = concat(
    var.prometheus_url != null ? [{
      name      = "Prometheus"
      type      = "prometheus"
      access    = "proxy"
      url       = var.prometheus_url
      isDefault = true
    }] : [],
    var.loki_url != null ? [{
      name   = "Loki"
      type   = "loki"
      access = "proxy"
      url    = var.loki_url
    }] : []
  )
}

resource "helm_release" "grafana" {
  name       = var.release_name
  repository = "https://grafana.github.io/helm-charts"
  chart      = "grafana"
  version    = var.chart_version
  namespace  = var.namespace

  values = [
    yamlencode({
      adminUser = var.admin_user
      adminPassword = var.admin_password != null ? var.admin_password : "admin"
      persistence = {
        enabled      = var.persistence_enabled
        size         = var.storage_size
        storageClass = var.storage_class != "" ? var.storage_class : null
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
      datasources = length(local.datasources_list) > 0 ? {
        "datasources.yaml" = {
          apiVersion = 1
          datasources = local.datasources_list
        }
      } : null
    })
  ]

  depends_on = [
    kubernetes_namespace_v1.grafana,
    kubernetes_secret_v1.grafana_admin
  ]
}
