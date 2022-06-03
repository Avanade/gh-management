CREATE PROCEDURE [dbo].[PR_ProjectsApproval_Populate]
    @ProjectId int
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
WHERE T.ApproverUserPrincipalName IS NOT NULL
AND P.Id = @ProjectId

exec PR_ProjectApprovals_Select_ById @ProjectId