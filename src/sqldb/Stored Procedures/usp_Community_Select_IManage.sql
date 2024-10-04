CREATE PROCEDURE [dbo].[usp_Community_Select_IManage]
  @UserPrincipalName [VARCHAR](100)
AS
BEGIN
  SELECT 
    [C].[Id],
    [C].[Name],
    [C].[Url],
    [C].[Description],
    [C].[Notes],
    [C].[ApprovalStatusId],
    [C].[TradeAssocId],
    [C].[IsExternal],
    [C].[Created],
    [C].[CreatedBy],
    [C].[Modified],
    [C].[ModifiedBy],
    [AS].[Name] AS [ApprovalStatus]
  FROM [dbo].[Community] AS [C]
  INNER JOIN [dbo].[ApprovalStatus] AS [AS] ON [C].[ApprovalStatusId] = [AS].[Id]
  WHERE [C].[CreatedBy] = @UserPrincipalName
END