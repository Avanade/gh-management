CREATE PROCEDURE [dbo].[usp_Community_Select_RelatedCommunity]
	@CommunityId [INT]
AS
BEGIN
  SELECT DISTINCT
    [CR].[Id], 
    [CR].[Name], 
    [CR].[Url], 
    [CR].[IsExternal]
  FROM [dbo].[Community] AS [C]
  INNER JOIN [dbo].[CommunityTag] AS [CT] ON [C].[Id] = [CT].[CommunityId]
  LEFT JOIN [dbo].[CommunityTag] AS [CTR] ON [CT].[Tag] = [CTR].[Tag]
  LEFT JOIN [dbo].[Community] AS [CR] ON [CTR].[CommunityId] = [CR].[Id] AND [CR].[Id] <> @CommunityId
  WHERE [C].[Id] = @CommunityId
  AND [CR].[Id] IS NOT NULL
END