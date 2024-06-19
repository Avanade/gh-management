CREATE PROCEDURE [dbo].[PR_CommunityApproversList_SelectAllActive]
(
  @Category VARCHAR(100) 
)
AS
BEGIN
 SELECT [Id]
      ,[ApproverUserPrincipalName]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[CommunityApproversList]
  WHERE [Disabled] = 0 AND [Category] = @Category
END