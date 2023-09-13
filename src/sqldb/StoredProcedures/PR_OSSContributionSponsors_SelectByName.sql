CREATE PROCEDURE PR_OSSContributionSponsors_SelectByName
(
    @Name VARCHAR(50)
)
AS
BEGIN

SELECT * FROM [dbo].[OSSContributionSponsors]
WHERE [Name]=@Name

END