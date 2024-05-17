CREATE PROCEDURE [dbo].[PR_Communities_Update_Status]
  @CommunityId INT

AS

IF EXISTS (
	SELECT CA.ApprovalStatusId 
	FROM CommunityApprovalRequests CAR
	LEFT JOIN CommunityApprovals CA ON CAR.RequestId = CA.Id
	WHERE CAR.CommunityId = @CommunityId AND CA.ApprovalStatusId <> 1
) -- IF THERE IS A RESPONSE
BEGIN
	UPDATE Communities
	SET ApprovalStatusId = (
		SELECT TOP 1 CA.ApprovalStatusId 
		FROM CommunityApprovalRequests CAR
		LEFT JOIN CommunityApprovals CA ON CAR.RequestId = CA.Id
		WHERE CAR.CommunityId = @CommunityId AND CA.ApprovalDate IS NOT NULL ORDER BY ApprovalDate
		)
	WHERE Id = @CommunityId
END