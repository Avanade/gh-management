CREATE PROCEDURE [dbo].[PR_Projects_Insert]
(
	@Name varchar(50),
	@CoOwner varchar(100) = NULL,
	@Description varchar(1000),
	@IsArchived bit = 0,
	@ConfirmAvaIP bit = 0,
	@ConfirmEnabledSecurity bit = 0,
	@ConfirmNotClientProject bit = 0,
	@CreatedBy varchar(100) = NULL,
	@VisibilityId int = 1,
	@AssetCode varchar(50) = NULL,
	@TFSProjectReference varchar(150) = NULL,
	@AssetUrl varchar(150) = NULL,
	@MaturityRating varchar(20) = NULL,
	@ECATTReference varchar(150) = NULL,
	@Created DATETIME = NULL
) AS

IF @Created is null
	SET @Created = getdate()

DECLARE @ResultTable table(Id int);

INSERT INTO Projects (
	[Name],
	CoOwner,
	[Description],
	IsArchived,
	ConfirmAvaIP,
	ConfirmEnabledSecurity,
	ConfirmNotClientProject,
	Created,
	CreatedBy,
	Modified,
	ModifiedBy,
	VisibilityId,
	AssetCode,
	TFSProjectReference,
	AssetUrl,
	MaturityRating,
	ECATTReference)
OUTPUT INSERTED.Id INTO @ResultTable
VALUES (
	@Name,
	@CoOwner,
	@Description,
	@IsArchived,
	@ConfirmAvaIP,
	@ConfirmEnabledSecurity,
	@ConfirmNotClientProject,
	@Created,
	@CreatedBy,
	GETDATE(),
	@CreatedBy,
	@VisibilityId,
	@AssetCode,
	@TFSProjectReference,
	@AssetUrl,
	@MaturityRating,
	@ECATTReference
)

DECLARE @Id AS int

SELECT @Id = Id FROM @ResultTable

IF @CreatedBy IS NOT NULL
	EXEC [PR_UserAccess_Insert] @Id, @CreatedBy

IF @CoOwner IS NOT NULL
	EXEC [PR_UserAccess_Insert] @Id, @CoOwner

SELECT @Id [ItemId]