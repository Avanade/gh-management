CREATE PROCEDURE [dbo].[usp_Community_Update_Status]
  @CommunityId [INT]
AS
IF EXISTS (
	SELECT 
    [CA].[ApprovalStatusId] 
  FROM [dbo].[CommunityApprovalRequest] AS [CAR]
	LEFT JOIN [dbo].[ApprovalRequest] AS [CA] ON [CAR].[ApprovalRequestId] = [CA].[Id]
	WHERE [CAR].[CommunityId] = @CommunityId AND [CA].[ApprovalStatusId] <> 1
)
BEGIN
	UPDATE [dbo].[Community]
	SET [ApprovalStatusId] = (
    SELECT TOP 1 [CA].[ApprovalStatusId]
    FROM [dbo].[CommunityApprovalRequest] AS [CAR]
    LEFT JOIN [dbo].[ApprovalRequest] [CA] ON [CAR].[ApprovalRequestId] = [CA].[Id]
    WHERE [CAR].CommunityId = @CommunityId AND [CA].[ApprovalDate] IS NOT NULL ORDER BY [ApprovalDate]
  )
	WHERE Id = @CommunityId
END