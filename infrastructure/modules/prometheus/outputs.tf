output "namespace" {
  description = "Kubernetes namespace where Prometheus is deployed"
  value       = var.namespace
}

output "release_name" {
  description = "Helm release name"
  value       = var.release_name
}

output "service_name" {
  description = "Name of the Prometheus service"
  value       = "${var.release_name}-kube-prometheus-prometheus"
}

output "service_endpoint" {
  description = "Endpoint for accessing Prometheus (default port 9090)"
  value       = "${var.release_name}-kube-prometheus-prometheus.${var.namespace}.svc.cluster.local:9090"
}
