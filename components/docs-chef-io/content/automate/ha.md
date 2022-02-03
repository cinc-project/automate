+++
title = "HA Overview"

draft = false

gh_repo = "automate"
[menu]
  [menu.automate]
    title = "HA Overview"
    parent = "automate/install"
    identifier = "automate/install/ha.md HA Overview"
    weight = 200
+++

## Introduction

### What is High Availability (HA)?

**High availability (HA)** refers to a system or application (such as a network, a server array, or cluster) that offers a high level of operational performance and quality over a relevant time with maximum potential uptime and accessibility for the content stored on it.

While a more basic system will be adequate to serve content to a low or medium number of users, it may include a single point of failure. This means that if one server goes down, whether due to traffic overload or any number of other issues, the entire site or application could become unavailable.

HA simply means the application remains available with no interruption. We achieve high availability when an application continues to operate when one or more underlying components fail. For example, a router, switch, firewall, or server that fails.

Thus, HA is designed to avoid loss of service by reducing or managing failures and minimizing unscheduled downtime (when your system or network is not available for use or is unresponsive) that happens due to power outages or failure of a component.

*This para required* \-- ? "Availability" includes two periods of time: how much time a service is accessible and how much time the system needs to respond to user requests. When it comes to measuring availability, several factors are salient. These include recovery time and both scheduled and unscheduled maintenance periods. Typically, availability as a whole is expressed as a percentage of uptime defined by service level agreements (SLAs). A score of 100 percent characterizes a system that never fails or experiences zero downtime by being 100% operational.

### What are High Availability (HA) clusters?

A **cluster** is a group of inter-connected computers that work together to perform intensive tasks.  In a cluster, each computer is referred to as a "node".

**HA clusters** are servers grouped to operate as a single, unified system supporting server applications that can be reliably utilized with a minimum amount of downtime. They are also called **failover clusters** as they share the same storage but use a different network and can run the same workloads of the primary system they support. HA clusters are tested regularly to confirm nodes are always in operational mode.

If a server in the cluster fails, another server or node can take over immediately to help ensure the application or service supported by the cluster remains operational. HA clusters help ensure no single point of failure for critical IT and reduce or eliminate downtime.

Thus, HA clusters strive to support the system or application or services run reliably with minimal downtime.

## Chef Automate High Availability (HA)

### What does Chef Automate HA Brings to you?

The Chef Automate HA equates to reliability, intending to increase functionality, efficiency, and productivity. It is built on the following three characteristics, **Redundancy**, **Monitoring**, and **Failover**. The Chef Automate HA architecture is an approach of defining the components, modules, or implementation of services of a system that ensures optimal operational performance, even at times of high loads. It aids in addressing three major issues, server failure, zone failure, and cloud failure. It also allows you to automate and test everything in your network.

### What are Chef Automate Clusters?

The Chef Automate Clusters is a professional services solution offering installation, high availability, system uptime/ scale-out performance, maintenance, and disaster recovery capabilities.

It includes the Chef Infra Server API to simplify the Chef Infrastructure and is built for customers with more than 10,000 chef-client nodes. You can configure it in the private data center or preferred cloud.
