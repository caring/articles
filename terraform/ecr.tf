resource "aws_ecr_repository" "ecr" {
  name = local.service_name
  tags = local.tags
}
