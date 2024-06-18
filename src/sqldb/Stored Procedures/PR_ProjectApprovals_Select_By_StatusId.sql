CREATE PROCEDURE [dbo].[PR_ProjectApprovals_Select_By_StatusId]
(
	@ApprovalStatusId INT
)
AS
BEGIN
  SELECT
	PA.Id AS ProjectApprovalId,
	CONVERT(varchar(36), PA.ApprovalSystemGUID, 1) AS ItemId,
	T.Name AS ApprovalType,
	P.TFSProjectReference AS RepoLink,
	P.Name AS RepoName,
	U.Name AS Requester,
	PA.Created
  FROM
    ProjectApprovals PA
    INNER JOIN Projects P ON P.Id = PA.ProjectId
	INNER JOIN ApprovalTypes T ON T.Id = PA.ApprovalTypeId
	INNER JOIN Users U ON U.UserPrincipalName = PA.CreatedBy
  WHERE  
	PA.ApprovalStatusId = @ApprovalStatusId
END