CREATE PROCEDURE [dbo].[PR_Repositories_TotalCount_BySearchTerm] (
	@search VARCHAR(50) = '',
	@Filter VARCHAR(MAX) = '',
  	@FilterType TINYINT = 0
)
AS
BEGIN
    SET NOCOUNT ON
	SELECT COUNT(*) AS Total FROM (
		SELECT COUNT(Id) AS Total FROM Projects AS p
		LEFT JOIN [dbo].[RepoTopics] AS rt ON rt.ProjectId = p.Id
		INNER JOIN STRING_SPLIT(@search, ' ') AS ss ON (
			p.Name LIKE '%'+ss.[value]+'%' OR rt.Topic LIKE '%'+ss.[value]+'%'
		)
		WHERE 
          @FilterType = 0
        OR
          (rt.Topic IN (SELECT value FROM STRING_SPLIT(@Filter, ',')) AND @FilterType = 1)
        OR
          (rt.Topic IS NULL AND @FilterType = 2)
		GROUP BY Id
	) AS Total
END