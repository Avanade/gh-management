CREATE TABLE [dbo].[ApprovalRequestApprovers]
(   
    [ApprovalRequestId] INT NOT NULL,
    [ApproverEmail] VARCHAR(100) NOT NULL
    CONSTRAINT PK_ApprovalRequestApprover PRIMARY KEY (ApprovalRequestId, ApproverEmail),
    CONSTRAINT FK_ApprovalRequestApprover_A FOREIGN KEY (ApprovalRequestId) REFERENCES ProjectApprovals(Id),
    CONSTRAINT FK_ApprovalRequestApprover_Users FOREIGN KEY (ApproverEmail) REFERENCES Users(UserPrincipalName)
)