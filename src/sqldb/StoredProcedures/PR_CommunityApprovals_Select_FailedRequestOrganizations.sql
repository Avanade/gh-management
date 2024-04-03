CREATE PROCEDURE [dbo].[PR_CommunityApprovals_Select_FailedRequestOrganizations]
AS
BEGIN
	SELECT
    O.Id                                          [Id],
    RO.Id                                         [RegionId],
    RO.Name                                       [RegionName],
    O.ClientName                                  [ClientName],
    O.ProjectName                                 [ProjectName],
    O.WBS                                         [WBS],
    O.CreatedBy                                   [Username],
    STRING_AGG(CA.ApproverUserPrincipalName, ',') [Approvers],
    STRING_AGG(CA.Id, ',')                        [RequestIds]
	FROM CommunityApprovals CA
	INNER JOIN OrganizationApprovalRequests AS OAR ON OAR.RequestId = CA.Id
	INNER JOIN Organizations O ON O.Id = OAR.OrganizationId
  INNER JOIN RegionalOrganizations RO ON RO.Id = O.Region
	INNER JOIN Users UC ON O.CreatedBy = UC.UserPrincipalName
	WHERE
		CA.ApprovalSystemGUID IS NULL
		AND DATEDIFF(MI, CA.Created, GETDATE()) >=5
  GROUP BY O.Id, RO.Id, RO.Name, O.ClientName, O.ProjectName, O.WBS, O.CreatedBy
END