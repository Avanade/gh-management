CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Insert]
  @Id [INT],
  @Name [VARCHAR](100)
AS
BEGIN
	SET NOCOUNT ON;
	IF NOT EXISTS (SELECT [Name] FROM [dbo].[RegionalOrganization] WHERE [Id] = @Id)
    BEGIN
      INSERT INTO [dbo].[RegionalOrganization]
      (
        [Id],
        [Name]
      )
      VALUES
      (
        @Id,
        @Name
      )
    END
  ELSE
    BEGIN
      UPDATE [dbo].[RegionalOrganization] 
      SET [Name] = @Name
      WHERE [Id] = @Id
    END
END