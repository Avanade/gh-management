CREATE PROCEDURE [dbo].[PR_ProjectApprovals_Select_Failed]

AS

SELECT
	PA.Id, ProjectId, P.[Name] [ProjectName],
	P.CoOwner [ProjectCoowner], P.[Description] [ProjectDescription],
	U1.Name [RequesterName], U1.GivenName [RequesterGivenName], U1.SurName [RequesterSurName], U1.UserPrincipalName [RequesterUserPrincipalName],
	U2.Name [CoownerName], U2.GivenName [CoownerGivenName], U2.SurName [CoownerSurName], U2.UserPrincipalName [CoownerUserPrincipalName],
	PA.ApprovalTypeId, T.[Name] ApprovalType,
	PA.ApprovalDescription,
	S.Name [RequestStatus],
	PA.[ApprovalDate], PA.[ApprovalRemarks],
	p.[ConfirmAvaIP],
	p.[ConfirmEnabledSecurity],
	p.[newcontribution], 
	C.Name AS OSSsponsor,
	p.[Avanadeofferingsassets],
	p.[Willbecommercialversion], 
	p.[OSSContributionInformation]
    
FROM 
	ProjectApprovals PA
	INNER JOIN ApprovalTypes T ON PA.ApprovalTypeId = T.Id
	INNER JOIN Projects P ON PA.ProjectId = P.Id
	INNER JOIN Users U1 ON P.CreatedBy = U1.UserPrincipalName
	LEFT JOIN Users U2 ON P.CoOwner = U2.UserPrincipalName
	INNER JOIN ApprovalStatus S ON S.Id = PA.ApprovalStatusId
	INNER JOIN OSSContributionSponsors C ON P.OSSContributionSponsorId = C.Id
WHERE  
	ApprovalSystemGUID IS NULL
	AND DATEDIFF(MI, PA.Created, GETDATE()) >=5