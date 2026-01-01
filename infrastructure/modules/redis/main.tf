resource "kubernetes_namespace" "redis" {
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

resource "kubernetes_secret" "redis_password" {
  count = var.password != null ? 1 : 0
  metadata {
    name      = "${var.cluster_name}-password"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    password = base64encode(var.password)
  }
  depends_on = [kubernetes_namespace.redis]
}

resource "kubernetes_config_map" "redis_config" {
  metadata {
    name      = "${var.cluster_name}-config"
    namespace = var.namespace
    labels    = var.labels
  }
  data = var.password != null ? {
    "redis.conf" = "requirepass ${var.password}\n"
  } : {
    "redis.conf" = "# Redis configuration\n"
  }
  depends_on = [kubernetes_namespace.redis]
}

resource "kubernetes_service" "redis" {
  metadata {
    name      = var.cluster_name
    namespace = var.namespace
    labels    = merge(var.labels, { app = var.cluster_name })
  }
  spec {
    type = "ClusterIP"
    selector = {
      app = var.cluster_name
    }
    port {
      port        = 6379
      target_port = 6379
      name       = "redis"
    }
  }
  depends_on = [kubernetes_namespace.redis]
}

resource "kubernetes_stateful_set" "redis" {
  metadata {
    name      = var.cluster_name
    namespace = var.namespace
    labels    = merge(var.labels, { app = var.cluster_name })
  }
  spec {
    service_name = kubernetes_service.redis.metadata[0].name
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
          name  = "redis"
          image = var.image
          port {
            container_port = 6379
            name          = "redis"
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
          volume_mount {
            name       = "config"
            mount_path = "/etc/redis"
          }
          command = var.password != null ? [
            "redis-server",
            "/etc/redis/redis.conf"
          ] : ["redis-server"]
        }
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.redis_config.metadata[0].name
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
    kubernetes_namespace.redis,
    kubernetes_service.redis,
    kubernetes_config_map.redis_config
  ]
}
