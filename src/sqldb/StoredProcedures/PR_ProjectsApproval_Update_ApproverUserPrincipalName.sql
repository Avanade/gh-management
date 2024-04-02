CREATE PROCEDURE [dbo].[PR_ProjectsApproval_Update_ApproverUserPrincipalName]
	@Id uniqueidentifier,
	@ApproverEmail varchar(100),
 	@Username VARCHAR(100)
AS

UPDATE [ProjectApprovals]
SET  
		ApproverUserPrincipalName = @ApproverEmail,
 		Modified = GETDATE(),
		ModifiedBy = @Username
WHERE ApprovalSystemguid= @Id
SELECT

	PA.ApproverUserPrincipalName,
	PA.Id, PA.ProjectId, P.[Name] [ProjectName],
	P.[Description] [ProjectDescription],
	U1.Name [RequesterName], U1.GivenName [RequesterGivenName], U1.SurName [RequesterSurName], U1.UserPrincipalName [RequesterUserPrincipalName],
	PA.ApprovalTypeId, T.[Name] ApprovalType,

	PA.ApprovalDescription,
	S.Name [RequestStatus],
	PA.[ApprovalDate], PA.[ApprovalRemarks],
	p.[ConfirmAvaIP],
	p.[ConfirmEnabledSecurity],
	p.[ConfirmNotClientProject],
	p.[newcontribution], 
	p.[OSSsponsor], 
	p.[Avanadeofferingsassets],
	p.[Willbecommercialversion], 
	p.[OSSContributionInformation]
    
FROM 
    ProjectApprovals PA
	INNER JOIN ApprovalTypes T ON PA.ApprovalTypeId = T.Id
	INNER JOIN Projects P ON PA.ProjectId = P.Id
	INNER JOIN Users U1 ON PA.CreatedBy = U1.UserPrincipalName
	INNER JOIN ApprovalStatus S ON S.Id = PA.ApprovalStatusId
WHERE  
      pa.ApprovalSystemguid= @Id 