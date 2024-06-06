CREATE PROCEDURE [dbo].[usp_ContributionArea_Insert]
  @Name [VARCHAR](100),
  @CreatedBy [VARCHAR](100)
AS
BEGIN
  DECLARE @Id AS [INT]

  INSERT INTO [dbo].[ContributionArea]
  (
    [Name],
    [Created],
    [CreatedBy]
  )
  VALUES
  (
      @Name,
      GETDATE(),
      @CreatedBy
	)

  SET @Id = SCOPE_IDENTITY()
  SELECT @Id AS [Id]
END