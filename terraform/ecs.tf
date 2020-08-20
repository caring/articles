resource "aws_cloudwatch_log_group" "service_log_group" {
  name              = "${local.service_name}_log_group"
  retention_in_days = 30

  tags = local.tags
}

module ecs_service {
  source = "git@github.com:caring/tf-modules.git//aws/ecs_service_with_nlb?ref=v1.6.0"
  vpc_id = data.terraform_remote_state.network.outputs.vpc_id
  tags   = local.tags

  ##################################
  # load balancer variables
  ##################################
  internal            = "false"
  nlb_listener_port   = var.ecs_task_port
  # nlb_protocol          = "TCP" (Default)
  listener_protocl    = "TLS"
  nlb_subnets         = data.terraform_remote_state.network.outputs.vpc_public_subnet_ids
  nlb_security_groups = [
    data.terraform_remote_state.network.outputs.nlb_sg_id,
    data.terraform_remote_state.network.outputs.ecs_sg_id
  ]
  ssl_certificate_arn = module.ssl_cert.valid_cert_arn

  ##################################
  # service variables
  ##################################
  service_name       = local.service_name
  ecs_cluster_id     = data.terraform_remote_state.network.outputs.ecs_cluster_id
  container_port     = var.ecs_task_port
  container_count    = 2
  # iam_role              = ecs_iam_role_arn
  #health_check_grace_period_seconds = (using default)
  #cpu                   = (using default)
  #memory                = (using default)
  task_definition    = data.template_file.task_definition.rendered
  task_role_arn      = aws_iam_role.ecs_task_role.arn
  execution_role_arn = aws_iam_role.ecs_task_role.arn

  ##################################
  # logging variables
  ##################################
  log_retention = 30
}

data "aws_secretsmanager_secret" "sentry_dsn" {
  name       = "${local.service_name}_sentry_dsn"
  depends_on = [
    aws_secretsmanager_secret.sentry_dsn_string
  ]
}

data "aws_secretsmanager_secret" "rds_db_pass" {
  name       = "${local.service_name}_rds_db_pass"
  depends_on = [
    aws_secretsmanager_secret.rds_db_pass
  ]
}

data "aws_secretsmanager_secret" "db_migration_src" {
  name       = "${local.service_name}_db_migration_src"
  depends_on = [
    aws_secretsmanager_secret.db_migration_src_string
  ]
}


data template_file "task_definition" {
  template = file("${path.module}/templates/task-definition.json")

  vars = {
    ##########################
    # environment information
    ##########################
    aws_default_region = var.aws_region
    environment_name   = local.deploy_env_name
    subdomain          = local.subdomain

    #########################
    # container info
    #########################
    image_url      = "${aws_ecr_repository.ecr.repository_url}:latest"
    container_name = local.service_name
    container_port = var.ecs_task_port

    ########################
    # Logging (ECS)
    #######################
    log_group_region = var.aws_region
    log_group_name   = aws_cloudwatch_log_group.service_log_group.name

    #########################################
    # below are env vars used by the service
    #########################################

    #####################
    # App
    ####################
    service_name = local.service_name
    port         = var.ecs_task_port

    
    ####################
    # DB
    ####################
    db_host           = module.rds_db.rds_instance_address
    db_port           = "3306"
    db_user           = local.service_name
    db_pwd            = data.aws_secretsmanager_secret.rds_db_pass.arn
    db_schema         = local.service_name
    db_migrations_src = data.aws_secretsmanager_secret.db_migration_src.arn
    waitfordbhost     = module.rds_db.rds_instance_address
    

    #########################
    # Logging (Application)
    #########################
    log_name              = local.service_name
    log_level             = var.log_level[terraform.workspace]
    log_enable_dev        = var.log_enable_dev[terraform.workspace]
    log_stream_monitoring = module.firehose.firehose_monitoring_stream_name
    log_stream_reporting  = module.firehose.firehose_reporting_stream_name
    log_disable_kinesis   = var.log_disable_kinesis[terraform.workspace]
    log_flush_interval    = var.log_flush_interval[terraform.workspace]
    log_buffer_size       = var.log_buffer_size[terraform.workspace]

    ##################
    # Tracing
    ##################
    trace_destination_dns  = var.trace_destination_dns[ terraform.workspace ]
    trace_destination_port = var.trace_destination_port[ terraform.workspace ]
    trace_disable          = var.trace_disable[ terraform.workspace ]
    trace_sample_rate      = var.trace_sample_rate[ terraform.workspace ]

    #################
    # sentry
    #################
    sentry_dsn     = data.aws_secretsmanager_secret.sentry_dsn.arn
    sentry_env     = local.env_name
    sentry_disable = var.sentry_disable[ terraform.workspace ]
  }
}
