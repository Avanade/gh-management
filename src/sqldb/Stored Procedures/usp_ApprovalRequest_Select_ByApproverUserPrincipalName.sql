CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Select_ByApproverUserPrincipalName]
	@ApproverUserPrincipalName [VARCHAR](100)
AS
BEGIN
	SELECT
        CONVERT(varchar(36), [PA].[ApprovalSystemGUID], 1) AS [ApprovalSystemGUID],
        [R].[Name] AS [ProjectName],
        [PA].[CreatedBy] AS [RequestedBy],
        [R].[Description] AS [Description],
        [R].[Newcontribution] AS [NewContribution],
        [OCS].[Name] AS [OSSContributionSponsor],
        [R].[Avanadeofferingsassets] AS [IsAvanadeOfferingAssets],
        [R].[Willbecommercialversion] AS [WillBeCommercialVersion],
        [R].[OSSContributionInformation] AS [OSSContributionInformation],
        [PA].[ApprovalRemarks] AS [Remarks],
        [A].[Name] AS [Status],
        [PA].[RespondedBy] AS [RespondedBy],
        [PA].[ApprovalDate] AS [ApprovalDate],
        [PA].[ApprovalSystemDateSent] AS [ApprovalSystemDateSent]
    FROM
        [dbo].[Repository] AS [R]
        INNER JOIN [dbo].[RepositoryApproval] AS [PA] ON [R].[Id] = [PA].[RepositoryId]
        INNER JOIN [dbo].[ApprovalStatus] AS [A] ON [PA].[ApprovalStatusId] = [A].[Id]
        INNER JOIN [dbo].[OSSContributionSponsor] AS [OCS] ON [R].[OSSContributionSponsorId] = [OCS].[Id]
        INNER JOIN [dbo].[ApprovalRequestApprover] AS [ARA] ON [PA].[Id] = [ARA].[RepositoryApprovalId]
    WHERE
        [ARA].[ApproverUserPrincipalName] = @ApproverUserPrincipalName
END