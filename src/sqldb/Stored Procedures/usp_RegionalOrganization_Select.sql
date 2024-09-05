CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select]
AS
BEGIN
    SELECT
      [Id],
      [Name],
      [IsCleanUpMembersEnabled],
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
END