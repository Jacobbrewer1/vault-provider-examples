# Vault Provider Examples

## Setup

To run the examples you will need to have a Vault server running as well as the MySQL database. I have included
a `docker-compose.yml` file which will start the Vault server and the MySQL database. To start the services run the
following command:

```bash
docker-compose up -d
```

Now that the services are running, we need to set up the configuration on Vault. To do this we need to run the following
steps:

1. Enable the database secrets engine (This can be done through the UI or CLI)
    1. Vault is running at `http://localhost:8200` with the token of `root`
2. Create a new MySQL connection
    1. `Database plugin` will be the `mysql-database-plugin`
    2. `Connection name` can be anything you want
    3. `Connection will be verified` should be checked
    4. `Use custom password policy` should be unchecked
    5. `Connection URL` should be `{{username}}:{{password}}@tcp(mysql:3306)/`
        1. You can replace the `mysql` with the IP of the docker container too. The IP can be found by using the docker
           inspect command
    6. `Username` should be `root`. In production, you should create a vault user that has root permissions and then
       rotate the root password to something else. As this is a demo, we will use the root user
    7. `Password` should be `Password01`
    8. We can then leave the rest of the fields as default
3. Create a new role
    1. `Role name` can be anything you want
    2. `Connection name` should be the connection name you created in step 2
    3. `Type of role` should be `dynamic`
    4. `Generated credentials’s Time-to-Live (TTL)` should be set to `10s`. This is for the demo purposes. In
       production, you should set this to a higher value
    5. `Generated credentials’s maximum Time-to-Live (Max TTL)` should be set to `30s`. This is for the demo purposes.
       In production, you should set this to a higher value
    6. `Creaation statements` should
       be `CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT ALL PRIVILEGES ON *.* TO '{{name}}'@'%';`
    7. `Revocation statements` should be left as default
4. Enable the [approle auth method](https://developer.hashicorp.com/vault/docs/auth/approle)
    1. We will do this through the CLI; you can do this through the UI if you wish
    2. Get exec into the vault container
    3. Run the following command to enable the approle auth method `vault auth enable approle`
    4. Run the following command to create a named
       role `vault write auth/approle/role/my-role token_type=batch secret_id_ttl=365d token_ttl=365d token_max_ttl=365d`.
       This will create a new role called `my-role` with the token type of `batch` and the secret id ttl, token ttl and
       token max ttl of `365d`. This is for the demo purposes. In production, you should set this to a lower value and
       figure out a way of rotating the secret id without any downtime to the application.
    5. Fetch the role id of the new AppRole `vault read auth/approle/role/my-role/role-id`. This will return the role id
    6. Get a secret id issued against the AppRole `vault write -f auth/approle/role/my-role/secret-id`. This will return
       the secret id. You should store this in a secure location as you will not be able to retrieve it again
5. Go to the UI and create a new policy
    1. The policy should be able to read the database credentials. The policy should look like this:
       ```hcl
       path "database/creds/{{ the role name you created in step 3 }}" {
         capabilities = ["read"]
       }
       ```
6. Attach the policy to the AppRole
    1. In the UI go to Access
    2. Click on `approle`
    3. Find the tab named `Entities`
    4. There will now be an entity there, click on the entity
    5. Click the tab named `Policies`
    6. Click on `Edit entity`
    7. Under policies, you can now search to add the policy you created in step 5
    8. Click save
7. Create a new config file for the demo app
    1. Create a new file called `config.json` in the root of the project
    2. Add the following to the file:
       ```json
       {
         "vault": {
           "address": "http://localhost:8200",
           "app_role_id": "{{ the role id you got in step 4.5 }}",
           "app_role_secret_id": "{{ the secret id you got in step 4.6 }}"
         },
         "database": {
           "credentials_path": "database/creds/{{ the role name you created in step 3 }}",
           "host": "localhost:3306",
           "schema": "schema"
         }
       }
       ```
8. Run the demo app and watch the logs over the next 2 minutes. You should see the app renewing the lease every 10
   seconds
