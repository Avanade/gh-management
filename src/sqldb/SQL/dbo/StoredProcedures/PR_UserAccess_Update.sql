CREATE PROCEDURE PR_UserAccess_Update
(	   @Id  int 
      ,@ProjectId	int 
      ,@Username	varchar(100) 
      ,@Created		datetime 
      ,@CreatedBy	varchar(50) 
      ,@Modified	datetime 
	  ,@ModifiedBy  varchar(50) 
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

UPDATE [dbo].[UserAccess]
   SET [Id] =@Id
      ,[ProjectId] =@ProjectId
      ,[Username] = @Username
      ,[Created] = @Created
      ,[CreatedBy] = @CreatedBy
      ,[Modified] =@Modified
      ,[ModifiedBy] = @ModifiedBy
 WHERE 
	  [Id] =@Id
END