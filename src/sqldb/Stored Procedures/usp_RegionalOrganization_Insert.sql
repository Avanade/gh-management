CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Insert]
  @Id INT,
  @Name VARCHAR(100),
  @IsCleanUpMembersEnabled BIT = 1,
  @IsIndexRepoEnabled BIT = 1,
  @IsCopilotRequestEnabled BIT = 1,
  @IsAccessRequestEnabled BIT = 1,
  @IsEnabled BIT = 1,
  @CreatedBy VARCHAR(100)
AS
BEGIN
	SET NOCOUNT ON
  IF EXISTS (SELECT * FROM [dbo].[RegionalOrganization] WHERE [Id] = @Id)
  BEGIN
    EXEC [dbo].[usp_RegionalOrganization_Update] 
      @Id, @Name, @IsCleanUpMembersEnabled, 
      @IsIndexRepoEnabled, @IsCopilotRequestEnabled, 
      @IsAccessRequestEnabled, 1, @CreatedBy
  END
  ELSE
  BEGIN
    INSERT INTO [dbo].[RegionalOrganization]
    (
        [Id],
        [Name],
        [IsCleanUpMembersEnabled],
        [IsIndexRepoEnabled],
        [IsCopilotRequestEnabled],
        [IsAccessRequestEnabled],
        [IsEnabled],
        [CreatedBy]
    )
    VALUES
    (
        @Id,
        @Name,
        @IsCleanUpMembersEnabled,
        @IsIndexRepoEnabled,
        @IsCopilotRequestEnabled,
        @IsAccessRequestEnabled,
        @IsEnabled,
        @CreatedBy
    )
  END
END