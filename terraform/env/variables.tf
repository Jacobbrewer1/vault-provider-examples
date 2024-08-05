variable "vault_addr" {
  description = "Vault address"
  type        = string
  default     = "http://localhost:8200"
}

variable "vault_token" {
  description = "Vault token"
  type        = string
  default     = "root" # Set in the docker compose file
}

variable "database_username" {
  description = "Database username"
  type        = string
  default     = "root"
}

variable "database_password" {
  description = "Database password"
  type        = string
  default     = "Password01" # Set in the docker compose file
}
