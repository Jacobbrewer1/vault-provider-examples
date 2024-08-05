terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 4.3"
    }
  }
}

locals {
  env = "prod"
}

provider "vault" {
  address = var.vault_addr
  token   = var.vault_token
}

module "vault" {
  source            = "../module"
  auth_applications = [
    {
      application = "database-example"
      password    = "some-password"
      policies    = [
        vault_policy.database-policy.name,
      ]
    },
  ]
  database_host          = "mysql:3306"
  database_password      = var.database_password
  database_username      = var.database_username
  enable_database_engine = true
  env                    = local.env
}
