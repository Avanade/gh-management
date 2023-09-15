CREATE TABLE [dbo].[Approvers]
(   
    [ApprovalTypeId] INT NOT NULL,
    [ApproverEmail] VARCHAR(100) NOT NULL,
    CONSTRAINT PK_Approver PRIMARY KEY (ApprovalTypeId, ApproverEmail),
    CONSTRAINT FK_Approvers_ApprovalTypes FOREIGN KEY (ApprovalTypeId) REFERENCES ApprovalTypes(Id),
    CONSTRAINT FK_Approvers_Users FOREIGN KEY (ApproverEmail) REFERENCES Users(UserPrincipalName)
)