CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Update_ApproverUserPrincipalName]
  @ApprovalSystemGUID [UNIQUEIDENTIFIER],
  @ApproverUserPrincipalName [VARCHAR](100),
  @UserPrincipalName [VARCHAR](100)
AS
BEGIN
  SET NOCOUNT ON

  UPDATE [dbo].[ApprovalRequest]
	SET  
    [ApproverUserPrincipalName] = @ApproverUserPrincipalName,
    [Modified] = GETDATE(),
    [ModifiedBy] = @UserPrincipalName
	WHERE [ApprovalSystemGUID] = @ApprovalSystemGUID;

  SELECT
    [AR].[ApproverUserPrincipalName],
    [AR].[Id],
    [C].[Id] AS [CommunityId],
    [C].[Name] AS [ProjectName],
    [C].[Description] AS [ProjectDescription],
    [U].[Name] AS [RequesterName],
    [U].[GivenName] AS [RequesterGivenName],
    [U].[SurName] AS [RequesterSurName],
    [U].[UserPrincipalName] AS [RequesterUserPrincipalName],
    [AR].[ApprovalStatusId],
    [RAT].[Name] AS [ApprovalType],
    [C].[Url],
    [C].[Notes],
    [AR].[ApprovalDescription],
    [AS].[Name] AS [RequestStatus],
    [AR].[ApprovalDate],
    [AR].[ApprovalRemarks]
  FROM
    [dbo].[ApprovalRequest] AS [AR]
    INNER JOIN [dbo].[RepositoryApprovalType] AS [RAT] ON [AR].[ApprovalStatusId] = [RAT].[Id]
    INNER JOIN [dbo].[CommunityApprovalRequest] AS [CAR] ON [AR].[Id] = [CAR].[ApprovalRequestId]
    INNER JOIN [dbo].[Community] AS [C] ON [CAR].[CommunityId] = [C].[Id]
    INNER JOIN [dbo].[User] AS [U] ON [AR].[CreatedBy] = [U].[UserPrincipalName]
    INNER JOIN [dbo].[ApprovalStatus] AS [AS] ON [AS].[Id] = [AR].[ApprovalStatusId]
  WHERE 
		[AR].[ApprovalSystemGUID] = @ApprovalSystemGUID;
END