CREATE PROCEDURE PR_CommunityApprovals_Populate
	@CommunityId INT
AS

INSERT INTO CommunityApprovals
	(
		CommunityId,
		ApproverUserPrincipalName,
		ApprovalStatusId,
		ApprovalDescription,
		CreatedBy,
		ModifiedBy
	)
	
SELECT @CommunityId, CAL.ApproverUserPrincipalName, 1, 'For Approval - ' + C.[Name], C.CreatedBy, C.CreatedBy
FROM Communities C, CommunityApproversList CAL
WHERE C.Id = @CommunityId

UPDATE Communities SET ApprovalStatusId = 2, Modified = GETDATE() WHERE Id = @CommunityId

exec PR_CommunityApprovals_Select_ById @CommunityId
