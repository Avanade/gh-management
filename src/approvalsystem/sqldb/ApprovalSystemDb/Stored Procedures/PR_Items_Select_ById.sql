CREATE PROCEDURE [dbo].[PR_Items_Select_ById]
	@Id UNIQUEIDENTIFIER
AS
	SELECT
		A.[Name] [Application],
		AM.[Name] [Module],
		[Subject], Body, DateSent,
		DateResponded, IsApproved, ApproverRemarks,
		T.ApproveText, T.RejectText,
		AM.CallbackUrl
	FROM Items I
	INNER JOIN ApplicationModules AM ON I.ApplicationModuleId = AM.Id
	INNER JOIN Applications A ON AM.ApplicationId = A.Id
	INNER JOIN ApprovalTypes T ON AM.ApprovalTypeId = T.Id
	WHERE I.Id = @Id