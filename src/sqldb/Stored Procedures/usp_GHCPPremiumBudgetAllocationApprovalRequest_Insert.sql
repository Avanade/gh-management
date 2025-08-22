CREATE PROCEDURE [dbo].[usp_GHCPPremiumBudgetAllocationApprovalRequest_Insert]
	@GHCPPremiumBudgetAllocationId [INT],
	@ApprovalRequestId [INT]
AS
BEGIN
	INSERT INTO [dbo].[GHCPPremiumBudgetAllocationApprovalRequest]
	(
		[GHCPPremiumBudgetAllocationId],
        [ApprovalRequestId]
	)
	VALUES
    (
		@GHCPPremiumBudgetAllocationId,
        @ApprovalRequestId
	)
END