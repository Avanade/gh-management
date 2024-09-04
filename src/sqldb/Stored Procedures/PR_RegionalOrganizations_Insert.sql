CREATE PROCEDURE [dbo].[PR_RegionalOrganizations_Insert]
(
    @Id INT,
    @Name VARCHAR(100),
    @IsCleanUpMembersEnabled BIT,
    @IsIndexRepoEnabled BIT,
    @IsCopilotRequestEnabled BIT,
    @IsAccessRequestEnabled BIT,
    @IsEnabled BIT        
)
AS
BEGIN
	SET NOCOUNT ON
    INSERT INTO [dbo].[RegionalOrganizations]
    (
        [Id],
        [Name],
        [IsCleanUpMembersEnabled],
        [IsIndexRepoEnabled],
        [IsCopilotRequestEnabled],
        [IsAccessRequestEnabled],
        [IsEnabled]
    )
    VALUES
    (
        @Id,
        @Name,
        @IsCleanUpMembersEnabled,
        @IsIndexRepoEnabled,
        @IsCopilotRequestEnabled,
        @IsAccessRequestEnabled,
        @IsEnabled
    )
END