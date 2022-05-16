CREATE PROCEDURE PR_UserAccess_Insert
	-- Add the parameters for the stored procedure here
( 
            @ProjectId int
           ,@Username varchar(100)

)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

INSERT INTO [dbo].[UserAccess]
           ( 
            [ProjectId],
            [Username],
            [Created],
            [CreatedBy],
            [Modified],
            [ModifiedBy]
            )
     VALUES
           ( 
            @ProjectId,
            @Username,
            GETDATE(),
            @Username,
            GETDATE(),
            @Username
            )
END