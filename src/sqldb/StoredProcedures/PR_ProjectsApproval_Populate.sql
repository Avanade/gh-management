CREATE PROCEDURE [dbo].[PR_ProjectsApproval_Populate]
    @ProjectId INT
AS

INSERT INTO ProjectApprovals
	(
		ProjectId,
		ApprovalTypeId,
		ApproverUserPrincipalName,
		ApprovalStatusId,
		ApprovalDescription,
		CreatedBy,
		ModifiedBy
	)
	
SELECT @ProjectId, T.Id, T.ApproverUserPrincipalName, 1, 'For Review - ' + T.[Name], P.CreatedBy, P.CreatedBy
FROM Projects P, ApprovalTypes T
WHERE T.ApproverUserPrincipalName IS NOT NULL AND T.IsActive = 1
AND P.Id = @ProjectId

UPDATE Projects SET ApprovalStatusId = 2, Modified = GETDATE() WHERE Id = @ProjectId

EXEC PR_ProjectApprovals_Select_ById @ProjectId