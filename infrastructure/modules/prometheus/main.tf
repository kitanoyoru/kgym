resource "kubernetes_namespace_v1" "prometheus" {
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

resource "helm_release" "prometheus" {
  name       = var.release_name
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"
  version    = var.chart_version != "" ? var.chart_version : null
  namespace  = var.namespace

  values = [
    yamlencode({
      prometheus = {
        prometheusSpec = {
          retention = var.retention
          storageSpec = {
            volumeClaimTemplate = {
              spec = {
                accessModes      = ["ReadWriteOnce"]
                storageClassName = var.storage_class != "" ? var.storage_class : null
                resources = {
                  requests = {
                    storage = var.storage_size
                  }
                }
              }
            }
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
        }
      }
      grafana = {
        enabled = false
      }
    })
  ]

  depends_on = [kubernetes_namespace_v1.prometheus]
}
