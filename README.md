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

1. Run Terraform against the Vault Docker container to set up the initial configuration:

```bash
cd terraform/env && terraform init && terraform apply
```

2. Now run the example application to set up the database and the Vault secrets:

```bash
cd cmd/database && go build -o db-example && ./db-example
```

3. Now the example is running you can watch the logs to see the application running. The logs will show that the Vault
   lease is being renewed every 15 seconds. We would increase this time in a production environment.
