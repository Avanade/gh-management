﻿CREATE PROCEDURE [dbo].[PR_Projects_Insert]
(
	@Name varchar(50),
	@CoOwner varchar(100),
	@Description varchar(1000),
	@ConfirmAvaIP bit,
	@ConfirmEnabledSecurity bit,
	@CreatedBy varchar(100)
) AS

DECLARE @ResultTable table(Id int);

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
OUTPUT INSERTED.Id INTO @ResultTable
VALUES (
	@Name,
	@CoOwner,
	@Description,
	@ConfirmAvaIP,
	@ConfirmEnabledSecurity,
	GETDATE(),
	@CreatedBy,
	GETDATE(),
	@CreatedBy
)

DECLARE @Id AS int

SELECT @Id = Id FROM @ResultTable

EXEC [PR_UserAccess_Insert] @Id, @CreatedBy

EXEC [PR_UserAccess_Insert] @Id, @CoOwner