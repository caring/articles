
resource "random_string" "rds_db_password" {
  length           = 16
  special          = true
  override_special = "!#+"
}

resource "aws_secretsmanager_secret" "rds_db_pass" {
  name                    = "${local.service_name}_rds_db_pass"
  description             = "The password of the MYSQL RDS instance that stores permissions."
  recovery_window_in_days = 0
  tags                    = local.tags
}

resource "aws_secretsmanager_secret_version" "rds_db_pass_version" {
  secret_id      = aws_secretsmanager_secret.rds_db_pass.id
  secret_string  = random_string.rds_db_password.result
  version_stages = [
    module.workspace_context.workspace_deploy_env[ terraform.workspace ]
  ]
  depends_on     = [
    aws_secretsmanager_secret.rds_db_pass,
    random_string.rds_db_password
  ]
}

resource "aws_db_parameter_group" "stored-procedures-enabled" {
  name        = "${local.service_name}-stored-procedures-enabled"
  description = "This parameter group is just the default group slightly modified to allow stored procedures"
  family      = "mysql8.0"

  parameter {
    name  = "log_bin_trust_function_creators"
    value = "1"
  }

  tags = local.tags
}

module rds_db {
  source = "git@github.com:caring/tf-modules.git//aws/rds?ref=v1.6.0"

  rds_instance_name     = "${local.service_name}-db"
  rds_allocated_storage = var.rds_disk_size
  rds_engine_type       = "mysql"
  rds_engine_version    = "8.0.17"
  db_parameter_group    = aws_db_parameter_group.stored-procedures-enabled.id
  rds_instance_class    = var.rds_instance_class[ terraform.workspace ]
  skip_final_snapshot   = true
  database_name         = local.service_name
  database_user         = local.service_name
  database_password     = random_string.rds_db_password.result
  rds_security_group_id = data.terraform_remote_state.network.outputs.mysql_sg_id

  tags = local.tags

  subnet_az1 = data.terraform_remote_state.network.outputs.vpc_private_subnet_ids[ 0 ]
  subnet_az2 = data.terraform_remote_state.network.outputs.vpc_private_subnet_ids[ 1 ]
  subnet_az3 = data.terraform_remote_state.network.outputs.vpc_private_subnet_ids[ 2 ]
}

