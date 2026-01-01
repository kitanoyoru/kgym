resource "kubernetes_namespace" "loki" {
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

resource "helm_release" "loki" {
  name       = var.release_name
  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki"
  version    = var.chart_version != "" ? var.chart_version : null
  namespace  = var.namespace

  values = [
    yamlencode({
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
    })
  ]

  depends_on = [kubernetes_namespace.loki]
}
