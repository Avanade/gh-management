# Language and Platform Selection

- *Status*: Approved
- *Date*: 18 March 2022
- *RFC*: Approved without RFC due to business need.

## Context

We need a tool to ensure licence compliance, security compliance, and legal compliance, across various disparate tools and workflows.

Current tools on the market are hard to maintain, don't connect with existing processes (only maintaining IAM), and don't meet our needs.

As a stopgap, for requesting new repositories to add a layer of governance, a PowerApp with a Go backend was created. This backend is poorly tested, and doesn't ensure desired state configuration of repositories - e.g. maintaining private state, or transferring repositories between organizations.

## Decision

- Move the entire application into Go
- Model the workflow as a state machine
- Allow the ability to plug in  alternate legal, licence, and security workflows

## Consequences

- Easier
  - Easier to maintain the application
  - Easier to plug in and remove providers as our vendor landscape changes
- More difficult
  - Slightly harder for employee discovery
  - Additional overhead for SSO and integration with our enterprise AD