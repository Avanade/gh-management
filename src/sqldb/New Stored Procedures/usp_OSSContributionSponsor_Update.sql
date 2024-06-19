CREATE PROCEDURE [dbo].[usp_OSSContributionSponsor_Update]
  @Id [INT],
  @Name [VARCHAR](50),
  @IsArchived [BIT]
AS
BEGIN
  UPDATE [dbo].[OSSContributionSponsor]
  SET
    [Name] = @Name,
    [IsArchived] = @IsArchived
  WHERE
    [Id] = @Id
END
