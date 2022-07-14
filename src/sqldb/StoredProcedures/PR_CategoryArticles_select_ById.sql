
CREATE PROCEDURE [dbo].[PR_CategoryArticles_select_ById]
@Id int
as 
begin
 

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
  FROM [dbo].[CategoryArticles] CA inner join Category C on CA.CategoryId = c.Id
  where CategoryId = @Id

 

end