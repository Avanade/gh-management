/****** Object:  StoredProcedure [dbo].[PR_CategoryArticles_select_ByArticlesID]    Script Date: 9/12/2022 6:09:45 PM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO


create PROCEDURE [dbo].[PR_CategoryArticles_select_ByArticlesID]
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
  where CA.[Id] = @Id

 

end
GO
