CREATE PROCEDURE [dbo].[usp_CommunityApprover_Select_Active]
  @GuidanceCategory [VARCHAR](100)
AS
BEGIN
  SELECT
    [Id],
    [ApproverUserPrincipalName],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy]
  FROM [dbo].[CommunityApprover]
  WHERE [IsDisabled] = 0 AND [GuidanceCategory] = @GuidanceCategory
END 