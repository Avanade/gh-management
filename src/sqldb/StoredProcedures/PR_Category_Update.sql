CREATE PROCEDURE [dbo].[PR_Category_Update]
(
			@Id INT,
			@Name VARCHAR(50),	
			@CreatedBy  VARCHAR(50),
			@ModifiedBy  VARCHAR(50) 
		 
) AS
BEGIN
UPDATE [dbo].[Category]
   SET [Name] = @Name,
       [Created] =GETDATE(),
	   [CreatedBy] = @CreatedBy,
	   [Modified] = GETDATE(),
	   [ModifiedBy] = @ModifiedBy
 WHERE  [Id] = @Id 
END
