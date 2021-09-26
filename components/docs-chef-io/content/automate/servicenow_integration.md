+++
title = "ServiceNow Integration App"
draft = false

gh_repo = "automate"

[menu]
  [menu.automate]
    title = "Integration App"
    parent = "automate/integrations/servicenow"
    identifier = "automate/integrations/servicenow/servicenow_integration.md ServiceNow Integration App"
    weight = 10
+++

[ServiceNow](https://www.servicenow.com/) provides cloud-based services that automate enterprise IT operations. ServiceNow specializes in IT service management (ITSM) applications and provides forms-based workflow application development. It supports third-party application and data integrations. The most common integrations are for configuration management, incident management, problem management, change management, user administration, and single sign-on authentication.

The Chef Automate Integration App for ServiceNow, also called the Integration App, integrates Chef Automate with ServiceNow. This app is a ServiceNow-certified scoped application available from the ServiceNow store. It integrates existing Chef Automate infrastructure and compliance functionality with ServiceNow enterprise services. Once installed and configured, this integration enables Chef Automate to create and update a ServiceNow Configuration Management Database (CMDB) with data from nodes managed by Chef Automate. Chef Automate aggregates information about infrastructure nodes, the Chef Infra Client runs, and Chef Compliance scans, helping you observe your infrastructure in real time.

The Integration App works by exposing the REST API endpoints for communication between Chef Automate and ServiceNow. Chef Automate sends HTTPS JSON data feeds to the app in ServiceNow to create and update the _ServiceNow CMDB_, client run, and compliance report tables.

{{< figure src="/images/automate/snow_integration_dataflow_diagram.png" alt="Data Flow Diagram" >}}

## Prequisites

### User Requirements

- Your unique ServiceNow URL. It has the format: `https://ven12345.service-now.com`.
- Once the Integration app is installed, you will need to have the `x_chef_automate.api` role to configure it. Your ServiceNow administrator can enable this for you.

### System and User Requirements

- A running [Chef Automate](https://www.chef.io/automate/) instance.
- Chef Automate is provisioned with a valid SSL/TLS certificate from a trusted certificate authority (CA).
- A running [ServiceNow](https://www.servicenow.com/) instance.
- The ServiceNow is reachable on port 443.

### Required ServiceNow Plugins

Install following ServiceNow plugins from the Service Management dashboard:

- System Import Sets `com.glide.system_import_set`, min version 1.0.0.
- Configuration Management (CMDB) 1.1.
- Configuration Management for Scoped Apps (com.snc.cmdb.scoped) 1.0.0.

{{< figure src="/images/automate/snow_integration_plugins.png" alt="Plugins" >}}

## Install the Integration App

Install the Integration App from the [ServiceNow Store](https://store.servicenow.com)

1. Navigate to the ServiceNow store at <https://store.servicenow.com>.
2. Search for **Chef Automate**.
3. Select the **Chef Automate Integration App**.
4. Select **Get** and follow the instructions by specifying your ServiceNow credentials.
5. Open your **ServiceNow Service Management** application.
6. Select **System Applications** > **All Available Applications** > **All** menu.
7. Find the application using the filter criteria and search bar.
   You can search for the application by its `Name` or `ID`. If you cannot find an application, you may have to request it from the ServiceNow Store.
8. Select **Install**.

## Setup

Connect the Integration App in ServiceNow to Chef Automate by creating a data feed in from Chef Automate or a connection in Service Now. Setting up and configuring the Integration App requires the `x_chef_automate.api` role to configure it. If you are not a ServiceNow administrator, ask one to enable it for you.

### Create a Data Feed in Chef Automate

Set up a data feed to send data from Chef Automate to the Integration App:

1. Confirm that you have the `x_chef_automate.api` role.
1. Navigate to **Settings** from the **Chef Automate** menu.
1. Select the **Data Feeds** link from **Settings**.
1. Select the **Create Data Feed** button.
1. Enter the following information:

   - **Name**: A unique name for this integration.
   - **Data Feed URL**: ServiceNow application Datafeed API URL ending with `api/x_chef_automate/asset`.
   - **Username**: The ServiceNow user.
   - **Password**: The Chef Automate user password.

   You must specify the fully qualified domain name (FQDN) of the ServiceNow instance to configure the **Data Feed URL** in Chef Automate, which has the format: `FQDN/api/x_chef_automate/asset` . For example, <https://venxxx.service-now.com/api/x_chef_automate/asset>.

1. Select the **Test Data Feed** button. The application checks that the values specified are correct and the connectivity between Chef Automate and the application establishes.

   A successful test displays the **service-now automate connectivity passed** message. Else, an error message displays detailing the connectivity or credentials issues and helps you in resolving the error.

1. Select **Create Data Feed**.

{{< figure src="/images/automate/snow_integration_create_data_feed.png" alt="Create Data Feed">}}

### Create a Connection from ServiceNow

Follow these steps connect the Integration App to Chef Automate:

1. Confirm that you have the `x_chef_automate.api` role.
1. Navigate to the **ServiceNow** instance.
1. Locate **Chef Automate** from the search bar.
1. Select the **Automate Instances** module.
1. Select the **New** button.
1. Enter the following details:

   - **Name**: A unique name for this integration.
   - **Instance URL**: The Chef Automate URL.
   - **Automate API token**: A Chef Automate API token with data-feed-service authorization.
   - **ServiceNow user**: The ServiceNow user name.
   - **ServiceNow password**: The ServiceNow user password.

1. Select the **Test Connectivity** button. The Integration App checks that the values are correct and test the connection with Chef Automate.

   - A successful test displays a **service-now automate connectivity passed** message.
   - An unsuccessful test results in an error message with suggestions for resolving the connection or credentials problems.

1. Select **Submit**. The ServiceNow app creates the Chef Automate instance detail.

{{< figure src="/images/automate/snow_integration_configure_sn.png" alt="Chef Automate Instances">}}

## Integration App Overview

When you install and setup the Integration App, the **Chef Automate** appears in the ServiceNow navigation pane.

{{< figure src="/images/automate/snow_integration_navigation.png" alt="Navigation">}}

- **Automate Instances** module allows the user to configure the integration with Chef Automate.
- **Server** module displays a list of servers in the **CMDB** module.

{{< figure src="/images/automate/snow_integration_discovery.png" alt="Discovery Source">}}

### Server Module

Integration App updates the CMDB file systems and software tables. It adds associated data from Chef Automate to the servers for:

- Client Runs
- Attributes
- Compliance Reports

#### Discovery Source

The Integration App augments the existing CMDB servers and inserts new servers into the ServiceNow CMDB. The Integration App uses the ServiceNow Discovery Identification and Reconciliation Engine (IRE)to insert or update servers.

The Integration App distinguishes between discovered Configuration Items (CIs) and imported CIs. CIs discovered by the Integration App have a default value of **Chef Automate** in the **Discovery Source** field in the `cmdb_servers` table.

#### Client Runs

You can view the data for a Chef Infra Client run by selecting a server from the list. The **Client Run** record displays information for:

- Client run cookbooks: cookbooks executed during the Chef Infra Client run.
- Client run lists: run lists executed during the Chef Infra Client run.
- Client run recipes: recipes executed during the Chef Infra Client run.

The complete Chef Infra Client run details are available for each server. Chef Infra Client runs are available from the **Client Runs** module.

#### Attributes

You can view the current attributes for a server by selecting the attributes record. All Ohai attributes are available for each server.

#### Compliance Reports

You can view the **Compliance** report detail from a server record by selecting the name on an individual **Compliance** report. The Compliance Report record displays related information for:

- Compliance report profiles: all profiles executed during the compliance scan.
- Compliance report results: all results from the compliance scan.

Also, you can view the **Compliance** report profile to display the individual results for each profile. The entire Chef Compliance Report details are available for each server. Compliance reports are also available from the **Compliance** report module.

For more information about configuration, see the [ServiceNow Administrator Reference]({{< relref "servicenow_reference" >}}).

## Uninstalling

To uninstall the Integration App:

1. Navigate to the **System Applications** > **Applications** in ServiceNow.
1. Open the **Downloads** tab and select the **Chef Automate Incident Creation**.
1. Navigate to **Related Links**.
1. Select **Uninstall**.
