CREATE PROCEDURE [dbo].[PR_RegionalOrganizations_SelectById]
(
    @Id int
)
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