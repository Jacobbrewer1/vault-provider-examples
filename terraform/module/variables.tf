variable "env" {
  description = "Environment"
  type        = string
}

variable "auth_applications" {
  description = "Auth application"
  type        = list(object({
    application = string
    password    = string
    policies    = list(string)
  }))
}

variable "enable_database_engine" {
  description = "Enable the database secret engine"
  type        = bool
}

variable "database_username" {
  description = "Database username"
  type        = string
  default     = "root"
}

variable "database_password" {
  description = "Database password"
  type        = string
  default     = "root"
}

variable "database_host" {
  description = "Database host"
  type        = string
  default     = "localhost:3306"
  validation {
    condition     = can(regex("^[a-zA-Z0-9.-]+:[0-9]+$", var.database_host))
    error_message = "Database host must be in the format 'hostname:port'"
  }
}
