CREATE PROCEDURE [dbo].[usp_User_IsExisting]
	@UserPrincipalName [VARCHAR](100)
AS
BEGIN
  IF EXISTS (
    SELECT [UserPrincipalName]
    FROM [dbo].[User]
    WHERE [UserPrincipalName] = @UserPrincipalName
  )
    BEGIN
      SELECT '1' AS [Result]
      RETURN 1
    END
  ELSE
    BEGIN
      SELECT '0' AS [Result]
      RETURN 0
    END
END