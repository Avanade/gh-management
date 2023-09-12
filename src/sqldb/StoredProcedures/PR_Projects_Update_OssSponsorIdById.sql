CREATE PROCEDURE PR_Projects_Update_OssSponsorIdById
(
    @Id INT,
    @OSSContributionSponsorId INT
)
AS
BEGIN

UPDATE
    [dbo].[Projects]
SET
    [OSSContributionSponsorId] = @OSSContributionSponsorId
WHERE
    [Id] = @Id AND [RepositorySource] = 'GitHub' AND [OSSContributionSponsorId] IS NULL

END