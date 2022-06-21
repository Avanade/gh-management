CREATE PROCEDURE  [dbo].[PR_CommunitySponsors_Select]
 
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON

    -- Insert statements for procedure here
    SELECT [Id]
      ,[CommunityId]
      ,[UserPrincipalName]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[CommunitySponsors]
END
