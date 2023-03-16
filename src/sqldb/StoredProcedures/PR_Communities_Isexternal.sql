
create PROCEDURE  [dbo].[PR_Communities_Isexternal]
(
    -- Add the parameters for the stored procedure here
    @isexternal int,
	@UserPrincipalName varchar(100)
)
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
		select 
		c.id , 
		c.name  
		FROM 
			[dbo].[Communities] c
			INNER JOIN ApprovalStatus T ON c.ApprovalStatusId = T.Id
			 
		WHERE 
	
			isexternal =@isexternal 
			and  	
			(	
			c.CreatedBy = @UserPrincipalName
			or
			 T.Id =5)
END

