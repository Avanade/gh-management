CREATE PROCEDURE PR_CommunityApprovals_Populate
	@CommunityId INT
AS
DECLARE @output TABLE (RequestId INT)
INSERT INTO CommunityApprovals
	(
		ApproverUserPrincipalName,
		ApprovalStatusId,
		ApprovalDescription,
		CreatedBy,
		ModifiedBy
	)
OUTPUT inserted.Id INTO @output
SELECT CAL.ApproverUserPrincipalName, 1, 'For Approval - ' + C.[Name], C.CreatedBy, C.CreatedBy
FROM Communities C, CommunityApproversList CAL
WHERE C.Id = @CommunityId
AND CAL.Disabled =0 AND CAL.Category = 'community'

UPDATE Communities SET ApprovalStatusId = 2, Modified = GETDATE() WHERE Id = @CommunityId

INSERT INTO CommunityApprovalRequests
	(
		CommunityId,
		RequestId
	)
SELECT @CommunityId, RequestId FROM @output
EXEC PR_CommunityApprovals_Select_ById @CommunityId