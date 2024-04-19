CREATE PROCEDURE [dbo].[PR_Projects_Select_ByUserPrincipalName]
(
	@UserPrincipalName VARCHAR(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here

SELECT [p].[Id],
       [p].[Name],
       [p].[AssetCode],
       [p].[Organization],
       [CoOwner],
       [Description],
       [ConfirmAvaIP],
       [ConfirmEnabledSecurity],
       (SELECT COUNT(*) FROM ProjectApprovals WHERE ProjectId = P.Id AND RespondedBy IS NULL) [TotalPendingRequest],
       [ApprovalStatusId],
       [IsArchived],
       [Created],
       [CreatedBy],
       [Modified],
       [ModifiedBy],
       [v].[Name] AS 'Visibility',
       [p].[RepositorySource],
	     [p].[TFSProjectReference],
       [p].[ECATTID],
       (SELECT STRING_AGG(r.Topic, ',') FROM dbo.RepoTopics AS r WHERE r.ProjectId=p.Id) AS "Topics"
  FROM 
   [dbo].[RepoOwners] AS RO
   LEFT JOIN [dbo].[Projects] AS p  ON  RO.ProjectId = p.Id
   LEFT JOIN [dbo].[Visibility] AS v ON p.VisibilityId = v.Id
 
  WHERE  
   RO.UserPrincipalName =@UserPrincipalName
  ORDER BY [Created] DESC
END