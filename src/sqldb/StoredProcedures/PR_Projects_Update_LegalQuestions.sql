CREATE PROCEDURE PR_Projects_Update_LegalQuestions
(
		@Id Int,
		@ModifiedBy varchar(100),
		@Newcontribution varchar(50),
		@OSSsponsor varchar(50),
		@Avanadeofferingsassets varchar(50),
		@Willbecommercialversion varchar(50),
		@OSSContributionInformation varchar(50)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
UPDATE 
		[dbo].[Projects]
   SET 
		[Newcontribution] = @Newcontribution,
		[OSSsponsor] = @OSSsponsor,
		[Avanadeofferingsassets] = @Avanadeofferingsassets,
		[Willbecommercialversion] = @Willbecommercialversion,
		[OSSContributionInformation] = @OSSContributionInformation,
		[Modified] = GETDATE(),
		[ModifiedBy] = @ModifiedBy
 WHERE  
		[Id] = @Id
END