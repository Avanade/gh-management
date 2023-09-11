CREATE PROCEDURE PR_Projects_Update_OssSponsorIdById
(
    @Id INT,
    @OSSContributionSponsorId INT
)
AS
BEGIN

UPDATE
    [dbo].[Project]
SET
    [OSSContributionSponsorId] = @OSSContributionSponsorId
WHERE
    [Id] = @Id

END