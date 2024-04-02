CREATE PROCEDURE [dbo].[PR_CommunityApprovals_Select_FailedRequestOrganizationAccess]
AS
BEGIN
	SELECT
    OA.Id                                         [Id],
    RO.Name                                       [RegionalOrgName],
    UC.GitHubUser                                 [GitHubUsername],
    UC.UserPrincipalName                          [UserPrincipalName],
    STRING_AGG(CA.ApproverUserPrincipalName, ',') [Approvers],
    STRING_AGG(CA.Id, ',')                        [RequestIds]
	FROM CommunityApprovals CA
	INNER JOIN OrganizationAccessApprovalRequests AS OAAR ON OAAR.RequestId = CA.Id
	INNER JOIN OrganizationAccess OA ON OA.Id = OAAR.OrganizationAccessId
  INNER JOIN RegionalOrganizations RO ON RO.Id = OA.OrganizationId
	INNER JOIN Users UC ON OA.UserPrincipalName = UC.UserPrincipalName
	WHERE
		CA.ApprovalSystemGUID IS NULL
		AND DATEDIFF(MI, CA.Created, GETDATE()) >=5
  GROUP BY OA.Id, RO.Name, UC.GitHubUser, UC.UserPrincipalName
END