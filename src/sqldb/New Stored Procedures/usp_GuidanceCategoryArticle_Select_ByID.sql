CREATE PROCEDURE [dbo].[usp_GuidanceCategoryArticle_Select_ByID]
  @Id [INT]
AS 
BEGIN
  SELECT 
    [CA].[Id],
    [CA].[Name],
	  [CA].[URL],
	  [CA].[Body],
	  [CA].[GuidanceCategoryId],
    [CA].[Created],
    [CA].[CreatedBy],
    [CA].[Modified],
    [CA].[ModifiedBy],
	  C.[Name] [CategoryName]
  FROM [dbo].[GuidanceCategoryArticle] [CA] INNER JOIN [dbo].[GuidanceCategory] AS [C] ON [CA].[GuidanceCategoryId] = [C].[Id]
  WHERE [CA].[Id] = @Id
END
GO
