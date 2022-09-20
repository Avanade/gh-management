create PROCEDURE [dbo].[PR_CommunityApproversList_select]
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
END