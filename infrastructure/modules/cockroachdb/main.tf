data "http" "operator_crds" {
  url = "https://raw.githubusercontent.com/cockroachdb/cockroach-operator/master/install/crds.yaml"
}

data "kubectl_file_documents" "operator_crds" {
  content = data.http.operator_crds.response_body
}

resource "kubectl_manifest" "operator_crds" {
  for_each  = data.kubectl_file_documents.operator_crds.manifests
  yaml_body = each.value
}

resource "kubernetes_namespace" "operator" {
  count = var.create_namespace ? 1 : 0
  metadata {
    name = var.operator_namespace
  }
}

data "http" "operator" {
  url = "https://raw.githubusercontent.com/cockroachdb/cockroach-operator/master/install/operator.yaml"
}

data "kubectl_file_documents" "operator" {
  content = data.http.operator.response_body
}

resource "kubectl_manifest" "operator" {
  for_each  = data.kubectl_file_documents.operator.manifests
  yaml_body = each.value
  depends_on = [
    kubectl_manifest.operator_crds,
    kubernetes_namespace.operator
  ]
}

resource "kubernetes_namespace" "cockroachdb" {
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

resource "kubernetes_manifest" "cockroachdb_cluster" {
  manifest = {
    apiVersion = "crdb.cockroachlabs.com/v1alpha1"
    kind       = "CrdbCluster"
    metadata = {
      name      = var.cluster_name
      namespace = var.namespace
      labels    = var.labels
    }
    spec = {
      dataStore = {
        pvc = {
          spec = {
            accessModes = ["ReadWriteOnce"]
            resources = {
              requests = {
                storage = var.storage_size
              }
            }
            storageClassName = var.storage_class != "" ? var.storage_class : null
          }
        }
      }
      nodes = var.node_count
      image = var.image
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
      tlsEnabled = var.tls_enabled
    }
  }
  depends_on = [
    kubectl_manifest.operator
  ]
}
