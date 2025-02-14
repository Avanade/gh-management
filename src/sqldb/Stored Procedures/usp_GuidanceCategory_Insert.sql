CREATE PROCEDURE [dbo].[usp_GuidanceCategory_Insert]
	@Name [VARCHAR](50),
	@CreatedBy [VARCHAR](50),
	@ModifiedBy [VARCHAR](50)
AS
BEGIN   
	INSERT INTO [dbo].[GuidanceCategory]
	(
		[Name],
		[Created],
		[CreatedBy],
		[Modified],
		[ModifiedBy]
	)
	OUTPUT 
		[INSERTED].Id
	VALUES
	(
		@Name,
		GETDATE(),
		@CreatedBy,
		GETDATE(),
		@ModifiedBy
	)
END