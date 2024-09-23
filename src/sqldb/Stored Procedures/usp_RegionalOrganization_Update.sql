CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Update]
  @Id INT,
  @Name VARCHAR(100),
  @IsCleanUpMembersEnabled BIT,
  @IsIndexRepoEnabled BIT,
  @IsCopilotRequestEnabled BIT,
  @IsAccessRequestEnabled BIT,
  @IsEnabled BIT,
  @ModifiedBy VARCHAR(100)
AS
BEGIN
	SET NOCOUNT ON
  UPDATE
    [dbo].[RegionalOrganization]
  SET
    [IsCleanUpMembersEnabled] = @IsCleanUpMembersEnabled,
    [IsIndexRepoEnabled] = @IsIndexRepoEnabled,
    [IsCopilotRequestEnabled] = @IsCopilotRequestEnabled,
    [IsAccessRequestEnabled] = @IsAccessRequestEnabled,
    [IsEnabled] = @IsEnabled,
    [ModifiedBy] = @ModifiedBy,
    [Modified] = GETDATE()
  WHERE
    [Id] = @Id
END