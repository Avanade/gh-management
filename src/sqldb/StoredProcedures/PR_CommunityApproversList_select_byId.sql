
create PROCEDURE [dbo].[PR_CommunityApproversList_select_byId]
(
			@Id int 
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