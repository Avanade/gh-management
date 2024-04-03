CREATE PROCEDURE [dbo].[PR_CommunityApprovals_Select_FailedRequestCommunities]
AS
BEGIN
	SELECT
		CA.Id							[Id],
		C.Id 	 			 		 	[CommunityId],
		C.Name	 			 		 	[CommunityName],
		C.Url 				 			[CommunityUrl],
		C.Description 		 			[CommunityDescription],
		C.Notes 			 			[CommunityNotes],
		C.TradeAssocId 		 			[CommunityTradeAssocId],
		C.CommunityType 	 			[CommunityType],
		UC.Name 			 		 	[RequesterName],
		UC.GivenName 		 		 	[RequesterGivenName],
		UC.SurName 			 		 	[RequesterSurName],
		UC.UserPrincipalName 		 	[RequesterUserPrincipalName],
		CA.ApproverUserPrincipalName 	[ApproverUserPrincipalName],
		CA.ApprovalDescription			[ApprovalDescription]
	FROM CommunityApprovals CA
	INNER JOIN CommunityApprovalRequests AS CAR ON CAR.RequestId = CA.Id
	INNER JOIN Communities C ON C.Id = CAR.CommunityId
	INNER JOIN Users UC ON C.CreatedBy = UC.UserPrincipalName
	WHERE
		CA.ApprovalSystemGUID IS NULL
		AND DATEDIFF(MI, CA.Created, GETDATE()) >=5
END