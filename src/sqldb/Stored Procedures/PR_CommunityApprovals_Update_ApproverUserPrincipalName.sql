CREATE PROCEDURE [dbo].[PR_CommunityApprovals_Update_ApproverUserPrincipalName]
	@Id uniqueidentifier,
	@ApproverEmail varchar(100),
 	@Username VARCHAR(100)
AS
BEGIN
    SET NOCOUNT ON


	UPDATE [CommunityApprovals]
	SET  
			ApproverUserPrincipalName = @ApproverEmail,
			Modified = GETDATE(),
			ModifiedBy = @Username
	WHERE ApprovalSystemguid= @Id
	
	SELECT 

		CA.ApproverUserPrincipalName,
		CA.Id, C.[Id] AS CommunityId, C.[Name] [ProjectName],
		C.[Description] [ProjectDescription],
		U1.Name [RequesterName], U1.GivenName [RequesterGivenName], U1.SurName [RequesterSurName], U1.UserPrincipalName [RequesterUserPrincipalName],
		CA.[ApprovalStatusId], T.[Name] ApprovalType,
		C.Url,	C.Notes,
		CA.ApprovalDescription,
		S.Name [RequestStatus],
		CA.[ApprovalDate], CA.[ApprovalRemarks]
	
		
	FROM 
		[CommunityApprovals] CA
		INNER JOIN ApprovalTypes T ON CA.[ApprovalStatusId] = T.Id
		INNER JOIN [dbo].[CommunityApprovalRequests] CAR ON CA.Id = CAR.RequestId
		INNER JOIN [dbo].[Communities] C ON CAR.[CommunityId] = C.Id
		INNER JOIN Users U1 ON CA.CreatedBy = U1.UserPrincipalName
		INNER JOIN ApprovalStatus S ON S.Id = CA.ApprovalStatusId
	WHERE 
		CA.ApprovalSystemguid = @Id
END