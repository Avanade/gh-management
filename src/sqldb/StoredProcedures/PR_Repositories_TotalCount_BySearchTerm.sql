SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[PR_Repositories_TotalCount_BySearchTerm] (
	@search varchar(50) = ''
)
AS
BEGIN
    SET NOCOUNT ON
	SELECT COUNT(Id) AS Total FROM Projects WHERE Name LIKE '%'+@search+'%'
END