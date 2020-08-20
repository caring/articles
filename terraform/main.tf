module "workspace_context" {
  source = "git@github.com:caring/tf-modules.git//aws/workspace?ref=v1.6.0"
}

provider "aws" {
  region = "us-east-1"
  version = "= 3.1.0"
  assume_role {
    role_arn = module.workspace_context.workspace_iam_roles[terraform.workspace]
  }
}

terraform {
  required_version = "= 0.12.28"

  backend "s3" {
    bucket         = "caring-tf-state"
    key            = "services/articles/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "caring-tf-state-lock"
  }
}
