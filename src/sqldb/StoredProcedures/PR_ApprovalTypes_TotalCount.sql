/****** Object:  StoredProcedure [dbo].[PR_ApprovalTypes_Select]    Script Date: 6/24/2022 1:49:08 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[PR_ApprovalTypes_TotalCount]
AS
BEGIN
	SELECT COUNT(Id) AS 'Total' FROM [dbo].[ApprovalTypes]
END