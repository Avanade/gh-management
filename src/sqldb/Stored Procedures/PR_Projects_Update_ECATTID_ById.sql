CREATE PROCEDURE PR_Projects_Update_ECATTID_ById
(
		@Id INT,
		@ECATTID INT,
		@ModifiedBy VARCHAR(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
UPDATE 
		[dbo].[Projects]
   SET 
		[ECATTID] = @ECATTID,
		[Modified] = GETDATE(),
		[ModifiedBy] = @ModifiedBy
 WHERE  
		[Id] = @Id
END

