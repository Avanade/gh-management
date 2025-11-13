CREATE PROCEDURE [dbo].[usp_GHCPPremiumBudgetAllocationApprovalRequest_Select_ByGHCPPremiumBudgetAllocationId]
  @GHCPPremiumBudgetAllocationId [INT]
AS
BEGIN
  SELECT 
    [AR].[Id],
    [AR].[ApproverUserPrincipalName],
    [A].[Name] AS [ApprovalStatus],
    [AR].[ApprovalRemarks],
    [AR].[ApprovalDate],
    [AR].[ApprovalDescription]
  FROM [dbo].[GHCPPremiumBudgetAllocationApprovalRequest] AS [GPBAAR]
  LEFT JOIN [dbo].[ApprovalRequest] AS [AR] ON [GPBAAR].[ApprovalRequestId] = [AR].[Id]
  LEFT JOIN [dbo].[ApprovalStatus] AS [A] ON [A].[Id] = [AR].[ApprovalStatusId]
  WHERE [GPBAAR].[GHCPPremiumBudgetAllocationId] = @GHCPPremiumBudgetAllocationId
END
