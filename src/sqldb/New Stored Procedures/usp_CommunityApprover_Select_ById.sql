CREATE PROCEDURE [dbo].[usp_CommunityApprover_Select_ById]
  @Id [INT]
AS
BEGIN
  SELECT
    [Id],
    [ApproverUserPrincipalName],
    [GuidanceCategory],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy],
    [IsDisabled]
  FROM [dbo].[CommunityApprover]
  WHERE  [Id] = @Id
END