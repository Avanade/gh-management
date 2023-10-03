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
	PA.ApprovalDescription,
	S.Name [RequestStatus],
	PA.[ApprovalDate], PA.[ApprovalRemarks],
	p.[ConfirmAvaIP],
	p.[ConfirmEnabledSecurity],
	p.[ConfirmNotClientProject],
	p.[newcontribution], 
	C.Name AS OSSsponsor, 
	p.[Avanadeofferingsassets],
	p.[Willbecommercialversion], 
	p.[OSSContributionInformation],
	PA.RespondedBy
    
FROM 
    ProjectApprovals PA
	INNER JOIN ApprovalTypes T ON PA.ApprovalTypeId = T.Id
	INNER JOIN Projects P ON PA.ProjectId = P.Id
	INNER JOIN Users U1 ON PA.CreatedBy = U1.UserPrincipalName
	INNER JOIN ApprovalStatus S ON S.Id = PA.ApprovalStatusId
	INNER JOIN OSSContributionSponsors C ON P.OSSContributionSponsorId = C.Id
WHERE  
    P.[Id] = @Id
END
