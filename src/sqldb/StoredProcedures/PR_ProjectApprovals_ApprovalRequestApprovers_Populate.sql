CREATE PROCEDURE [dbo].[PR_ProjectApprovals_ApprovalRequestApprovers_Populate]
(
    @ProjectId INT,
	@RequestedBy VARCHAR(100)
)
AS
BEGIN
    INSERT INTO ProjectApprovals
    (
		ProjectId,
		ApprovalTypeId,
		-- ApproverUserPrincipalName,
		ApprovalStatusId,
		ApprovalDescription,
		CreatedBy
	)
    SELECT 
        @ProjectId, 
        T.Id, 
        -- T.ApproverUserPrincipalName, 
        1, 
        'For Review - ' + T.[Name], 
        @RequestedBy
    FROM 
        ApprovalTypes AS T
    WHERE 
        T.IsActive = 1 AND 
        T.IsArchived = 0

    INSERT INTO ApprovalRequestApprovers
    (
        ApprovalRequestId,
        ApproverEmail
    )
    SELECT
        PA.Id,
        A.ApproverEmail
    FROM
        ProjectApprovals AS PA 
        INNER JOIN
        ApprovalTypes AS T ON PA.ApprovalTypeId = T.Id
        INNER JOIN 
        Approvers AS A ON T.Id = A.ApprovalTypeId
    WHERE
        PA.ProjectId = @ProjectId

    UPDATE Projects SET ApprovalStatusId = 2, Modified = GETDATE() WHERE Id = @ProjectId

    EXEC PR_ProjectApprovals_Select_ById @ProjectId
END