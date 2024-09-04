CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select_ById]
  @Id [INT]
AS
BEGIN
  SELECT
    [Id],
    [Name],
    [IsCleanUpMembersEnabled],
    [IsIndexRepoEnabled],
    [IsCopilotRequestEnabled],
    [IsAccessRequestEnabled],
    [IsEnabled]
  FROM 
    [dbo].[RegionalOrganizations] 
  WHERE 
    [Id] = @Id
END
