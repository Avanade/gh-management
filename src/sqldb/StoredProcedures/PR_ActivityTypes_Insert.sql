/****** Object:  StoredProcedure [dbo].[PR_ActivityTypes_Insert]    Script Date: 6/16/2022 3:45:40 PM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[PR_ActivityTypes_Insert] 
(
	@Name VARCHAR(100)
)
AS
BEGIN
	DECLARE @Id AS INT
	SET @Id = (SELECT Id FROM [dbo].[ActivityTypes] WHERE Name=@Name)

	IF @Id IS NULL
	BEGIN
		INSERT INTO [dbo].[ActivityTypes] (Name) VALUES (@Name)
		SET @Id = SCOPE_IDENTITY()
	END
	SELECT @Id Id
END
GO

