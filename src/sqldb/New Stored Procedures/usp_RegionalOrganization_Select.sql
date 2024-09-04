CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select]
  @Id INT,
  @Name VARCHAR(100),
  @IsCleanUpMembersEnabled BIT,
  @IsIndexRepoEnabled BIT,
  @IsCopilotRequestEnabled BIT,
  @IsAccessRequestEnabled BIT,
  @IsEnabled BIT
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
      [dbo].[RegionalOrganization] 
END