# gh-management
> This tool allows an organization to manage users, and their association to a public GitHub, and a private GitHub enterprise (for InnerSource.)

[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)
![GitHub issues](https://img.shields.io/github/issues/Avanade/gh-management)
![GitHub](https://img.shields.io/github/license/Avanade/gh-management)
![GitHub Repo stars](https://img.shields.io/github/stars/Avanade/gh-management?style=social)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](https://avanade.github.io/code-of-conduct/)
[![Incubating InnerSource](https://img.shields.io/badge/Incubating-Ava--Maturity-%23FF5800?labelColor=yellow)](https://avanade.github.io/maturity-model/)

## Overview
<!-- TODO: Update overview -->

This project contains a PowerApp, which is the main entry point for the application. Backend functionality is provided through a Go application, which handles associations with GitHub, and events (including leavers.)

Microsoft has an excellent [GitHub management portal](https://github.com/microsoft/opensource-portal) to allow for self-service at scale - but this provides significantly more functionality than some organizations require.

It's recommended to call the leaving API as part of your JML process when users leave the organisation, but a Power Automate example is provided as an alternative (this is significantly less performant.)

This repository allows for basic self-service and automation of common workflows using PowerApps for:
- Automatically
  - Ensuring all users on GitHub are active employees
  - Checking each repository has an active business sponsor
  - Checking each repository has mandatory code scans
- Self-service
  - Request to join the GitHub organization
  - Request a new repository, and tracking of the ticket through approvals
  - Request permission for a code contribution of <50 lines

## Example Call Flow

  ```mermaid
  sequenceDiagram
    participant Client
    participant Business Logic
    participant Database Service
    participant Database Provider
    participant GitHubAPI Service
    participant GitHubAPI
    Client->>Business Logic: HTTP request
    Business Logic->>Business Logic: input validation
    Business Logic->>Database Service: data request
    Database Service->>Database Provider: execute stored procedure
    Database Provider->>Database Service: raw data result
    Database Service->>Database Service: parse raw data
    Database Service->>Business Logic: parsed data
    Business Logic->>GitHubAPI Service: GitHubAPI request
    GitHubAPI Service->>GitHubAPI Service: Prepare request
    GitHubAPI Service->>GitHubAPI: Execute API call
    GitHubAPI->>GitHubAPI Service: raw data result
    GitHubAPI Service->>GitHubAPI Service: parse raw data
    GitHubAPI Service->>Business Logic: parsed data
    Business Logic->>Business Logic:data calculations/transforms
    Business Logic->>Business Logic: package data per user-facing contract
    Business Logic->>Client: Request result
  ```

## Licensing
gh-management is available under the [MIT Licence](./LICENCE).

## Solutions Referenced
<!-- TODO: Update referenced solutions -->
- [Microsoft PowerApps](https://docs.microsoft.com/en-us/powerapps/WT.mc_id=AI-MVP-5004204)
- [Microsoft PowerAutomate](https://docs.microsoft.com/en-us/power-automate/?WT.mc_id=AI-MVP-5004204)

## Documentation
The `docs` folder contains [more detailed documentation](./docs/start-here.md), along with setup instructions.

## Contact
Feel free to [raise an issue on GitHub](https://github.com/Avanade/gh-management/issues), or see our [security disclosure](./SECURITY.md) policy.

## Contributing
Contributions are welcome. See information on [contributing](./CONTRIBUTING.md), as well as our [code of conduct](https://avanade.github.io/code-of-conduct/).

If you're happy to follow these guidelines, then check out the [getting started](./docs/start-here.md) guide.

## Who is Avanade?

[Avanade](https://www.avanade.com) is the leading provider of innovative digital and cloud services, business solutions and design-led experiences on the Microsoft ecosystem, and the power behind the Accenture Microsoft Business Group.