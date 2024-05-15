CREATE PROCEDURE [dbo].[PR_CommunityApprovals_Select_ById]
(
	@Id INT
)
AS

SELECT
CA.Id,
C.Id [CommunityId],
C.[Name] [CommunityName],
C.Url [CommunityUrl],
C.Description [CommunityDescription],
C.Notes [CommunityNotes],
C.TradeAssocId [CommunityTradeAssocId],
C.CommunityType [CommunityType],
UC.[Name] [RequesterName],
UC.GivenName [RequesterGivenName],
UC.SurName [RequesterSurName],
UC.UserPrincipalName [RequesterUserPrincipalName],
CA.[ApproverUserPrincipalName],
CA.[ApprovalDescription],
s.Name [ApprovalStatus]
FROM CommunityApprovalRequests CAR
LEFT JOIN CommunityApprovals CA ON CAR.RequestId = CA.Id
LEFT JOIN Communities C ON CAR.CommunityId = C.Id
LEFT JOIN Users UC ON C.CreatedBy = UC.UserPrincipalName
LEFT JOIN ApprovalStatus S ON S.Id = CA.ApprovalStatusId
WHERE CAR.CommunityId = @Id

