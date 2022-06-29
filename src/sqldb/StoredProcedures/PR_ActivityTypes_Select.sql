/****** Object:  StoredProcedure [dbo].[PR_ActivityTypes_Select]    Script Date: 6/16/2022 3:46:09 PM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[PR_ActivityTypes_Select]
AS
BEGIN
    SET NOCOUNT ON

    SELECT * FROM [dbo].[ActivityTypes]
END
GO