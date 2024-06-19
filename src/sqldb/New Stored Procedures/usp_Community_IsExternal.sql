CREATE PROCEDURE [dbo].[usp_Community_IsExternal]
  @IsExternal [INT],
  @UserPrincipalName [VARCHAR](100)
AS
BEGIN
  SELECT
    [C].[Id],
    [C].[Name]
  FROM
    [dbo].[Community] AS [C]
    INNER JOIN [dbo].[ApprovalStatus] AS [T] ON [T].[Id] = [C].[ApprovalStatusId]
  WHERE
    [C].[IsExternal] = @IsExternal
    AND (	
      [C].[CreatedBy] = @UserPrincipalName OR [T].[Id] = 5
    )
END