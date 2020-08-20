module "ssl_cert" {
  source            = "git@github.com:caring/tf-modules.git//aws/acm?ref=v1.6.0"
  provider_iam_role = module.workspace_context.workspace_iam_roles["caring-main"]
  hosted_zone       = "caring.com"
  domain_name       = "${local.dns_record[terraform.workspace]}.caring.com"
  tags              = local.tags
}
