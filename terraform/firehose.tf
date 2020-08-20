module "firehose" {
  source = "git@github.com:caring/tf-modules.git//aws/firehose_monitoring_reporting?ref=v1.6.0"

  service_name              = local.service_name
  task_role_name            = aws_iam_role.ecs_task_role.name
  storage_aws_iam_user_arn  = var.storage_aws_iam_user_arn
  storage_aws_external_id   = var.storage_aws_external_id
  snowpipe_sqs_arn          = var.snowpipe_sqs_arn
  tags = local.tags
}
