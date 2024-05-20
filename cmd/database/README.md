# Database provider

In this example you can find the example of how to interact with the Vault database provider.

## Configuration

The app will require some config to be able to connect to the database. This could be provided to the app however you
wish. For this example I have provided it with a `config.json` file.

```json
{
  "vault": {
    "address": "https://vaulthost.com",
    "app_role_id": "some-app-role-id",
    "app_role_secret_id": "some-app-role-secret-id"
  },
  "database": {
    "credentials_path": "database/creds/my-mysql-role",
    "host": "localhost:3306",
    "schema": "schema"
  }
}
```

In the config file you can see that the `vault` section contains the address of the Vault server and the app role id and
secret id. The `database` section contains the path to the database credentials in Vault, the host of the database and
the schema to connect to.

## Vault

What Vault does under the hood is it will generate a new set of credentials for the database and return them to the
application. The application can then use these credentials to connect to the database.

## Common things

A lot of people will get stuck when they receive the unauthorised error when trying to connect to the database. A common
reason for this is because they have not aligned the app role entity with the MySQL role in Vault. This can be done
through a policy. Please look at the vault documentation for more information on this.

## Custom Vault Package

What my custom vault wrapper package does is it will allow you to interact with Vault in a more object-oriented way. It
will allow you to create a new Vault object and then call the method to get the credentials. This will return a new
`Credentials` object which will contain the username, password and the lease duration. You can then call the renewal
method on the `Credentials` object to renew the lease. You will have to pass an `renewal` method to the method; what
happens is the `renewal` method will be called every time the lease is expired and new credentials are generated.

