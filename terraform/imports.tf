#####################################################
#################### Caller Info ####################
#####################################################
data "aws_caller_identity" "current" {}

#####################################################
############## Core Network Information #############
#####################################################
data "terraform_remote_state" "network" {
  backend   = "s3"
  workspace = terraform.workspace
  config    = {
    bucket  = "caring-tf-state"
    key     = "tf-core/base/terraform.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}

#####################################################
################ Secrets Information ################
#####################################################
data "terraform_remote_state" "secrets" {
  backend   = "s3"
  workspace = terraform.workspace
  config    = {
    bucket  = "caring-tf-state"
    key     = "tf-core/secrets/terraform.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}
