CREATE PROCEDURE [dbo].[PR_CommunityApprovals_Select_ById]
(
	@Id int
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
C.IsExternal [CommunityIsExternal],
UC.[Name] [RequesterName],
UC.GivenName [RequesterGivenName],
UC.SurName [RequesterSurName],
UC.UserPrincipalName [RequesterUserPrincipalName],
CA.[ApproverUserPrincipalName],
CA.[ApprovalDescription]
FROM CommunityApprovals CA
INNER JOIN Communities C ON CA.CommunityId = C.Id
INNER JOIN Users UC ON C.CreatedBy = UC.UserPrincipalName
WHERE C.Id = @Id
