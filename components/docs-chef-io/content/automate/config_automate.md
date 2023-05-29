+++
title = "Automate Node config for HA"
draft = false
gh_repo = "automate"
[menu]
  [menu.automate]
    title = "Automate Node config for HA"
    parent = "automate/deploy_high_availability/configuration"
    identifier = "automate/deploy_high_availability/configuration/config_automate.md Automate Node config for HA"
    weight = 210
+++

{{< warning >}}
{{% automate/ha-warn %}}
{{< /warning >}}

Automate provides various configuration options that can be patched to customize its behavior and meet specific requirements. This guide outlines some out of all the available configurations that you can modify.
The config.toml file is the main configuration file for Chef Automate. It allows you to customize various aspects of the Chef Automate server. Below are some important configuration options that you can patch in the config.toml file:
To check all available automate config please run `chef-automate dev default-config`.

### Patching Automate FQDN (Fully Qualified Domain Name)

Click [here](/automate/configuration/#chef-automate-fqdn) to learn more.

### Configure Data Feed

Click [here](/automate/datafeed/#configuring-global-data-feed-behavior) to learn more.

### Auto Upgrade ON/OFF

Click [here](/automate/configuration/#upgrade-strategy) to learn more.

### Configure External Opensearch

To know about OpenSearch configuration click [here](automate/install/#configuring-external-opensearch)

### Adding resolvers for external OpenSearch
To know about adding OpenSearch resolvers click [here](automate/install/#adding-resolvers-for-opensearch)

### Backup externally-deployed OpenSearch to local filesystem


Click [here](automate/install/#backup-externally-deployed-opensearch-to-local-filesystem) for more information.

### Backup externally deployed OpenSearch to AWS S3


Click [here](automate/install/#backup-externally-deployed-opensearch-to-aws-s3) for more information.

### Backup externally deployed OpenSearch to GCS

Click [here](automate/install/#backup-externally-deployed-opensearch-to-gcs) for more information.

### Configuring External PostgresSQL

Click [here](automate/install/#configuring-an-external-postgresql-database) for more information on external PostgreSQL configuration.

### Adding resolvers for PostgreSQL database

Click [here](automate/install/#adding-resolvers-for-postgresql-database) for more information on external PostgreSQL configuration.

### Load Balancer Certificate and Private Key

Click [here](/automate/configuration/#load-balancer-certificate-and-private-key) for more information

### Proxy Settings

Click [here](/automate/configuration/#proxy-settings) for more information

### Global Log Level

Click [here](/automate/configuration/#global-log-level) for more information

### Load Balancer
Click [here](/automate/configuration/#load-balancer) for more information

### Buffer Size

Configure message buffer ingest size:

Click [here](/automate/configuration/#buffer-size) for more information

### Compliance Configuration

Click [here](/automate/configuration/#compliance-configuration) for more information

### Configure Inflight Data Collector Request Maximum

Click [here](/automate/configuration/#configure-inflight-data-collector-request-maximum) for more information

### Sign-out on Browser Close
Click [here](/automate/configuration/#sign-out-on-browser-close) for more information

### Disclosure Banner

Click [here](/automate/configuration/#disclosure-banner) for more information

### Disclosure Panel

Click [here](/automate/configuration/#disclosure-panel) for more information

### Content Security Policy Header

Click [here](/automate/configuration/#content-security-policy-header) for more information

### Session Timeout

Click [here](/automate/session_timeout/) for more information

###  Telemetry

Click [here](/automate/telemetry/) for more information

### Invalid Login Attempts

Click [here](/automate/invalid_login_attempts/) for more information

### Troubleshooting

Click [here](/automate/configuration/#troubleshooting) for more information