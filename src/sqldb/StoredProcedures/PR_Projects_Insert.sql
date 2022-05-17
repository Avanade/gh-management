CREATE PROCEDURE [dbo].[PR_Projects_Insert]
(
	@Name varchar(50),
	@CoOwner varchar(100),
	@Description varchar(1000),
	@ConfirmAvaIP bit,
	@ConfirmEnabledSecurity bit,
	@Username varchar(100)
) AS

INSERT INTO Projects (
	[Name],
	CoOwner,
	[Description],
	ConfirmAvaIP,
	ConfirmEnabledSecurity,
	Created,
	CreatedBy,
	Modified,
	ModifiedBy)
VALUES (
	@Name,
	@CoOwner,
	@Description,
	@ConfirmAvaIP,
	@ConfirmEnabledSecurity,
	GETDATE(),
	@Username,
	GETDATE(),
	@Username

)