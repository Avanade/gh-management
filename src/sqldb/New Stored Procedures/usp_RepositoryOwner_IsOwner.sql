CREATE PROCEDURE [dbo].[usp_RepositoryOwner_IsOwner]
  @RepositoryId [INT],
	@UserPrincipalName [VARCHAR](100)
AS
BEGIN
  IF EXISTS (
    SELECT * FROM [dbo].[RepositoryOwner]
    WHERE [RepositoryId] = @RepositoryId AND [UserPrincipalName] = @UserPrincipalName
  )
  BEGIN
    SELECT 1 AS [IsOwner]
  END
  ELSE
  BEGIN
    SELECT 0 AS [IsOwner]
  END
END