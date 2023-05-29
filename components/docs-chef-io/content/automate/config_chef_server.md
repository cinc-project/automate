+++
title = "Chef Server Node config for HA"
draft = false
gh_repo = "automate"
[menu]
  [menu.automate]
    title = "Chef Server Node config for HA"
    parent = "automate/deploy_high_availability/configuration"
    identifier = "automate/deploy_high_availability/configuration/config_chef_server.md Chef Server Node config for HA"
    weight = 210
+++

{{< warning >}}
{{% automate/ha-warn %}}
{{< /warning >}}

## Configure chef-server to send data to Automate

```toml
[global.v1.external.automate]
enable = true
node = "https://<automate server url>"
[global.v1.external.automate.auth]
token = "<data-collector token>"
[global.v1.external.automate.ssl]
server_name = "<server name from the automate server ssl cert>"
root_cert = """<pem format root CA cert>
"""
[auth_n.v1.sys.service]
# It is fine to use an A2 data collector token.
a1_data_collector_token = "<data-collector token>"
[erchef.v1.sys.data_collector]
enabled = true
```

### Chef Infra Configuration In Chef Automate

Click [here](/automate/chef_infra_in_chef_automate) for more information

## Patching Automate FQDN (Fully Qualified Domain Name)

Click [here](/automate/configuration/#chef-automate-fqdn) to learn more.

## Auto Upgrade ON/OFF

Click [here](/automate/configuration/#upgrade-strategy) to learn more.

## Configure External Opensearch

To know about OpenSearch configuration click [here](automate/install/#configuring-external-opensearch)

## Adding resolvers for external OpenSearch
To know about adding OpenSearch resolvers click [here](automate/install/#adding-resolvers-for-opensearch)

## Backup externally deployed OpenSearch to local filesystem

Click [here](automate/install/#backup-externally-deployed-opensearch-to-local-filesystem) for more information.

## Backup externally deployed OpenSearch to AWS S3

Click [here](automate/install/#backup-externally-deployed-opensearch-to-aws-s3) for more information.

## Backup externally deployed OpenSearch TO GCS

Click [here](automate/install/#backup-externally-deployed-opensearch-to-gcs) for more information.

## Configuring External PostgresSQL

Click [here](automate/install/#configuring-an-external-postgresql-database) for more information on external PostgreSQL configuration.

## Adding resolvers for PostgreSQl database

Click [here](automate/install/#adding-resolvers-for-postgresql-database) for more information on external PostgreSQL configuration.

#### Load Balancer Certificate and Private Key

Click [here](/automate/configuration/#load-balancer-certificate-and-private-key) for more information

#### Proxy Settings

Click [here](/automate/configuration/#proxy-settings) for more information

#### Global Log Level

Click [here](/automate/configuration/#global-log-level) for more information

#### Load Balancer
Click [here](/automate/configuration/#load-balancer) for more information

### Troubleshooting

Click [here](/automate/configuration/#troubleshooting) for more information