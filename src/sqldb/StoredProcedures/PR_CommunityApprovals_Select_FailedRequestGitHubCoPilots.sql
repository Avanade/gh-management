CREATE PROCEDURE [dbo].[PR_CommunityApprovals_Select_FailedRequestGitHubCoPilots]
AS
BEGIN
	SELECT
    GC.Id                                           [Id],
    RO.Id                                           [RegionId],
    RO.Name                                         [RegionName],
    UC.GitHubId                                     [GitHubId],
    UC.GitHubUser                                   [GitHubUsername],
    UC.UserPrincipalName                            [Username],
    STRING_AGG(CA.ApproverUserPrincipalName, ',')   [Approvers],
    STRING_AGG(CA.Id, ',')                          [RequestIds]
	FROM CommunityApprovals CA
	INNER JOIN GitHubCopilotApprovalRequests GCAR ON GCAR.RequestId = CA.Id
	INNER JOIN GitHubCopilot GC ON GC.Id = GCAR.GitHubCopilotId
  INNER JOIN RegionalOrganizations RO ON RO.Id = GC.Region
	INNER JOIN Users UC ON GC.CreatedBy = UC.UserPrincipalName
	WHERE
		CA.ApprovalSystemGUID IS NULL
		AND DATEDIFF(MI, CA.Created, GETDATE()) >=5
  GROUP BY GC.Id, RO.Id, RO.Name, UC.GitHubId, UC.GitHubUser, UC.UserPrincipalName
END