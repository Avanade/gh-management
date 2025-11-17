CREATE TABLE [dbo].[GHCPPremiumBudgetAllocationApprovalRequest] (
    [GHCPPremiumBudgetAllocationId] [INT] NOT NULL,
    [ApprovalRequestId] [INT] NOT NULL,
    CONSTRAINT [PK_GHCPPremiumBudgetAllocationApprovalRequest] PRIMARY KEY ([GHCPPremiumBudgetAllocationId], [ApprovalRequestId]),
    CONSTRAINT [FK_GHCPPremiumBudgetAllocationApprovalRequest_PremiumBudgetAllocation] FOREIGN KEY ([GHCPPremiumBudgetAllocationId]) REFERENCES [dbo].[GHCPPremiumBudgetAllocation]([Id]),
    CONSTRAINT [FK_GHCPPremiumBudgetAllocationApprovalRequest_ApprovalRequest] FOREIGN KEY ([ApprovalRequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)
