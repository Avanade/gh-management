ALTER PROCEDURE PR_Communities_Select_Related

	@CommunityId INT

AS

SELECT DISTINCT
CR.Id, CR.[Name], CR.[Url], CR.IsExternal
FROM Communities C
INNER JOIN CommunityTags CT ON C.Id = CT.CommunityId
LEFT JOIN CommunityTags CTR ON CT.Tag = CTR.Tag
LEFT JOIN Communities CR ON CTR.CommunityId = CR.Id AND CR.Id <> @CommunityId
WHERE C.Id = @CommunityId
AND CR.Id IS NOT NULL