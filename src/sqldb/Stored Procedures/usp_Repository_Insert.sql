CREATE PROCEDURE [dbo].[usp_Repository_Insert]
	@Name [VARCHAR](50),
	@GithubId [INT],
	@CoOwner [VARCHAR](100) = NULL,
	@Description [VARCHAR](1000),
	@IsArchived [BIT] = 0,
	@ConfirmAvaIP [BIT] = 0,
	@ConfirmEnabledSecurity [BIT] = 0,
	@ConfirmNotClientProject [BIT] = 0,
	@CreatedBy [VARCHAR](100) = NULL,
	@Organization [VARCHAR](100),
	@VisibilityId [INT] = 1,
	@AssetCode [VARCHAR](50) = NULL,
	@TFSProjectReference [VARCHAR](150) = NULL,
	@AssetUrl [VARCHAR](150) = NULL,
	@MaturityRating [VARCHAR](20) = NULL,
	@ECATTReference [VARCHAR](150) = NULL
AS
BEGIN
  DECLARE @Id AS [INT]

  INSERT INTO [dbo].[Repository]
  (
    [GithubId],
    [Name],
    [Description],
    [IsArchived],
    [ConfirmAvaIP],
    [ConfirmEnabledSecurity],
    [ConfirmNotClientProject],
    [Created],
    [CreatedBy],
    [Organization],
    [Modified],
    [ModifiedBy],
    [VisibilityId],
    [AssetCode],
    [TFSProjectReference],
    [AssetUrl],
    [MaturityRating],
    [ECATTReference]
  )
  VALUES 
  (
    @GithubId,
    @Name,
    @Description,
    @IsArchived,
    @ConfirmAvaIP,
    @ConfirmEnabledSecurity,
    @ConfirmNotClientProject,
    GETDATE(),
    @CreatedBy,
    @Organization,
    GETDATE(),
    @CreatedBy,
    @VisibilityId,
    @AssetCode,
    @TFSProjectReference,
    @AssetUrl,
    @MaturityRating,
    @ECATTReference
  )

  SET @Id = SCOPE_IDENTITY()
	SELECT @Id [Id]
END