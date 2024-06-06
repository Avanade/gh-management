CREATE PROCEDURE [dbo].[usp_ContributionArea_Update]
  @Id [INT],
  @Name [VARCHAR](100),
  @ModifiedBy [VARCHAR](100)
AS
BEGIN
  SET NOCOUNT ON
  
  UPDATE [dbo].[ContributionArea]
  SET 
    [Name] = @Name,
    [Modified] = GETDATE(),
    [ModifiedBy] =  @ModifiedBy
  WHERE [Id] = @Id
END