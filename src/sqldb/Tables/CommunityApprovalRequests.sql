CREATE TABLE [dbo].[CommunityApprovalRequests]
(
	[CommunityId] INT NOT NULL, 
    [RequestId] INT NOT NULL,
    PRIMARY KEY (CommunityId, RequestId),
    CONSTRAINT [FK_CommunityApprovalRequests_Communities] FOREIGN KEY (CommunityId) REFERENCES Communities(Id),
    CONSTRAINT [FK_CommunityApprovalRequests_CommunityApprovals] FOREIGN KEY (RequestId) REFERENCES CommunityApprovals(Id)
)
