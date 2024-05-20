CREATE PROCEDURE [dbo].[PR_ProjectApprovals_Select_By_ProjectId]
  (
  @Id INT
)
AS
BEGIN
  SELECT
      PA.Id, PA.ProjectId,
      P.[Name] [ProjectName],
      PA.ApprovalTypeId, T.[Name] ApprovalType,
      PA.ApprovalDescription,
      S.Name [RequestStatus],
      PA.[ApprovalDate], PA.[ApprovalRemarks]
    FROM
      ProjectApprovals PA
      INNER JOIN Projects P ON PA.ProjectId = P.Id
      INNER JOIN ApprovalTypes T ON PA.ApprovalTypeId = T.Id
      INNER JOIN ApprovalStatus S ON S.Id = PA.ApprovalStatusId
    WHERE
      PA.Created = (SELECT TOP(1) Created FROM ProjectApprovals WHERE ProjectId = @Id ORDER BY Created DESC)
    ORDER BY PA.Created DESC
END