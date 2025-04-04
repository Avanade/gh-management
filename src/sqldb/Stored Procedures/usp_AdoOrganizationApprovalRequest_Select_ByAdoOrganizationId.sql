CREATE PROCEDURE [dbo].[usp_AdoOrganizationApprovalRequest_Select_ByAdoOrganizationId]
  @AdoOrganizationId [INT]
AS
BEGIN
  SELECT 
    [AR].[Id],
    [AR].[ApproverUserPrincipalName],
    [A].[Name] AS [ApprovalStatus],
    [AR].[ApprovalRemarks],
    [AR].[ApprovalDate],
    [AR].[ApprovalDescription]
  FROM [dbo].[AdoOrganizationApprovalRequest] AS [AOAR]
  LEFT JOIN [dbo].[ApprovalRequest] AS [AR] ON [AOAR].[ApprovalRequestId] = [AR].[Id]
  LEFT JOIN [dbo].[ApprovalStatus] AS [A] ON [A].[Id] = [AR].[ApprovalStatusId]
  WHERE [AOAR].[AdoOrganizationId] = @AdoOrganizationId
END
