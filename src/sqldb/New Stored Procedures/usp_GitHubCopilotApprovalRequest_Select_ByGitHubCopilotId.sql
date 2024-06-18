CREATE PROCEDURE [dbo].[usp_GitHubCopilotApprovalRequest_Select_ByGitHubCopilotId]
  @GitHubCopilotId [INT]
AS
BEGIN
  SELECT 
    [AR].[Id],
    [AR].[ApproverUserPrincipalName],
    [A].[Name] AS [ApprovalStatus],
    [AR].[ApprovalRemarks],
    [AR].[ApprovalDate],
    [AR].[ApprovalDescription]
  FROM [dbo].[GitHubCopilotApprovalRequest] AS [GCAR]
  LEFT JOIN [dbo].[ApprovalRequest] AS [AR] ON [GCAR].[ApprovalRequestId] = [AR].[Id]
  LEFT JOIN [dbo].[ApprovalStatus] AS [A] ON [A].[Id]  = [AR].[ApprovalStatusId]
  WHERE [GCAR].[GitHubCopilotId] = @GitHubCopilotId
END
