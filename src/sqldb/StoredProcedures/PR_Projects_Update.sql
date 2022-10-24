CREATE PROCEDURE PR_Projects_Update
(
		@Id INT,
		@Name VARCHAR(50),
		@CoOwner VARCHAR(100),
		@Description VARCHAR(1000),
		@ConfirmAvaIP BIT,
		@ConfirmEnabledSecurity BIT,
		@ConfirmNotClientProject BIT,
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
		[Name] = @Name,
		[CoOwner] = @CoOwner,
		[Description] = @Description,
		[ConfirmAvaIP] = @ConfirmAvaIP,
		[ConfirmEnabledSecurity] = @ConfirmEnabledSecurity,
		[ConfirmNotClientProject] = @ConfirmNotClientProject,
		[Modified] = GETDATE(),
		[ModifiedBy] = @ModifiedBy
 WHERE  
		[Id] = @Id
END