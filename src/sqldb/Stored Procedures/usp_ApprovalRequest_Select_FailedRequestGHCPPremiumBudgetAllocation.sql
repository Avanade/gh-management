CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Select_FailedRequestGHCPPremiumBudgetAllocation]
AS
BEGIN
  SELECT
    [GPBA].[Id] AS [Id],
    [GPBA].[OrganizationGitHubID] AS [OrganizationGitHubId],
    [GPBA].[OrganizationGitHubLogin] AS [OrganizationGitHubLogin],
    [GPBA].[UserGitHubID] AS [UserGitHubId],
    [GPBA].[UserGitHubLogin] AS [UserGitHubLogin],
    [GPBA].[Amount] AS [Amount],
    [GPBA].[Created] AS [Requested],
    [GPBA].[CreatedBy] AS [RequestedBy],
    STRING_AGG([AR].[ApproverUserPrincipalName], ',') AS [Approvers],
    STRING_AGG([AR].[Id], ',') AS [RequestIds]
  FROM [dbo].[ApprovalRequest] AS [AR]
    INNER JOIN [dbo].[GHCPPremiumBudgetAllocationApprovalRequest] AS [GPBAAR] ON [GPBAAR].[ApprovalRequestId] = [AR].[Id]
    INNER JOIN [dbo].[GHCPPremiumBudgetAllocation] AS [GPBA] ON [GPBA].[Id] = [GPBAAR].[GHCPPremiumBudgetAllocationId]
    INNER JOIN [dbo].[User] AS [U] ON [GPBA].[CreatedBy] = [U].[UserPrincipalName]
  WHERE
		[AR].[ApprovalSystemGUID] IS NULL
    AND DATEDIFF(MI, [AR].[Created], GETDATE()) >= 5
    GROUP BY [GPBA].[Id], [GPBA].[OrganizationGitHubID], [GPBA].[OrganizationGitHubLogin], [GPBA].[UserGitHubID], [GPBA].[UserGitHubLogin], [GPBA].[Amount], [GPBA].[Created], [GPBA].[CreatedBy]
END