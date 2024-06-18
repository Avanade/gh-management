CREATE PROCEDURE [dbo].[usp_Community_Select_AllApproved]
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
    [T].[Name] AS [ApprovalStatus]
  FROM [dbo].[Community] AS [C]
    INNER JOIN [dbo].[ApprovalStatus] AS [T] ON [C].ApprovalStatusId = [T].[Id]
  WHERE 
    [C].[ApprovalStatusId] = 5
END