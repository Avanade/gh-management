CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Update_StatusByOrganizationId]
  @OrganizationId [INT],
  @ApprovalStatusId [INT]
AS
UPDATE [dbo].[ApprovalRequest]
SET
    [ApprovalStatusId] = @ApprovalStatusId
WHERE [Id] IN (
  SELECT [ApprovalRequestId] FROM [dbo].[OrganizationApprovalRequest] WHERE [OrganizationId] = @OrganizationId
)
