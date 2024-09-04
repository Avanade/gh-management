CREATE PROCEDURE [dbo].[PR_RegionalOrganizations_Update]
(
    @Id INT,
    @IsCleanUpMembersEnabled BIT,
    @IsIndexRepoEnabled BIT,
    @IsCopilotRequestEnabled BIT,
    @IsAccessRequestEnabled BIT,
    @IsEnabled BIT
)
AS
BEGIN
    SET NOCOUNT ON;

    UPDATE
        [dbo].[RegionalOrganization]
    SET
        [IsCleanUpMembersEnabled] = @IsCleanUpMembersEnabled,
        [IsIndexRepoEnabled] = @IsIndexRepoEnabled,
        [IsCopilotRequestEnabled] = @IsCopilotRequestEnabled,
        [IsAccessRequestEnabled] = @IsAccessRequestEnabled,
        [IsEnabled] = @IsEnabled
    WHERE
        [Id] = @Id
END