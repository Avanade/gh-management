CREATE PROCEDURE [dbo].[usp_RepositoryOwner_Insert]
  @RepositoryId [INT],
  @UserPrincipalName [VARCHAR](100)
AS
BEGIN
  IF NOT EXISTS (
    SELECT *
    FROM [dbo].[RepositoryOwner]
    WHERE [RepositoryId] = @RepositoryId AND [UserPrincipalName] = @UserPrincipalName
  )
  BEGIN
   INSERT INTO [dbo].[RepositoryOwner]
    (
      [RepositoryId],
      [UserPrincipalName]
    )
    VALUES
    (
      @RepositoryId,
      @UserPrincipalName
    )
  END
END