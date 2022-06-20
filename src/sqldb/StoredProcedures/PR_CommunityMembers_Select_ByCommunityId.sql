Create PROCEDURE [dbo].[PR_CommunityMembers_Select_ByCommunityId]
(
  @CommunityId INT
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
SELECT CM.[Id],
       CM.[CommunityId],
       CM.[UserPrincipalName],
       U.[Name]
  FROM 
       [dbo].[CommunityMembers] CM
  INNER JOIN Users U ON U.UserPrincipalName = CM.[UserPrincipalName]
  WHERE
       CM.[CommunityId] = @CommunityId

END
