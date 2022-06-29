/****** Object:  StoredProcedure [dbo].[PR_ContributionAreas_Insert]    Script Date: 6/16/2022 3:47:37 PM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[PR_ContributionAreas_Insert]
(
	@Name varchar(100),
	@CreatedBy varchar(100)
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
GO

