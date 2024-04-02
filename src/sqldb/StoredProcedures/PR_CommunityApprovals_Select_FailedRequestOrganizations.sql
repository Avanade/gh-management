CREATE PROCEDURE [dbo].[PR_CommunityApprovals_Select_FailedRequestOrganizations]
AS
BEGIN
	SELECT
		RO.Id                         [Region],
    O.ClientName                  [ClientName],
    O.ProjectName                 [ProjectName],
    O.WBS                         [WBS],
    O.CreatedBy                   [Username],
    O.Id                          [Id],
    CA.ApproverUserPrincipalName  [ApproverUserPrincipalName],
    RO.Name                       [RegionName],
    CA.Id                         [RequestId]
	FROM CommunityApprovals CA
	INNER JOIN OrganizationApprovalRequests AS OAR ON OAR.RequestId = CA.Id
	INNER JOIN Organizations O ON O.Id = OAR.OrganizationId
  INNER JOIN RegionalOrganizations RO ON RO.Id = O.Region
	INNER JOIN Users UC ON O.CreatedBy = UC.UserPrincipalName
	WHERE
		CA.ApprovalSystemGUID IS NULL
		AND DATEDIFF(MI, CA.Created, GETDATE()) >=5
END