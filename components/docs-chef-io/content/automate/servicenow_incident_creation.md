+++
title = "ServiceNow Incident App"

draft = false

gh_repo = "automate"

[menu]
  [menu.automate]
    title = "Incident App"
    identifier = "automate/integrations/servicenow/servicenow_incident_creation.md Chef Automate Incident Creation"
    parent = "automate/integrations/servicenow"
    weight = 20
+++

The Chef Automate Incident Creation app for ServiceNow, also called the Incident App is a certified app available from the [ServiceNow](https://store.servicenow.com) store. Use this app to automate the generation of incidents in your ServiceNow Incident Management environment for Chef Client run and Chef InSpec scan failures in Chef Automate. You can integrate the data from the Incident app with ServiceNOw resolution and tracking tools, improving the effectiveness of both ServiceNow and Chef Automate applications while also helping you reduce the impact of infrastructure and compliance failures.

The Incident App exposes the REST API endpoint for communication between Chef Automate and the ServiceNow instance. Chef Automate sends HTTPS JSON notifications to the Incident App in a ServiceNow instance to creates and update incident failures.

![ServiceNow and Chef Automate Flow](/images/automate/SNOW_Automate_diagram.png)

## Key Features of Chef Automate Incident Creation App

* Incident management for infrastructure and compliance automation
* Intelligent data management and event de-duplication
* Compliance-related integration opportunities with other capabilities on ServiceNow.

## Prerequisites

* The [Integration App]({{< relref "servicenow_integration" >}}) is already installed and configured.
* Your unique ServiceNow URL. It has the format: `https://ven12345.service-now.com`.
* You need the `x_chef_automate.api` role to set up the app. Ask an Integration App or ServiceNow administrator to enable this for you.

## Installation

Get the app from the [ServiceNow](https://store.servicenow.com) store and then install it to your account from the **Service Management** dashboard.

## Setup

You can setup automatic incident creation for:

* Chef Infra Client failures
* Chef InSpec scan failures

{{% servicenow_incidents %}}

For more information about configuration, see the [ServiceNow Administrator Reference]({{< relref "servicenow_reference" >}}).

## Uninstalling

To uninstall the Chef Automate Incident App:

1. Navigate to the **System Applications** > **Applications** in ServiceNow.
1. Open the **Downloads** tab and select the **Chef Automate Incident Creation**.
1. Navigate to **Related Links**.
1. Select **Uninstall**.
