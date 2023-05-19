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
FROM CommunityApprovals CA
INNER JOIN Communities C ON CA.CommunityId = C.Id
INNER JOIN Users UC ON C.CreatedBy = UC.UserPrincipalName
INNER JOIN ApprovalStatus S ON S.Id = CA.ApprovalStatusId
WHERE C.Id = @Id

