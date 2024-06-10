CREATE PROCEDURE [dbo].[usp_RepositoryApproval_Populate]
  @RepositoryId [INT],
	@RequestedBy [VARCHAR](100)
AS
BEGIN
  INSERT INTO [dbo].[RepositoryApproval]
  (
		[RepositoryId],
		[RepositoryApprovalTypeId],
		[ApprovalStatusId],
		[ApprovalDescription],
		[CreatedBy]
	)
  SELECT 
    @RepositoryId, 
    [T].[Id], 
    1, 
    'For Review - ' + [T].[Name], 
    @RequestedBy
  FROM [dbo].[RepositoryApprovalType] AS [T]
  WHERE 
    [T].[IsActive] = 1 AND 
    [T].[IsArchived] = 0

  INSERT INTO ApprovalRequestApprover
  (
    [RepositoryApprovalId],
    [ApproverUserPrincipalName]
  )
  SELECT
        [RA].[Id],
        [A].[ApproverUserPrincipalName]
  FROM [dbo].[RepositoryApproval] AS [RA] 
    INNER JOIN [dbo].[RepositoryApprovalType] AS [T] ON [RA].[RepositoryApprovalTypeId] = [T].[Id]
    INNER JOIN [dbo].[RepositoryApprover] AS [A] ON [T].[Id] = [A].[RepositoryApprovalTypeId]
  WHERE
    [RA].[RepositoryId] = @RepositoryId AND
    [RA].[ApprovalStatusId] = 1

  UPDATE [dbo].[Repository] 
  SET 
    [ApprovalStatusId] = 2,
    [Modified] = GETDATE() 
  WHERE [Id] = @RepositoryId

  EXEC usp_RepositoryApproval_Select_ById @RepositoryId
END
