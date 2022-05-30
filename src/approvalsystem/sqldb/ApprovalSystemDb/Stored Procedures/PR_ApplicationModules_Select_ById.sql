CREATE PROCEDURE PR_ApplicationModules_Select_ById

	@Id UNIQUEIDENTIFIER

AS

SELECT
dbo.UidToString(A.Id) [ApplicationId], A.[Name] [ApplicationName],
dbo.UidToString(AM.Id) [ApplicationModuleId], AM.[Name] [ApplicationModuleName],
AM.CallbackUrl, AM.RequireRemarks, AM.ApprovalTypeId, [AT].ApproveText, [AT].RejectText
FROM Applications A
INNER JOIN ApplicationModules AM ON A.Id = AM.ApplicationId
INNER JOIN ApprovalTypes [AT] ON AM.ApprovalTypeId = AT.Id
WHERE
A.IsActive = 1
AND AM.IsActive = 1
AND AM.Id = @Id