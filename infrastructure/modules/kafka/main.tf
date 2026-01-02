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

        auth = {
          enabled = true
          clientProtocol = "sasl"
          interBrokerProtocol = "sasl"
          saslMechanisms = ["plain"]
          plainUsers = var.sasl_users
        }

        extraEnvVars = [
          {
            name  = "KAFKA_CFG_PROCESS_ROLES"
            value = "broker,controller"
          },
          {
            name  = "KAFKA_CFG_LISTENERS"
            value = "SASL_PLAINTEXT://:9092,CONTROLLER://:9093"
          },
          {
            name  = "KAFKA_CFG_ADVERTISED_LISTENERS"
            value = "SASL_PLAINTEXT://${var.release_name}.${var.namespace}.svc.cluster.local:9092"
          },
          {
            name  = "KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP"
            value = "CONTROLLER:PLAINTEXT,SASL_PLAINTEXT:SASL_PLAINTEXT"
          },
          {
            name  = "KAFKA_CFG_CONTROLLER_LISTENER_NAMES"
            value = "CONTROLLER"
          },
          {
            name  = "KAFKA_CFG_CONTROLLER_QUORUM_VOTERS"
            value = "0@${var.release_name}-0.${var.release_name}-headless.${var.namespace}.svc.cluster.local:9093"
          },
          {
            name  = "KAFKA_CFG_INTER_BROKER_LISTENER_NAME"
            value = "SASL_PLAINTEXT"
          },
          {
            name  = "KAFKA_CFG_SASL_ENABLED_MECHANISMS"
            value = "PLAIN"
          },
          {
            name  = "KAFKA_CFG_SASL_MECHANISM_INTER_BROKER_PROTOCOL"
            value = "PLAIN"
          },
          {
            name  = "KAFKA_CFG_SECURITY_INTER_BROKER_PROTOCOL"
            value = "SASL_PLAINTEXT"
          }
        ]

        service = {
          type = "ClusterIP"
          ports = {
            internal = 9092
          }
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
