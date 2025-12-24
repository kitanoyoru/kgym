output "cluster_name" {
  description = "Name of the CockroachDB cluster"
  value       = var.cluster_name
}

output "namespace" {
  description = "Kubernetes namespace where CockroachDB is deployed"
  value       = var.namespace
}

output "service_name" {
  description = "Name of the CockroachDB public service (created by operator)"
  value       = "${var.cluster_name}-public"
}

output "service_endpoint" {
  description = "Endpoint for connecting to CockroachDB (default port 26257)"
  value       = "${var.cluster_name}-public.${var.namespace}.svc.cluster.local:26257"
}

output "http_endpoint" {
  description = "HTTP endpoint for CockroachDB admin UI (default port 8080)"
  value       = "${var.cluster_name}-public.${var.namespace}.svc.cluster.local:8080"
}

output "node_count" {
  description = "Number of CockroachDB nodes"
  value       = var.node_count
}
