CREATE PROCEDURE [dbo].[PR_Applications_Insert]
	@Name varchar(100),
	@IsActive bit = true
AS
	INSERT INTO Applications ([Name], IsActive)
	VALUES (@Name, @IsActive)
