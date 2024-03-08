CREATE PROCEDURE [dbo].[PR_Projects_Select_ById]
(
	@Id INT
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
SELECT p.[Id],
       [GithubId],
       p.[Name],
       [CoOwner],
       [Description],
       [ConfirmAvaIP],
       [ConfirmEnabledSecurity],
       [ApprovalStatusId],
       [IsArchived],
       [Created],
       [CreatedBy],
       [Modified],
       [ModifiedBy],
       [TFSProjectReference],
       [RepositorySource],
       [v].[Name] AS "Visibility",
       (SELECT STRING_AGG(r.Topic, ',') FROM dbo.RepoTopics AS r WHERE r.ProjectId=p.Id) AS "Topics"
  FROM 
       [dbo].[Projects] AS p
        LEFT JOIN [dbo].[Visibility] AS v ON p.VisibilityId = v.Id
  WHERE
      p.[Id] = @Id
END