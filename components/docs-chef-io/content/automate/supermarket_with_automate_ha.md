+++
title = "Supermarket with Automate HA"
draft = false
gh_repo = "automate"
[menu]
  [menu.automate]
    title = "Supermarket with Automate HA"
    parent = "automate/deploy_high_availability"
    identifier = "automate/deploy_high_availability/supermarket_with_automate_ha.md Supermarket with Automate HA"
    weight = 75
+++
Chef Supermarket is the site for community cookbooks. It provides a searchable cookbook repository and a friendly web UI. This page will discuss the integration of Supermarket with Automate HA setup using On-Premises along with AWS Managed Services, OpenSearch, and PostgreSQL.
Before integrating the supermarket, let's see how to register a new Oauth application with OC-ID under Automate.

## Register a New Oauth Application with OC-ID

When you install Chef Automate, it bundles the OC-ID component as an Oauth provider. The Oauth provider uses the Chef-Server credentials to log in to another application—for example, Supermarket. Follow the steps below to register the applications to use OC-ID as a medium to log in to the respective applications. Once you finish the registration, you will be authorized to use the application if the Che-Server login credentials are correct.

1. Create a file to list down the details of the application you want to register with OC-ID. In the file, mention the **name** and the **redirect_uri** for the file. The content of the created file looks like as shown below:

    ```cd
    [ocid.v1.sys.ocid.oauth_application_config]
        [[ocid.v1.sys.ocid.oauth_application_config.oauth_applications]]
            name = ""
            redirect_uri = ""
    ```

    To add more than one application to the OC-ID, keep adding the above code in the file. For example:

    ```cd
    [ocid.v1.sys.ocid.oauth_application_config]
        [[ocid.v1.sys.ocid.oauth_application_config.oauth_applications]]
            name = "application-1"
            redirect_uri = "https://application-1.com/auth/chef_oauth2/callback"
        [[ocid.v1.sys.ocid.oauth_application_config.oauth_applications]]
            name = "application-2"
            redirect_uri = "https://application-2.com/auth/chef_oauth2/callback"
    ```

    Using the above snippet, you can register two applications to the OC-ID.

1. Save the file in `.toml` format. For example: `register_oauth_application.toml`.
1. Now, patch the above file by running the below command from your habitat studio `/src` directory or your bastion system:

    ```bash
    chef-automate config patch register_oauth_application.toml
    ```

1. Verify whether the new configuration has been applied or not by running the following command:

    ```bash
    chef-automate config show
    ```

    The output of the above command will contain the values from the file you patched by running the patch command.

1. The above command will patch the new configuration and restart the OC-ID service. Under the hood, it does a few things as follows:

    * Once you patch the `.toml` file, it will update the values in the running automate configuration mentioned in the `.toml` file.

    * It will restart the **automate-cs-ocid** service, invoking that component's run hook.

    * The run hook has a default instruction to register a new application if mentioned in the automate configuration.

    * As part of the registration process, it creates a record for the application and generates a pair of (**uid+secret**) in the OC-ID database.

    * Now, your app is successfully registered as an oauth application under OC-ID service.

    * Then we have one more instruction in the run hook, i.e., to dump the generated configuration of all the registered oauth applications as a `.yaml` file in the path: `/hab/svc/automate-cs-ocid/data/registered_oauth_applications.yaml`. This file is the source of truth of the data of all registered oauth applications in that instance of automate.

1. Run the following `ctl` command to get the details of the registered oauth application created above.

    ```cd
    chef-automate config oc-id-show-app
    ```

    The output of the above command is as shown below:

    ```cd
    supermarket1:
    -   name: supermarket1
        redirect_uri: https://127.0.0.1:3001/auth/chef_oauth2/callback
        uid: 735c44e423787134839ce1bdb6b2ab8bd9eca5b656f0f4e69df3641ea494cdda
        secret: 4c371ceb46465b162c0b4a670573d80ac1d6adeebaa2638db53bb9f94d432340
        id:
    supermarket2:
    -   name: supermarket2
        redirect_uri: https://127.0.0.1:3002/auth/chef_oauth2/callback
        uid: 4e3339dc860ad9f499624034f31666c6f737ece10d40ad6ce0b8819efeeb52b0
        secret: 96c39da744af78ff16d2227d5701a18efb98e25b62baaf395b3b27efd280e37c
        id:
    supermarket3:
    -   name: supermarket3
        redirect_uri: https://127.0.0.1:3003/auth/chef_oauth2/callback
        uid: 0d35a770941cb9747418b9c9de0de76b7d3c7964930b55ae0ce6e4f930377df0
        secret: 5f18bb0a14f846304dc6addffce815a466c6b22446bdb48a3d1ab9acc96e6a09
        id:
    ```

1. Now, configure your application with the generated `uid` and `secret` to be authorized to use **OC-ID** as an oauth chef-server provider and use the chef-server credential to log in to the app.