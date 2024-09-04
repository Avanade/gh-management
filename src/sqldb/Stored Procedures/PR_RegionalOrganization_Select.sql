CREATE PROCEDURE [dbo].[PR_RegionalOrganizations_Select]
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
END