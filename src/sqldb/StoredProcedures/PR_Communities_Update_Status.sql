CREATE PROCEDURE [dbo].[PR_Communities_Update_Status]
  @CommunityId INT

AS

IF EXISTS (SELECT ApprovalStatusId FROM CommunityApprovals WHERE CommunityId = @CommunityId AND ApprovalStatusId <> 1) -- IF THERE IS A RESPONSE
BEGIN
	UPDATE Communities
	SET ApprovalStatusId = (SELECT TOP 1 ApprovalStatusId FROM CommunityApprovals WHERE CommunityId = @CommunityId AND ApprovalDate IS NOT NULL ORDER BY ApprovalDate)
	WHERE Id = @CommunityId
END
