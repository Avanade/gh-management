CREATE TABLE [dbo].[ApprovalRequestApprovers] (
    [ApprovalRequestId] [INT] NOT NULL,
    [ApproverEmail] [VARCHAR](100) NOT NULL,
    CONSTRAINT [PK_ApprovalRequestApprover] PRIMARY KEY ([ApprovalRequestId], [ApproverEmail]),
    CONSTRAINT [FK_ApprovalRequestApprover_ProjectApprovals] FOREIGN KEY ([ApprovalRequestId]) REFERENCES [dbo].[ProjectApprovals]([Id]),
    CONSTRAINT [FK_ApprovalRequestApprover_Users] FOREIGN KEY ([ApproverEmail]) REFERENCES [dbo].[Users]([UserPrincipalName])
)
