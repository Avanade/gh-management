CREATE PROCEDURE [dbo].[PR_ContributionAreas_Insert]
(
	@Name VARCHAR(100),
	@CreatedBy VARCHAR(100)
)
AS
BEGIN
	DECLARE @Id AS INT
    INSERT INTO [dbo].[ContributionAreas] (
		Name,
		Created,
		CreatedBy
	) VALUES (
		@Name,
		GETDATE(),
		@CreatedBy
	)
	SET @Id = SCOPE_IDENTITY()
	SELECT @Id Id
END