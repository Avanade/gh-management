CREATE PROCEDURE PR_UserAccess_Insert
	-- Add the parameters for the stored procedure here
( 
            @ProjectId INT
           ,@UserPrincipalName VARCHAR(100)

)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

INSERT INTO [dbo].[UserAccess]
           ( 
            [ProjectId],
            [UserPrincipalName],
            [Created],
            [CreatedBy],
            [Modified],
            [ModifiedBy]
            )
     VALUES
           ( 
            @ProjectId,
            @UserPrincipalName,
            GETDATE(),
            @UserPrincipalName,
            GETDATE(),
            @UserPrincipalName
            )
END