CREATE PROCEDURE [dbo].[usp_RepositoryApproval_Select_ById]
	@Id [INT]
AS
BEGIN
  SELECT
    [RA].[Id],
    [RA].[RepositoryId],
    [R].[Name] AS [RepositoryName],
    [R].[Description] AS [RepositoryDescription],
    (
      SELECT STRING_AGG([ApproverUserPrincipalName], ', ') 
      FROM [dbo].[ApprovalRequestApprover] 
      WHERE [RepositoryApprovalId] = [RA].[Id] 
      GROUP BY [RepositoryApprovalId]
    ) AS [Approvers],
    [U].[Name] AS [RequesterName],
    [U].[GivenName] AS [RequesterGivenName], 
    [U].[SurName] AS [RequesterSurName], 
    [U].[UserPrincipalName] AS [RequesterUserPrincipalName],
    [RA].[RepositoryApprovalTypeId],
    [RAT].[Name] AS [ApprovalType],
    [RA].[ApprovalDescription],
    [AS].[Name] AS [RequestStatus],
    [RA].[ApprovalDate],
    [RA].[ApprovalRemarks],
    [R].[ConfirmAvaIP],
    [R].[ConfirmEnabledSecurity],
    [R].[ConfirmNotClientProject],
    [R].[Newcontribution], 
    [CS].[Name] AS [OSSsponsor], 
    [R].[Avanadeofferingsassets],
    [R].[Willbecommercialversion], 
    [R].[OSSContributionInformation],
    [RA].[RespondedBy],
    [RA].[Created]
  FROM [dbo].[RepositoryApproval] AS [RA]
    INNER JOIN [dbo].[RepositoryApprovalType] AS [RAT] ON [RA].[RepositoryApprovalTypeId] = [RAT].[Id]
    INNER JOIN [dbo].[Repository] AS [R] ON [RA].[RepositoryId] = [R].[Id]
    INNER JOIN [dbo].[User] AS [U] ON [RA].[CreatedBy] = [U].[UserPrincipalName]
    INNER JOIN [dbo].[ApprovalStatus] AS [AS] ON [AS].[Id] = [RA].[ApprovalStatusId]
    INNER JOIN [dbo].[OSSContributionSponsor] AS [CS] ON [R].[OSSContributionSponsorId] = [CS].[Id]
  WHERE [R].[Id] = @Id
END

