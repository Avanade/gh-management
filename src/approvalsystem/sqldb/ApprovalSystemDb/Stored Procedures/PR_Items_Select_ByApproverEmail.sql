CREATE PROCEDURE [dbo].[PR_Items_Select_ByApproverEmail]
	@ApproverEmail varchar(100)
AS
	SELECT
		dbo.UidToString(A.Id) [ApplicationId], A.[Name] [Application],
		dbo.UidToString(AM.Id) [ApplicationModuleId], AM.[Name] [Module],
		dbo.UidToString(I.Id) [ItemId], [Subject], Body, DateSent,
		DateResponded, IsApproved, ApproverRemarks, I.Created,
		T.ApproveText, T.RejectText
	FROM Items I
	INNER JOIN ApplicationModules AM ON I.ApplicationModuleId = AM.Id
	INNER JOIN Applications A ON AM.ApplicationId = A.Id
	INNER JOIN ApprovalTypes T ON T.Id = AM.ApprovalTypeId
	WHERE I.ApproverEmail = @ApproverEmail