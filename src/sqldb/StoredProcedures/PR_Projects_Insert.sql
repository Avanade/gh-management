CREATE PROCEDURE [dbo].[PR_Projects_Insert]
(
	@Name varchar(50),
	@CoOwner varchar(100),
	@Description varchar(1000),
	@IsPrivate bit = 1,
	@IsArchived bit = 0,
	@ConfirmAvaIP bit,
	@ConfirmEnabledSecurity bit,
	@CreatedBy varchar(100),
	@Newcontribution varchar(50),
	@OSSsponsor varchar(50),
	@Avanadeofferingsassets varchar(50),
	@Willbecommercialversion varchar(50),
	@OSSContributionInformation varchar(50)
) AS

DECLARE @ResultTable table(Id int);

INSERT INTO Projects (
	[Name],
	CoOwner,
	[Description],
	IsPrivate,
	IsArchived,
	ConfirmAvaIP,
	ConfirmEnabledSecurity,
	Created,
	CreatedBy,
	Modified,
	ModifiedBy,	
		Newcontribution,
	OSSsponsor,
	Avanadeofferingsassets,
	Willbecommercialversion,
	OSSContributionInformation)
OUTPUT INSERTED.Id INTO @ResultTable
VALUES (
	@Name,
	@CoOwner,
	@Description,
	@IsPrivate,
	@IsArchived,
	@ConfirmAvaIP,
	@ConfirmEnabledSecurity,
	GETDATE(),
	@CreatedBy,
	GETDATE(),
	@CreatedBy,
	@Newcontribution,
	@OSSsponsor,
	@Avanadeofferingsassets,
	@Willbecommercialversion,
	@OSSContributionInformation
)

DECLARE @Id AS int

SELECT @Id = Id FROM @ResultTable

EXEC [PR_UserAccess_Insert] @Id, @CreatedBy

EXEC [PR_UserAccess_Insert] @Id, @CoOwner

SELECT @Id [ItemId]