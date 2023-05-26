CREATE PROCEDURE [dbo].[PR_CommunityApproversList_SelectAllActive]
AS
BEGIN
 SELECT [Id]
      ,[ApproverUserPrincipalName]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[CommunityApproversList]
  WHERE [Disabled] = 0
END