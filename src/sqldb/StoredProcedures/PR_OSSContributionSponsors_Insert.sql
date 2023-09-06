CREATE PROCEDURE PR_OSSContributionSponsors_Insert
	@Name VARCHAR(50)
AS
BEGIN

INSERT INTO
    [dbo].[OSSContributionSponsors] ( 
        [Name]
    )
VALUES
    ( 
        @Name
    )
END