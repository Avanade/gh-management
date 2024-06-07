CREATE PROCEDURE [dbo].[usp_Organization_Update]
  @Id [INT],
  @ApprovalStatusId [INT]
AS
UPDATE [dbo].[ApprovalRequest]
SET
    [ApprovalStatusId] = @ApprovalStatusId
WHERE [Id] = @Id
