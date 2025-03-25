CREATE PROCEDURE [dbo].[usp_OSSContributionSponsor_Insert]
	@Name [VARCHAR](50),
  @IsArchived [BIT] = 0
AS
BEGIN
  INSERT INTO [dbo].[OSSContributionSponsor]
  ( 
    [Name],
    [IsArchived]
  )
  OUTPUT
    [INSERTED].[Id]
  VALUES
  ( 
    @Name,
    @IsArchived
  )
END