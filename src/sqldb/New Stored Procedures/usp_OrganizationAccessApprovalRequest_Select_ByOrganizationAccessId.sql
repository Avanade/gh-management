CREATE PROCEDURE [dbo].[usp_OrganizationAccessApprovalRequest_Select_ByOrganizationAccessId]
  @Id [INT]
AS
BEGIN
  SELECT 
    [AR].[Id],
    [AR].[ApproverUserPrincipalName],
    [A].[Name] AS [ApprovalStatus],
    [AR].[ApprovalRemarks],
    [AR].[ApprovalDate],
    [AR].[ApprovalDescription]
  FROM [dbo].[OrganizationAccessApprovalRequest] AS [OAAR]
  LEFT JOIN [dbo].[ApprovalRequest] AS [AR] ON [OAAR].[ApprovalRequestId] = [AR].[Id]
  LEFT JOIN [dbo].[ApprovalStatus] AS [A] ON [A].[Id] = [AR].[ApprovalStatusId]
  WHERE [OAAR].[OrganizationAccessId] = @Id
END
