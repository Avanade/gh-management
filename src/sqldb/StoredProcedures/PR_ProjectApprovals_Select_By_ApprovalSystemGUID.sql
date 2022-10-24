CREATE PROCEDURE [dbo].[PR_ProjectApprovals_Select_By_ApprovalSystemGUID]
  (
  @ApprovalSystemGUID UNIQUEIDENTIFIER
)
AS
BEGIN
  SELECT
    PA.Id, PA.ProjectId,
    P.[Name] [ProjectName],
    PA.ApprovalTypeId, T.[Name] ApprovalType,
    PA.ApproverUserPrincipalName,
    PA.ApprovalDescription,
    S.Name [RequestStatus],
    PA.[ApprovalDate], PA.[ApprovalRemarks]
  FROM
    ProjectApprovals PA
    INNER JOIN Projects P ON PA.ProjectId = P.Id
    INNER JOIN ApprovalTypes T ON PA.ApprovalTypeId = T.Id
    INNER JOIN ApprovalStatus S ON S.Id = PA.ApprovalStatusId
  WHERE  
    PA.[ApprovalSystemGUID] = @ApprovalSystemGUID
END
