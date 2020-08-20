data "aws_iam_policy_document" "ecs_assume_role" {
  statement {
    actions = [
      "sts:AssumeRole"
    ]

    principals {
      type        = "Service"
      identifiers = [
        "ecs-tasks.amazonaws.com",
        "firehose.amazonaws.com"
      ]
    }
  }
}

resource "aws_iam_role" "ecs_task_role" {
  name               = "${local.service_name}-task-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_assume_role.json
  tags               = local.tags
  depends_on         = [
    data.aws_iam_policy_document.ecs_assume_role
  ]
}

data "aws_iam_policy_document" "ecs_task_role" {
  statement {
    sid       = "AuthSecretsSecretsAccess"
    effect    = "Allow"
    actions   = [
      "secretsmanager:GetSecretValue"
    ]
    resources = [
      
      "arn:aws:secretsmanager:${var.aws_region}:${data.aws_caller_identity.current.account_id}:secret:${local.service_name}_rds_db_pass-??????",
      "arn:aws:secretsmanager:${var.aws_region}:${data.aws_caller_identity.current.account_id}:secret:${local.service_name}_db_migration_src-??????",
      "arn:aws:secretsmanager:${var.aws_region}:${data.aws_caller_identity.current.account_id}:secret:${local.service_name}_sentry_dsn-??????"
      
    ]
  }
}

resource "aws_iam_policy" "ecs_task_role" {
  name       = local.service_name
  policy     = data.aws_iam_policy_document.ecs_task_role.json
  depends_on = [
    data.aws_iam_policy_document.ecs_task_role
  ]
}

resource "aws_iam_policy" "internal_profile_policy_copy" {
  name   = "${local.service_name}-profile-internal"
  policy = data.terraform_remote_state.network.outputs.tf_profile_internal_policy_document
}

resource "aws_iam_role_policy_attachment" "ecs_task_role_policy_attach" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.ecs_task_role.arn
  depends_on = [
    aws_iam_role.ecs_task_role,
    aws_iam_policy.ecs_task_role
  ]
}

resource "aws_iam_role_policy_attachment" "profile_internal_policy_attach" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.internal_profile_policy_copy.arn
  depends_on = [
    aws_iam_role.ecs_task_role,
    aws_iam_policy.internal_profile_policy_copy
  ]
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_ecr_policy" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  depends_on = [
    aws_iam_role.ecs_task_role
  ]
}
