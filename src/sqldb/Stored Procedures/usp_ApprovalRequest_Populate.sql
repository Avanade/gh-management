CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Populate]
	@CommunityId [INT]
AS
BEGIN
	DECLARE @output TABLE (RequestId INT)

	INSERT INTO [dbo].[ApprovalRequest]
		(
		[ApproverUserPrincipalName],
		[ApprovalStatusId],
		[ApprovalDescription],
		[CreatedBy],
		[ModifiedBy]
		)
	OUTPUT inserted.Id INTO @output
	SELECT
		[CA].[ApproverUserPrincipalName],
		1,
		'For Approval - ' + [C].[Name],
		[C].[CreatedBy],
		[C].[CreatedBy]
	FROM
		[dbo].[Community] AS [C],
		[dbo].[CommunityApprover] AS [CA]
	WHERE 
		[C].[Id] = @CommunityId
		AND [CA].[IsDisabled] = 0
		AND [CA].[GuidanceCategory] = 'community'

	UPDATE [dbo].[Community] SET [ApprovalStatusId] = 2, [Modified] = GETDATE() WHERE [Id] = @CommunityId

	INSERT INTO [dbo].[CommunityApprovalRequest]
		(
		[CommunityId],
		[ApprovalRequestId]
		)
	SELECT @CommunityId, [RequestId]
	FROM @output

	EXEC [dbo].[usp_ApprovalRequest_Select_ByCommunityId] @CommunityId
END