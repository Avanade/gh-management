
CREATE PROCEDURE [dbo].[PR_ContributionAreas_Update_ById]
(
    -- Add the parameters for the stored procedure here
		@Id INT,
		@Name VARCHAR(100),
        @ModifiedBy VARCHAR(100)
)
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON

    -- Insert statements for procedure here
   
UPDATE [dbo].[ContributionAreas]
   SET [Name] =@Name,
       [Modified] = GETDATE(),
       [ModifiedBy] =  @ModifiedBy
 WHERE [Id] =@Id
END
