create PROCEDURE [dbo].[PR_RepoTopics_Select_PopularTopics]
(
	@offSet INT = 0,
	@rowCount INT = 0
)
AS
BEGIN
    IF @rowCount = 0
        SELECT
            Topic,
            COUNT(Topic) AS Total
        FROM 
            RepoTopics
        GROUP BY Topic ORDER BY Total DESC
    ELSE
        SELECT
            Topic,
            COUNT(Topic) AS Total
        FROM 
            RepoTopics
        GROUP BY Topic ORDER BY Total DESC
        OFFSET @offSet ROWS
        FETCH NEXT @rowCount ROWS ONLY
END