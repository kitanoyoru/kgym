resource "kubernetes_namespace_v1" "minio" {
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

resource "kubernetes_secret_v1" "minio_credentials" {
  metadata {
    name      = "${var.cluster_name}-credentials"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    MINIO_ROOT_USER     = base64encode(var.access_key)
    MINIO_ROOT_PASSWORD = base64encode(var.secret_key)
  }
  depends_on = [kubernetes_namespace_v1.minio]
}

resource "kubernetes_service_v1" "minio_api" {
  metadata {
    name      = "${var.cluster_name}-api"
    namespace = var.namespace
    labels    = merge(var.labels, { app = var.cluster_name, service = "api" })
  }
  spec {
    type = "ClusterIP"
    selector = {
      app = var.cluster_name
    }
    port {
      port        = 9000
      target_port = 9000
      name       = "api"
    }
  }
  depends_on = [kubernetes_namespace_v1.minio]
}

resource "kubernetes_service_v1" "minio_console" {
  metadata {
    name      = "${var.cluster_name}-console"
    namespace = var.namespace
    labels    = merge(var.labels, { app = var.cluster_name, service = "console" })
  }
  spec {
    type = "ClusterIP"
    selector = {
      app = var.cluster_name
    }
    port {
      port        = 9001
      target_port = 9001
      name       = "console"
    }
  }
  depends_on = [kubernetes_namespace_v1.minio]
}

resource "kubernetes_stateful_set_v1" "minio" {
  metadata {
    name      = var.cluster_name
    namespace = var.namespace
    labels    = merge(var.labels, { app = var.cluster_name })
  }
  spec {
    service_name = kubernetes_service_v1.minio_api.metadata[0].name
    replicas     = var.replicas
    selector {
      match_labels = {
        app = var.cluster_name
      }
    }
    template {
      metadata {
        labels = merge(var.labels, { app = var.cluster_name })
      }
      spec {
        container {
          name  = "minio"
          image = var.image
          args = [
            "server",
            "/data",
            "--console-address",
            ":9001"
          ]
          port {
            container_port = 9000
            name          = "api"
          }
          port {
            container_port = 9001
            name          = "console"
          }
          env {
            name = "MINIO_ROOT_USER"
            value_from {
              secret_key_ref {
                name = kubernetes_secret_v1.minio_credentials.metadata[0].name
                key  = "MINIO_ROOT_USER"
              }
            }
          }
          env {
            name = "MINIO_ROOT_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret_v1.minio_credentials.metadata[0].name
                key  = "MINIO_ROOT_PASSWORD"
              }
            }
          }
          resources {
            requests = {
              cpu    = var.resources.requests.cpu
              memory = var.resources.requests.memory
            }
            limits = {
              cpu    = var.resources.limits.cpu
              memory = var.resources.limits.memory
            }
          }
          volume_mount {
            name       = "data"
            mount_path = "/data"
          }
          liveness_probe {
            http_get {
              path = "/minio/health/live"
              port = 9000
            }
            initial_delay_seconds = 30
            period_seconds       = 20
          }
          readiness_probe {
            http_get {
              path = "/minio/health/ready"
              port = 9000
            }
            initial_delay_seconds = 10
            period_seconds       = 10
          }
        }
      }
    }
    volume_claim_template {
      metadata {
        name = "data"
      }
      spec {
        access_modes       = ["ReadWriteOnce"]
        storage_class_name = var.storage_class != "" ? var.storage_class : null
        resources {
          requests = {
            storage = var.storage_size
          }
        }
      }
    }
  }
  depends_on = [
    kubernetes_namespace_v1.minio,
    kubernetes_secret_v1.minio_credentials,
    kubernetes_service_v1.minio_api,
    kubernetes_service_v1.minio_console
  ]
}
