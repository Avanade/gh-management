CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select_ById]
  @Id [INT]
AS
BEGIN
  SELECT
    [Id],
    [Name],
    [IsRegionalOrganization],
    [IsIndexRepoEnabled],
    [IsCopilotRequestEnabled],
    [IsAccessRequestEnabled],
    [IsEnabled],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy]
  FROM 
    [dbo].[RegionalOrganization] 
  WHERE 
    [Id] = @Id
END
