create PROCEDURE [dbo].[PR_Category_Update]
(
			@Id int,
			@Name varchar(50),	
			@CreatedBy  varchar(50),
			@ModifiedBy  varchar(50) 
		 
) AS
BEGIN
UPDATE [dbo].[Category]
   SET [Name] = @Name,
       [Created] =GETDATE(),
	   [CreatedBy] = @CreatedBy,
	   [Modified] = GETDATE(),
	   [ModifiedBy] = @ModifiedBy
 WHERE  [Id] = @Id

 
end
