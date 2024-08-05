resource "vault_mount" "db" {
  count = var.enable_database_engine ? 1 : 0 # Only create if enabled
  path  = "org-${var.env}-inst-db"
  type  = "database"
}

resource "vault_database_secret_backend_connection" "mysql" {
  count         = var.enable_database_engine ? 1 : 0 # Only create if enabled
  depends_on    = [vault_mount.db]
  backend       = vault_mount.db[0].path
  name          = "org-${var.env}-inst-mysql"
  allowed_roles = [
    "readwrite",
  ]

  mysql {
    max_open_connections = 5
    connection_url       = "${var.database_username}:${var.database_password}@tcp(${var.database_host})/"
  }
}

resource "vault_database_secret_backend_role" "role" {
  count               = var.enable_database_engine ? 1 : 0 # Only create if enabled
  depends_on          = [vault_mount.db, vault_database_secret_backend_connection.mysql]
  backend             = vault_mount.db[0].path
  name                = "readwrite"
  db_name             = vault_database_secret_backend_connection.mysql[0].name
  creation_statements = [
    "CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';",
    "GRANT ALL PRIVILEGES ON *.* TO '{{name}}'@'%';"
  ]
  default_ttl = "5" # 5 seconds
  max_ttl     = "15" # 15 seconds
}
