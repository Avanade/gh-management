CREATE PROCEDURE [dbo].[usp_OrganizationApprovalRequest_Select_ByOrganizationId]
  @OrganizationId [INT]
AS
BEGIN
  SELECT 
    [AR].[Id],
    [AR].[ApproverUserPrincipalName],
    [A].[Name] AS [ApprovalStatus],
    [AR].[ApprovalRemarks],
    [AR].[ApprovalDate],
    [AR].[ApprovalDescription]
  FROM [dbo].[OrganizationApprovalRequest] AS [OAR]
  LEFT JOIN [dbo].[ApprovalRequest] AS [AR] ON [OAR].[ApprovalRequestId] = [AR].[Id]
  LEFT JOIN [dbo].[ApprovalStatus] AS [A] ON [A].[Id] = [AR].[ApprovalStatusId]
  WHERE [OAR].[OrganizationId] = @OrganizationId
END
