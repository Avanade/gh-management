CREATE PROCEDURE PR_Items_Select_FailedCallbacks

AS

SELECT dbo.UidToString(Id) [Id] FROM Items WHERE IsCallbackFailed = 1