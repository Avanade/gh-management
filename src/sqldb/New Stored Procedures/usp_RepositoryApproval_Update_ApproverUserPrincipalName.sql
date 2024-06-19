CREATE PROCEDURE [dbo].[usp_RepositoryApproval_Update_ApproverUserPrincipalName]
	@Id [UNIQUEIDENTIFIER],
	@ApproverEmail [VARCHAR](100),
 	@Username [VARCHAR](100)
AS
BEGIN
	UPDATE [dbo].[RepositoryApproval]
	SET
		[Modified] = GETDATE(),
		[ModifiedBy] = @Username
	WHERE [ApprovalSystemGUID]= @Id;

	UPDATE [dbo].[ApprovalRequestApprover]
	SET
		[ApproverUserPrincipalName] = @ApproverEmail
	WHERE [RepositoryApprovalId] = 
    (
      SELECT [Id] 
      FROM [dbo].[RepositoryApproval]
      WHERE [ApprovalSystemGUID] = @Id
		) AND 
    [ApproverUserPrincipalName] = @Username;

	SELECT
		[ARA].[ApproverUserPrincipalName],
		[RA].[Id],
		[RA].[RepositoryId],
		[R].[Name] AS [RepositoryName],
		[R].[Description] AS [RepositoryDescription],
		[U].[Name] AS [RequesterName], 
		[U].[GivenName] AS [RequesterGivenName], 
		[U].[SurName] AS [RequesterSurName], 
		[U].[UserPrincipalName] AS [RequesterUserPrincipalName],
		[RA].[RepositoryApprovalTypeId], 
		[T].[Name] AS [ApprovalType],
		[RA].[ApprovalDescription],
		[S].[Name] AS [RequestStatus],
		[RA].[ApprovalDate], 
		[RA].[ApprovalRemarks],
		[R].[ConfirmAvaIP],
		[R].[ConfirmEnabledSecurity],
		[R].[ConfirmNotClientProject],
		[R].[Newcontribution], 
		[OCS].[Name] AS [OSSsponsor], 
		[R].[Avanadeofferingsassets],
		[R].[Willbecommercialversion], 
		[R].[OSSContributionInformation]
	FROM [dbo].[RepositoryApproval] AS [RA]
		INNER JOIN [dbo].[RepositoryApprovalType] AS [T] ON [RA].[RepositoryApprovalTypeId] = [T].[Id]
		INNER JOIN [dbo].[Repository] AS [R] ON [RA].[RepositoryId] = [R].[Id]
		INNER JOIN [dbo].[User] AS [U] ON [RA].[CreatedBy] = [U].[UserPrincipalName]
		INNER JOIN [dbo].[ApprovalStatus] AS [S] ON [S].[Id] = [RA].[ApprovalStatusId]
		INNER JOIN [dbo].[ApprovalRequestApprover] AS [ARA] ON [RA].[Id] = [ARA].[RepositoryApprovalId]
		INNER JOIN [dbo].[OSSContributionSponsor] AS [OCS] ON [R].[OSSContributionSponsorId] = [OCS].[Id]
	WHERE  
		[RA].[ApprovalSystemGUID]= @Id AND [ARA].[ApproverUserPrincipalName] = @ApproverEmail
END
