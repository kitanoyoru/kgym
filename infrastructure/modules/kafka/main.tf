resource "kubernetes_namespace_v1" "kafka" {
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

resource "helm_release" "kafka" {
  name       = var.release_name
  repository = "oci://registry-1.docker.io/bitnamicharts"
  chart      = "kafka"
  version    = var.chart_version != "" ? var.chart_version : null
  namespace  = var.namespace
  timeout    = var.timeout
  wait       = true

  depends_on = [
    kubernetes_namespace_v1.kafka
  ]

  values = [
    yamlencode(merge(
      {
        replicaCount = var.replicas

        kraft = {
          enabled = true
        }

        image = merge(
          {
            registry   = "docker.io"
            repository = var.image_repository
            pullPolicy = "IfNotPresent"
          },
          var.image_tag != "" ? {
            tag = var.image_tag
          } : {}
        )

        persistence = {
          enabled      = true
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
      }
    ))
  ]
}
