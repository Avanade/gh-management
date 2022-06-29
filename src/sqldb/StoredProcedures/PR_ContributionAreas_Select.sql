/****** Object:  StoredProcedure [dbo].[PR_ContributionAreas_Select]    Script Date: 6/16/2022 3:48:07 PM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[PR_ContributionAreas_Select]
AS
BEGIN
    SELECT * FROM [dbo].[ContributionAreas]
END
GO