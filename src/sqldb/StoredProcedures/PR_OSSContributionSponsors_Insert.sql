CREATE PROCEDURE PR_OSSContributionSponsors_Insert
	@Name VARCHAR(50),
    @IsArchived BIT = 0
AS
BEGIN

INSERT INTO
    [dbo].[OSSContributionSponsors] ( 
        [Name],
        [IsArchived]
    )
VALUES
    ( 
        @Name,
        @IsArchived
    )
END