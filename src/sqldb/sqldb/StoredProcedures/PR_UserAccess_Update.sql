CREATE PROCEDURE PR_UserAccess_Update
(	   
        @Id  int ,
        @ProjectId	int ,
        @Username	varchar(100) 

)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

UPDATE [dbo].[UserAccess]
   SET [Id] =@Id,
       [ProjectId] =@ProjectId,
       [Username] = @Username,
       [Created] = GETDATE(),
       [CreatedBy] = @Username,
       [Modified] =GETDATE(),
       [ModifiedBy] = @Username
    WHERE  
	   [Id] =@Id
END