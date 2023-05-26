CREATE PROCEDURE [dbo].[PR_CategoryArticles_select_ByArticlesID]
@Id INT
AS 
BEGIN
SELECT CA.[Id]
      ,CA.[Name]
	  ,CA.[URL]
	  ,CA.[Body]
	  ,CA.[CategoryId]
      ,CA.[Created]
      ,CA.[CreatedBy]
      ,CA.[Modified]
      ,CA.[ModifiedBy]
	  ,C.[Name] [CategoryName]
  FROM [dbo].[CategoryArticles] CA INNER JOIN Category C ON CA.CategoryId = c.Id
  WHERE CA.[Id] = @Id
END
GO
