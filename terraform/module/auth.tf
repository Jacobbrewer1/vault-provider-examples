# Create a new userpass auth method for applications to authenticate with
resource "vault_auth_backend" "userpass" {
  type = "userpass"
}

resource "vault_generic_endpoint" "auth_userpass" {
  for_each = {for app in var.auth_applications : app.application => app}
  depends_on           = [vault_auth_backend.userpass]
  path                 = "auth/userpass/users/${each.value.application}"
  ignore_absent_fields = true
  data_json = <<EOT
{
  "password": "${each.value.password}",
  "policies": ${jsonencode(each.value.policies)}
}
EOT
}
