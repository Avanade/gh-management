CREATE PROCEDURE PR_OSSContributionSponsors_SelectByIsArchived
(
    @IsArchived BIT = 0
)
AS
BEGIN

SELECT * FROM [dbo].[OSSContributionSponsors]
WHERE [IsArchived]=@IsArchived

ORDER BY [Id] ASC

END