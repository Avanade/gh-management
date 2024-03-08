CREATE TABLE [dbo].[OrganizationApprovalRequests]
(
	[OrganizationId] INT NOT NULL, 
    [RequestId] INT NOT NULL,
    PRIMARY KEY (OrganizationId, RequestId),
    CONSTRAINT [FK_OrganizationApprovalRequests_Organizations] FOREIGN KEY (OrganizationId) REFERENCES Organizations(Id)
)
