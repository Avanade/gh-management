
CREATE PROCEDURE [dbo].[PR_CommunityApproversList_select_byId]
(
    @Id INT
) AS
BEGIN
 SELECT [Id]
      ,[ApproverUserPrincipalName]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
      ,[Disabled]
  FROM [dbo].[CommunityApproversList]
 WHERE  [Id] = @Id
END