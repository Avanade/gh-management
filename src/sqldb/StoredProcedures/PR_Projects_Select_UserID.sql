Create PROCEDURE [dbo].[PR_Projects_Select_UserId]
(
	@Id Int
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
SELECT [Id],
       [Name],
       [CoOwner],
       [Description],
       [ConfirmAvaIP],
       [ConfirmEnabledSecurity],
       [Created],
       [CreatedBy],
       [Modified],
       [ModifiedBy]
  FROM 
       [dbo].[Projects]
  WHERE  
       [Id] =@Id
END
