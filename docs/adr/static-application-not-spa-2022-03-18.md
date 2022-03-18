# Static Application (not SPA)

- *Status*: Approved
- *Date*: 18 March 2022
- *RFC*: Approved without RFC due to business need.

## Context

In order for users to interact with the application, an interface is required. The Go Toolchain is able to serve React applications and other frontend frameworks, but these are built separately as part of build/deploy.

Providing realtime updates is not a key requirement of this application, and neither is fast interactivity. It's expected that most users will complete form data, and then come back to the application pending an approval or review.

The application also needs to support multiple APIs.

## Decision

- At this time, use the native Go net/http package
- Defer any use of an SPA until later

## Consequences

- Easier
    - Easier to maintain the application
    - Easier to develop at speed
    - Easier to handle authentication with AzureAD SSO and with GitHub
- More difficult
    - It will be harder to create the frontend as an SPA later as the underlying services will need to be refactored into APIs