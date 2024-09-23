CREATE PROCEDURE [dbo].[usp_Community_Select_ById]
  @Id [INT]
AS 
BEGIN
  SELECT 
    [Id],
    [Name],
    [Url],
    [Description],
    [Notes],
    [TradeAssocId],
    [IsExternal],
    [CommunityType],
    [ChannelId],
    [OnBoardingInstructions],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy],
    [ApprovalStatusId]
  FROM 
    [dbo].[Community]
  WHERE [Id] = @Id
END