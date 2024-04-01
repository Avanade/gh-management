CREATE TABLE [dbo].[OrganizationAccessApprovalRequests]
(
	[OrganizationAccessId] INT NOT NULL, 
    [RequestId] INT NOT NULL,
    PRIMARY KEY (OrganizationAccessId, RequestId),
    CONSTRAINT [FK_OrganizationAccessApprovalRequests_OrganizationAccess] FOREIGN KEY (OrganizationAccessId) REFERENCES OrganizationAccess(Id),
    CONSTRAINT [FK_OrganizationAccessApprovalRequests_CommunityApprovals] FOREIGN KEY (RequestId) REFERENCES CommunityApprovals(Id)
)
