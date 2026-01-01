locals {
  default_kubeconfig = pathexpand("~/.kube/config")
  kubeconfig_path = var.kubeconfig_path != null && var.kubeconfig_path != "" ? var.kubeconfig_path : local.default_kubeconfig
}

provider "kubernetes" {
  config_path    = local.kubeconfig_path
  config_context = var.kubeconfig_context
}

provider "kubectl" {
  config_path    = local.kubeconfig_path
  config_context = var.kubeconfig_context
}

provider "helm" {
  kubernetes {
    config_path    = local.kubeconfig_path
    config_context = var.kubeconfig_context
  }
}
