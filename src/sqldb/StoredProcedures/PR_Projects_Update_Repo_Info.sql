CREATE PROCEDURE PR_Projects_Update_Repo_Info
(
		@Id INT,
		@GithubId INT,
		@Name VARCHAR(50),
		@Description VARCHAR(1000),
		@IsArchived BIT,
		@VisibilityId INT,
		@Created DATETIME
)
AS
BEGIN
	SET NOCOUNT ON;
UPDATE 
		[dbo].[Projects]
   SET 
		[Name] = @Name,
		[GithubId] = @GithubId,
		[Description] = @Description,
		[IsArchived] = @IsArchived,
		[VisibilityId] = @VisibilityId,
		[Created] = @Created,
		[Modified] = GETDATE()
 WHERE  
		[Id] = @Id
END