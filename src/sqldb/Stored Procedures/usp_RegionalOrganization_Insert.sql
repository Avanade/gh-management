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