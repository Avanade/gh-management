CREATE PROCEDURE [dbo].[usp_CommunitySponsor_Select_ByCommunityId]
  @CommunityId [INT]
AS
BEGIN
  SET NOCOUNT ON

  SELECT
    [CS].[Id],
    [CS].[CommunityId],
    [CS].[UserPrincipalName],
    [U].[Name],
    [U].[GivenName],
    [U].[SurName],
    [CS].[Created],
    [CS].[CreatedBy],
    [CS].[Modified],
    [CS].[ModifiedBy]
  FROM [dbo].[CommunitySponsor] AS [CS]
  INNER JOIN [dbo].[User] AS [U] ON [CS].[UserPrincipalName] = [U].[UserPrincipalName]
  WHERE [CommunityId] = @CommunityId
END
