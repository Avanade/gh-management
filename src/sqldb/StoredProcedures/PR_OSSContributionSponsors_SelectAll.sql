CREATE PROCEDURE PR_OSSContributionSponsors_SelectAll
AS
BEGIN

SELECT * FROM [dbo].[OSSContributionSponsors]

ORDER BY [Id] ASC

END