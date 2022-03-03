
**The following checks will be managed by this app:**

Regular scans

-   Org membership
    -   check everyone is an employee
-   Repository
    -   Check repository has legal approval
    -   Check repository is secured
    -   Check repository has addressed advisories
    -   Check repository has business approver
    -   Check repository has maintainer
    -   Check open source repo is public
    -   Check innersource repo is set as internal
    -   Check innersource repo has internal badge set
    -   Check open source repo has a badge set
    -   Check unarchived innersource repo has more than one commit - or wipe
    -   Check unarchived open source repo has more than one commit - or wipe
    -   Check unarchived innersource repo has a commit in last year
    -   Check unarchived open source repo has a commit in last year
-   Repository membership
    -   Check innersource repo has no external users
    -   Check repo has at least one admin
    -   Check maintainer is a repo member
    -   Check open source repo members are all employees
    -   Convert non employees to external collaborators
    -   Confirm each external collaborator been approved in last 3 months
    -   Identify external collaborators expiring in next 2 weeks





**Events**

The following events are monitored:



-   GitHub
    -   Repository-Created
    -   Repository-Deleted
    -   Repository-Renamed
    -   Repository-Publicized
    -   Repository-Privatized
    -   Member-Added
        -   User
            -   New user added to repo
        -   Collaborator
            -   New external collaborator added to repo
                -   External collaborator added to private repo
    -   Member-Edited
    -   Member-Removed
        -   User removed from repo
        -   External collaborator removed from repo
    -   Organization-Member_added
        -   Confirm user is employee
    -   Organization-Member_removed
    -   secret_scanning_alert
    -   repository_vulnerability_alert
-   User
    -   Leaver
    -   Daily cron job on each user
        -   Confirm user is employee
    -   Daily cron job on each repo
        -   Confirm repo rules
        -   Confirm repo membership rules
-   O365
    -   SendNotification
-   Repo management
    -   Actions
        -   Initiate repo request
        -   Initiate move repo request
        -   Initiate repo move
    -   Reactions
        -   Move request approved
        -   Move request rejected
        -   repo move complete
        -   New user not an employee
        -   No business sponsor
        -   Legal approval added
        -   Repository has unapproved licences
        -   Repository has unapproved IP