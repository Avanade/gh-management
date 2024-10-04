CREATE PROCEDURE [dbo].[usp_CommunityApprover_Select]
  @GuidanceCategory [VARCHAR](100) 
AS
BEGIN
  SELECT 
    [Id],
    [ApproverUserPrincipalName],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy],
    [IsDisabled]
  FROM [dbo].[CommunityApprover]
  WHERE [GuidanceCategory] = @GuidanceCategory
END