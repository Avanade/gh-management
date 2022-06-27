create PROCEDURE  [dbo].[PR_CommunitySponsors_Select_By_CommunityId]
 @CommunityId int
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON

    -- Insert statements for procedure here
    SELECT CS.[Id]
      ,CS.[CommunityId]
      ,CS.[UserPrincipalName]
	  ,U.[Name]
	  ,U.[GivenName]
	  ,U.[SurName]
      ,CS.[Created]
      ,CS.[CreatedBy]
      ,CS.[Modified]
      ,CS.[ModifiedBy]
  FROM [dbo].[CommunitySponsors] CS
  INNER JOIN Users U ON CS.UserPrincipalName = U.UserPrincipalName
  where [CommunityId] = @CommunityId
END
