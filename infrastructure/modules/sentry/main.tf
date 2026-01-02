resource "kubernetes_namespace" "sentry" {
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

resource "kubernetes_secret" "sentry_secret_key" {
  metadata {
    name      = "${var.release_name}-secret-key"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    secret-key = base64encode(random_password.sentry_secret_key.result)
  }
  depends_on = [kubernetes_namespace.sentry]
}

resource "random_password" "sentry_secret_key" {
  length  = 50
  special = true
}

resource "kubernetes_secret" "sentry_postgresql" {
  count = var.postgresql_host != "" ? 1 : 0
  metadata {
    name      = "${var.release_name}-postgresql"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    host     = base64encode(var.postgresql_host)
    port     = base64encode(tostring(var.postgresql_port))
    database = base64encode(var.postgresql_database)
    username = base64encode(var.postgresql_user)
    password = base64encode(var.postgresql_password != "" ? var.postgresql_password : "")
  }
  depends_on = [kubernetes_namespace.sentry]
}

resource "kubernetes_secret" "sentry_redis" {
  count = var.redis_host != "" ? 1 : 0
  metadata {
    name      = "${var.release_name}-redis"
    namespace = var.namespace
    labels    = var.labels
  }
  data = var.redis_password != null ? {
    host     = base64encode(var.redis_host)
    port     = base64encode(tostring(var.redis_port))
    password = base64encode(var.redis_password)
  } : {
    host = base64encode(var.redis_host)
    port = base64encode(tostring(var.redis_port))
  }
  depends_on = [kubernetes_namespace.sentry]
}

resource "helm_release" "sentry" {
  name       = var.release_name
  repository = "https://sentry-kubernetes.github.io/charts"
  chart      = "sentry"
  version    = var.chart_version != "" ? var.chart_version : null
  namespace  = var.namespace
  timeout    = var.timeout
  wait       = true

  values = [
    yamlencode({
      user = {
        email    = var.user_email
        password = var.user_password
      }
      postgresql = var.postgresql_host != "" ? {
        enabled = false
        persistence = null
      } : {
        enabled = true
        persistence = {
          enabled      = true
          size         = var.storage_size
          storageClass = var.storage_class != "" ? var.storage_class : null
        }
      }
      externalPostgresql = var.postgresql_host != "" ? {
        host     = var.postgresql_host
        port     = var.postgresql_port
        database = var.postgresql_database
        username = var.postgresql_user
        password = var.postgresql_password != "" ? var.postgresql_password : ""
      } : null
      redis = var.redis_host != "" ? {
        enabled = false
        master = null
      } : {
        enabled = true
        master = {
          persistence = {
            enabled      = true
            size         = var.storage_size
            storageClass = var.storage_class != "" ? var.storage_class : null
          }
        }
      }
      externalRedis = var.redis_host != "" ? {
        host     = var.redis_host
        port     = var.redis_port
        password = var.redis_password != null ? var.redis_password : null
      } : null
      rabbitmq = {
        enabled = false
      }
      kafka = var.kafka_host != "" ? {
        enabled = false
      } : {
        enabled = true
      }
      externalKafka = var.kafka_host != "" ? {
        host = var.kafka_host
        port = var.kafka_port
      } : null
      web = {
        replicas = 1
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
      worker = {
        replicas = 1
        resources = {
          requests = {
            cpu    = var.worker_resources.requests.cpu
            memory = var.worker_resources.requests.memory
          }
          limits = {
            cpu    = var.worker_resources.limits.cpu
            memory = var.worker_resources.limits.memory
          }
        }
      }
      cron = {
        replicas = 1
        resources = {
          requests = {
            cpu    = var.worker_resources.requests.cpu
            memory = var.worker_resources.requests.memory
          }
          limits = {
            cpu    = var.worker_resources.limits.cpu
            memory = var.worker_resources.limits.memory
          }
        }
      }
      ingress = {
        enabled = false
      }
      service = {
        type = "ClusterIP"
      }
    })
  ]

  depends_on = [
    kubernetes_namespace.sentry,
    kubernetes_secret.sentry_secret_key
  ]
}
