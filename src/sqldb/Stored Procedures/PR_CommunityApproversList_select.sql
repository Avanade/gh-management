CREATE PROCEDURE [dbo].[PR_CommunityApproversList_select]
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
      ,[Disabled]
  FROM [dbo].[CommunityApproversList]
  WHERE [Category] = @Category
END