resource "kubernetes_namespace_v1" "loki" {
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

resource "kubernetes_config_map_v1" "loki_config" {
  count = var.config_file_path != null ? 1 : 0
  metadata {
    name      = "${var.release_name}-config"
    namespace = var.namespace
    labels    = var.labels
  }

  data = {
    "loki.yaml" = file(var.config_file_path)
  }

  depends_on = [kubernetes_namespace_v1.loki]
}

locals {
  common_values = {
    persistence = {
      enabled = true
      size    = var.storage_size
      storageClassName = var.storage_class != "" ? var.storage_class : null
    }
    deploymentMode = "SingleBinary"
    singleBinary = {
      replicas = var.replicas
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
    }
    read = {
      replicas = 0
    }
    write = {
      replicas = 0
    }
    backend = {
      replicas = 0
    }
    chunksCache = {
      replicas = 1
      resources = {
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
    resultsCache = {
      replicas = 1
      resources = {
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
  }

  values_with_config = merge(local.common_values, {
    loki = {
      config = {
        existingSecret = false
      }
      configMap = {
        name = kubernetes_config_map_v1.loki_config[0].metadata[0].name
      }
    }
  })

  values_without_config = merge(local.common_values, {
    loki = {
      auth_enabled = false
      useTestSchema = true
      storage = {
        type = "filesystem"
        bucketNames = {
          chunks = "chunks"
          ruler  = "ruler"
        }
      }
      commonConfig = {
        replication_factor = 1
      }
    }
  })
}

resource "helm_release" "loki" {
  name       = var.release_name
  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki"
  version    = var.chart_version != "" ? var.chart_version : null
  namespace  = var.namespace
  timeout    = var.timeout
  wait       = true

  values = var.config_file_path != null ? [
    yamlencode(local.values_with_config)
  ] : [
    yamlencode(local.values_without_config)
  ]

  depends_on = [
    kubernetes_namespace_v1.loki
  ]
}
