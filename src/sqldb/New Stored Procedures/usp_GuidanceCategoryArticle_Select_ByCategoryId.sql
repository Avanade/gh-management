CREATE PROCEDURE [dbo].[usp_GuidanceCategoryArticle_Select_ByCategoryId]
  @Id [INT]
AS 
BEGIN
  SELECT 
    CA.[Id],
    CA.[Name],
	  CA.[URL],
	  CA.[Body],
	  CA.[GuidanceCategoryId],
    CA.[Created],
    CA.[CreatedBy],
    CA.[Modified],
    CA.[ModifiedBy],
    C.[Name] AS [CategoryName]
  FROM [dbo].[GuidanceCategoryArticle] AS [CA] INNER JOIN [dbo].[GuidanceCategory] AS [C] ON [CA].[GuidanceCategoryId] = [C].[Id]
  WHERE [GuidanceCategoryId] = @Id
END