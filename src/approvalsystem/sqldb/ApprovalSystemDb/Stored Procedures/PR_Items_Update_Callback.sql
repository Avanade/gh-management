CREATE PROCEDURE PR_Items_Update_Callback
	@ItemId UNIQUEIDENTIFIER,
	@IsCallbackFailed BIT
AS

UPDATE Items
SET IsCallbackFailed = @IsCallbackFailed,
LastCallbackAttemptDate = GETDATE(),
CallbackAttemptCount = CallbackAttemptCount+1
WHERE Id = @ItemId