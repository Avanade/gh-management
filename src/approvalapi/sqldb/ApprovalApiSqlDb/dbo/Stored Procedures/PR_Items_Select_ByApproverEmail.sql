CREATE PROCEDURE [dbo].[PR_Items_Select_ByApproverEmail]
	@ApproverEmail varchar(100)
AS
	SELECT
		A.[Name] [Application],
		AM.[Name] [Module],
		[Subject], Body, DateSent,
		DateResponded, IsApproved, ApproverRemarks
	FROM Items I
	INNER JOIN ApplicationModules AM ON I.ApplicationModuleId = AM.Id
	INNER JOIN Applications A ON AM.ApplicationId = A.Id
	WHERE I.ApproverEmail = @ApproverEmail