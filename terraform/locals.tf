locals {
  service_name    = "articles"
  deploy_env_name = module.workspace_context.workspace_deploy_env[terraform.workspace]
  subdomain       = module.workspace_context.env_subdomain[terraform.workspace]
  env_name        = module.workspace_context.env_name[terraform.workspace]

  dns_record = {
    "caring-prod": local.service_name
    "caring-stg": "${local.service_name}.${local.subdomain}"
    "caring-dev": "${local.service_name}.${local.subdomain}"
  }

  tags = {
    service_name = local.service_name,
    env          = terraform.workspace,
    terraform    = "true",
    repo_name    = local.service_name
  }
}
