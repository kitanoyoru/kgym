output "namespace" {
  description = "Kubernetes namespace where Kafka is deployed"
  value       = var.namespace
}

output "release_name" {
  description = "Helm release name"
  value       = var.release_name
}

output "service_name" {
  description = "Name of the Kafka service"
  value       = "${var.release_name}"
}

output "service_endpoint" {
  description = "Endpoint for accessing Kafka (default port 9092)"
  value       = "${var.release_name}.${var.namespace}.svc.cluster.local:9092"
}


output "replicas" {
  description = "Number of Kafka broker replicas"
  value       = var.replicas
}
