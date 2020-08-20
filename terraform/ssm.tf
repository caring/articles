resource "aws_ssm_parameter" "dns" {
  name      = "/${local.service_name}/dns"
  type      = "String"
  value     = module.ecs_service.nlb_dns_name
  overwrite = true
  tags      = local.tags
}

resource "aws_secretsmanager_secret" "sentry_dsn_string" {
  name                    = "${local.service_name}_sentry_dsn"
  description             = "The DSN of our sentry connection for this service."
  recovery_window_in_days = 0
  tags                    = local.tags
}

resource "aws_secretsmanager_secret_version" "sentry_dsn_version" {
  secret_id      = aws_secretsmanager_secret.sentry_dsn_string.id
  secret_string  = "<sentry_dsn_string>"
  version_stages = [
    module.workspace_context.workspace_deploy_env[ terraform.workspace ]
  ]
  depends_on     = [
    aws_secretsmanager_secret.sentry_dsn_string
  ]
}

data "aws_secretsmanager_secret_version" "github_user" {
  secret_id     = data.terraform_remote_state.secrets.outputs.github_user_string_arn
  version_stage = module.workspace_context.workspace_deploy_env[ terraform.workspace ]
}

data "aws_secretsmanager_secret_version" "github_token" {
  secret_id     = data.terraform_remote_state.secrets.outputs.github_token_string_arn
  version_stage = module.workspace_context.workspace_deploy_env[ terraform.workspace ]
}


resource "aws_secretsmanager_secret" "db_migration_src_string" {
  name                    = "${local.service_name}_db_migration_src"
  description             = "The source location for the DB migration scripts"
  recovery_window_in_days = 0
  tags                    = local.tags
}

resource "aws_secretsmanager_secret_version" "db_migration_src_version" {
  secret_id      = aws_secretsmanager_secret.db_migration_src_string.id
  secret_string  = "github://${data.aws_secretsmanager_secret_version.github_user.secret_string}:${data.aws_secretsmanager_secret_version.github_token.secret_string}@caring/${local.service_name}/internal/db/migrations"
  version_stages = [
    module.workspace_context.workspace_deploy_env[ terraform.workspace ]
  ]
  depends_on     = [
    data.aws_secretsmanager_secret_version.github_user,
    data.aws_secretsmanager_secret_version.github_token
  ]
}


data "aws_secretsmanager_secret_version" "acm_ssl_cert" {
  secret_id = data.terraform_remote_state.secrets.outputs.acm_ssl_cert_arn
}
