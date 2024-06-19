CREATE PROCEDURE [dbo].[usp_RelatedCommunity_Select_ByParentCommunityId]
  @ParentCommunityId [INT]
AS
BEGIN
  SELECT   
    [ParentCommunityId],
    [RelatedCommunityId],
    [C].[IsExternal],
    [C].[Name]
  FROM [dbo].[RelatedCommunity] AS [RC]
  INNER JOIN [dbo].[Community] AS [C] ON [RC].[RelatedCommunityId] = [C].[Id]
  WHERE [ParentCommunityId] = @ParentCommunityId
END
