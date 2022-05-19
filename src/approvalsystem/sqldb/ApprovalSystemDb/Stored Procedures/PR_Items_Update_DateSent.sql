CREATE PROCEDURE PR_Items_Update_DateSent

	@Id uniqueidentifier

AS

UPDATE Items
SET
DateSent = GETDATE(),
Modified = GETDATE()
WHERE Id = @Id