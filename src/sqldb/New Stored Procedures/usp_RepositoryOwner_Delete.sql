CREATE PROCEDURE [dbo].[usp_RepositoryOwner_Delete]
  @RepositoryId [INT],
  @UserPrincipalName [VARCHAR](100)
AS
BEGIN
	DELETE FROM [dbo].[RepositoryOwner]
  WHERE [RepositoryId] = @RepositoryId AND [UserPrincipalName] = @UserPrincipalName
END