resource "vault_policy" "database-policy" {
  name   = "database-policy"
  policy = file("policies/database-policy.hcl")
}
