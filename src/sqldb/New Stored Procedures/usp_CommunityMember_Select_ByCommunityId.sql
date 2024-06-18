CREATE PROCEDURE [dbo].[usp_CommunityMember_Select_ByCommunityId]
  @CommunityId [INT]
AS
BEGIN
  SET NOCOUNT ON;

  SELECT
    [CM].[CommunityId],
    [CM].[UserPrincipalName],
    [U].[Name]
  FROM [dbo].[CommunityMember] AS [CM]
  INNER JOIN [dbo].[User] AS [U] ON [U].[UserPrincipalName] = [CM].[UserPrincipalName]
  WHERE
    [CM].[CommunityId] = @CommunityId;
END