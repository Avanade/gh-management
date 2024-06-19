CREATE PROCEDURE [dbo].[usp_Repository_Update_LegalQuestions]
  @Id [INT],
  @ModifiedBy [VARCHAR](100),
  @Newcontribution [VARCHAR](50),
  @OSSsponsor [INT],
  @Avanadeofferingsassets [VARCHAR](50),
  @Willbecommercialversion [VARCHAR](50),
  @OSSContributionInformation [VARCHAR](1000)
AS
BEGIN
	SET NOCOUNT ON;

  UPDATE 
    [dbo].[Repository]
  SET 
    [Newcontribution] = @Newcontribution,
    [OSSContributionSponsorId] = @OSSsponsor,
    [Avanadeofferingsassets] = @Avanadeofferingsassets,
    [Willbecommercialversion] = @Willbecommercialversion,
    [OSSContributionInformation] = @OSSContributionInformation,
    [Modified] = GETDATE(),
    [ModifiedBy] = @ModifiedBy
  WHERE  
    [Id] = @Id
END