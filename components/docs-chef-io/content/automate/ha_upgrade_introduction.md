+++
title = "Upgrade"

draft = false

gh_repo = "automate"
[menu]
  [menu.automate]
    title = "Upgrade"
    parent = "automate/deploy_high_availability"
    identifier = "automate/deploy_high_availability/ha_upgrade_introduction.md Upgrade HA"
    weight = 70
+++

{{< warning >}}
{{% automate/ha-warn %}}
{{< /warning >}}

Steps to upgrade the Chef Automate HA are as shown below:

- Download the latest cli 
  ```bash
   curl https://packages.chef.io/files/current/latest/chef-automate-cli/chef-automate_linux_amd64.zip | gunzip - > chef-automate && chmod +x chef-automate | cp -f chef-automate /usr/bin/chef-automate
   ```

- Download Airgapped Bundle, download latest Bundle with this:

  ```bash
  curl https://packages.chef.io/airgap_bundle/current/automate/latest.aib -o latest.aib
  ```
  Download specific version bundle with this:
  ```bash
  curl https://packages.chef.io/airgap_bundle/current/automate/<version>.aib -o automate-<version>.aib
  ```

  {{< note >}} 
  Chef Automate bundles are available for 365 days from the release of a version. However, the milestone release bundles are available for download forever.
  {{< /note >}}

- If we want to only upgrade FrontEnd Services i.e. Chef Automate and Chef Infra Server.
  ```bash
  chef-automate upgrade run --airgap-bundle latest.aib --upgrade-frontends
  ```

- If we want to only upgrade BackEnd Services i.e. Postgresql and OpenSearch.
  ```bash
  chef-automate upgrade run --airgap-bundle latest.aib --upgrade-backends
  ```

- To upgrade full Chef Automate HA System run this command from Bation Host: 
  ```bash
  chef-automate upgrade run --airgap-bundle latest.aib
  ```


{{< note >}}

  - BackEnd upgrades will restart the backend service, which take time for cluster to be in health state.
  - Upgrade command, currently only supports minor upgrade.  
{{< /note >}}

- To skip user confirmation prompt in upgrade, you can pass a flag
  ```bash 
    chef-automate upgrade run --airgap-bundle latest.aib --auto-approve
    OR 
    chef-automate upgrade run --airgap-bundle latest.aib --upgrade-backends --auto-approve
    OR
    chef-automate upgrade run --airgap-bundle latest.aib --upgrade-frontends --auto-approve
  ```

Upgrade will also check for new version of bastion workspace, if new version is available, it will promt for a confirmation for workspace upgrade before upgrading the Frontend or backend nodes, 

In case of yes, it will do workspace upgrade and no will skip this.
We can also pass a flag in upgade command to avoid prompt for workspace upgrade. 

  ```bash
   chef-automate upgrade run --airgap-bundle latest.aib --auto-approve --workspace-upgrade yes
      OR  
   chef-automate upgrade run --airgap-bundle latest.aib --auto-approve --workspace-upgrade no
  ```

{{< note >}}

  AMI Upgrade is only for AWS deployment ,as in Onpremise Deployment all the resources are managed by customers its self.  
{{< /note >}}

## AMI Upgrade Setup For AWS Deployment

### Requirements

1. Two identical clusters located in same/different AWS regions.
1. Amazon S3 access in both regions for Application backup.

In the above approach, there will be two identical clusters

- Primary Cluster (or Production Cluster)
- Secondary Cluster (Having upgraded AMI)

{{< note >}}

  The AWS deployment should be configured with S3.  
{{< /note >}}

### Steps to set up the AMI Upgraded Cluster

1. Deploy the Primary cluster following the deployment instructions by [clicking here](/automate/ha_aws_deploy_steps/#deployment).

1. Deploy the Secondary cluster into a same/different data center/region using the same steps as the Primary cluster

1. Do the backup configuration only when you have not provided the (backup information) configuration at the time of deployment. Refer backup section for [object storage](/automate/ha_backup_restore_aws_s3/).

    {{< note >}}
     During the deployment for the Primary and Secondary clusters, use the same S3 bucket name.
    {{< /note >}}

1. On Primary Cluster

    - Take a backup of Primary cluster from bastion by running below command:

    ```sh
    chef-automate backup create --no-progress > /var/log/automate-backups.log
    ```

    - Create a bootstrap bundle; this bundle captures any local credentials or secrets that aren't persisted to the database. To create the bootstrap bundle, run the following command in one of the Automate nodes:

    ```sh
    chef-automate bootstrap bundle create bootstrap.abb
    ```

    - Copy `bootstrap.abb` to all Automate and Chef Infra frontend nodes in the Secondary cluster.


1. On New AMI upgraded Cluster

    - Install `bootstrap.abb` on all the Frontend nodes (Chef-server and Automate nodes) by running the following command:

    ```cmd
    sudo chef-automate bootstrap bundle unpack bootstrap.abb
    ```

    - Run the following command in bastion to get the ID of the backups:

    ```sh
    chef-automate backup list
    ```

    -On Secondary Cluster Trigger restore command `chef-automate backup restore` on one of the Chef Automate nodes.

        - To run the restore command, you need the airgap bundle. Get the Automate HA airgap bundle from the location `/var/tmp/` in Automate instance. For example: `frontend-4.x.y.aib`.

        - In case the airgap bundle is not present at `/var/tmp`, it can be copied from the bastion node to the Automate frontend node

        - Run the command at the Automate node of Automate HA cluster to get the applied config:

        ```bash
        sudo chef-automate config show > current_config.toml
        ```

        - Add the OpenSearch credentials to the applied config. If using Chef Managed OpenSearch, add the config below into `current_config.toml` (without any changes).

        ```bash
        [global.v1.external.opensearch.auth.basic_auth]
            username = "admin"
            password = "admin"
        ```

        - In the New AMI Upgraded cluster, use the following sample command to restore the latest backup from any Chef Automate frontend instance.

        ```cmd
        id=$(chef-automate backup list | grep completed | tail -1 | awk '{print $1}')
        sudo chef-automate backup restore <backup-url-to-object-storage>/automate/$id/ --patch-config /path/to/current_config.toml --airgap-bundle /var/tmp/frontend-4.x.y.aib --skip-preflight --s3-access-key "Access_Key"  --s3-secret-key "Secret_Key"
        ```

### Switch to New Upgraded Cluster

Steps to switch to the New cluster are as follows:

- Start the services on all the Automate and Chef Infra frontend nodes, using the below command:

    ```sh
    systemctl start chef-automate
    ```
- Once the restore is successful ,you can destroy the Primary Cluster
