CREATE PROCEDURE [dbo].[usp_RepositoryApproval_Select_FailedRequest]
AS
BEGIN
  SELECT
    [RA].[Id],
    [RepositoryId],
    [R].[Name] AS [RepositoryName],
    [R].[CoOwner] AS [ProjectCoowner], 
    [R].[Description] AS [ProjectDescription],
    [U].[Name] AS [RequesterName], 
    [U].[GivenName] AS [RequesterGivenName], 
    [U].[SurName] AS [RequesterSurName], 
    [U].[UserPrincipalName] AS [RequesterUserPrincipalName],
    [C].[Name] AS [CoownerName], 
    [C].[GivenName] AS [CoownerGivenName], 
    [C].[SurName] AS [CoownerSurName], 
    [C].[UserPrincipalName] AS [CoownerUserPrincipalName],
    [RA].[RepositoryApprovalTypeId], 
    [T].[Name] AS [ApprovalType],
    [RA].[ApprovalDescription],
    [S].[Name] AS [RequestStatus],
    [RA].[ApprovalDate], 
    [RA].[ApprovalRemarks],
    [R].[ConfirmAvaIP],
    [R].[ConfirmEnabledSecurity],
    [R].[Newcontribution], 
    [CS].[Name] AS [OSSsponsor],
    [R].[Avanadeofferingsassets],
    [R].[Willbecommercialversion], 
    [R].[OSSContributionInformation]
  FROM [RepositoryApproval] AS [RA]
  	INNER JOIN [dbo].[RepositoryApprovalType] AS [T] ON [RA].[RepositoryApprovalTypeId] = [T].[Id]
	  INNER JOIN [dbo].[Repository] AS [R] ON [RA].[RepositoryId] = [R].[Id]
	  INNER JOIN [dbo].[User] AS [U] ON [R].[CreatedBy] = [U].[UserPrincipalName]
	  LEFT JOIN [dbo].[User] AS [C] ON [R].[CoOwner] = [C].[UserPrincipalName]
	  INNER JOIN [dbo].[ApprovalStatus] AS [S] ON [S].[Id] = [RA].[ApprovalStatusId]
	  INNER JOIN [dbo].[OSSContributionSponsor] AS [CS] ON [R].[OSSContributionSponsorId] = [CS].[Id]
  WHERE  
	  [ApprovalSystemGUID] IS NULL AND
	  DATEDIFF(MI, [RA].[Created], GETDATE()) >=5
END
