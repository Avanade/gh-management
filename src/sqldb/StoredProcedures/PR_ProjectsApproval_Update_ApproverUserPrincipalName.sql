CREATE PROCEDURE PR_ProjectsApproval_Update_ApproverUserPrincipalName
	@Id uniqueidentifier,
	@ApproverEmail varchar(100),
 	@Username VARCHAR(100)
AS
BEGIN
	UPDATE 
		ProjectApprovals
	SET
		Modified = GETDATE(),
		ModifiedBy = @Username
	WHERE 
		ApprovalSystemguid= @Id;

	UPDATE 
		ApprovalRequestApprovers
	SET
		ApproverEmail = @ApproverEmail
	WHERE
		ApprovalRequestId = (
			SELECT 
				Id 
			FROM 
				ProjectApprovals 
			WHERE 
				ApprovalSystemguid = @Id
		) AND ApproverEmail = @Username;

	SELECT
		ARA.ApproverEmail AS ApproverUserPrincipalName,
		PA.Id,
		PA.ProjectId,
		P.Name AS ProjectName,
		P.Description AS ProjectDescription,
		U.Name AS RequesterName, 
		U.GivenName AS RequesterGivenName, 
		U.SurName AS RequesterSurName, 
		U.UserPrincipalName AS RequesterUserPrincipalName,
		PA.ApprovalTypeId, 
		T.Name AS ApprovalType,
		PA.ApprovalDescription,
		S.Name AS RequestStatus,
		PA.ApprovalDate, 
		PA.ApprovalRemarks,
		P.ConfirmAvaIP,
		P.ConfirmEnabledSecurity,
		P.ConfirmNotClientProject,
		P.newcontribution, 
		OCS.Name AS OSSsponsor, 
		P.Avanadeofferingsassets,
		P.Willbecommercialversion, 
		P.OSSContributionInformation
	FROM 
		ProjectApprovals AS PA
		INNER JOIN ApprovalTypes AS T ON PA.ApprovalTypeId = T.Id
		INNER JOIN Projects AS P ON PA.ProjectId = P.Id
		INNER JOIN Users AS U ON PA.CreatedBy = U.UserPrincipalName
		INNER JOIN ApprovalStatus AS S ON S.Id = PA.ApprovalStatusId
		INNER JOIN ApprovalRequestApprovers AS ARA ON PA.Id = ARA.ApprovalRequestId
		INNER JOIN OSSContributionSponsors AS OCS ON P.OSSContributionSponsorId = OCS.Id
	WHERE  
		PA.ApprovalSystemguid= @Id AND ARA.ApproverEmail = @ApproverEmail
END