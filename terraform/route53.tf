module "service_dns_record" {
  source            = "git@github.com:caring/tf-modules.git//aws/route53_record_set?ref=v1.6.0"
  provider_iam_role = module.workspace_context.workspace_iam_roles["caring-main"]
  record_name       = "${local.dns_record[terraform.workspace]}.caring.com"
  records           = [module.ecs_service.nlb_dns_name]
  tags              = local.tags
}
