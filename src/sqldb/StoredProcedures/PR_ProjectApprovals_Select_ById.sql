CREATE PROCEDURE [dbo].[PR_ProjectApprovals_Select_ById]
(
	@Id INT
)
AS
BEGIN

SELECT
	PA.Id, PA.ProjectId, P.[Name] [ProjectName],
	P.[Description] [ProjectDescription],
	(SELECT STRING_AGG(ApproverEmail, ', ') FROM ApprovalRequestApprovers WHERE ApprovalRequestId = PA.Id GROUP BY ApprovalRequestId) Approvers,
	U1.Name [RequesterName], U1.GivenName [RequesterGivenName], U1.SurName [RequesterSurName], U1.UserPrincipalName [RequesterUserPrincipalName],
	PA.ApprovalTypeId, T.[Name] ApprovalType,
	PA.ApproverUserPrincipalName,
	PA.RespondedBy,
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
    P.[Id] = @Id
END
