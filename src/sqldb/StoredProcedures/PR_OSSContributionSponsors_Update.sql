CREATE PROCEDURE PR_ExternalLinks_Update
(
			@Id int,
			@Name VARCHAR(50),
            @IsArchived BIT
)
AS
BEGIN

UPDATE
    [dbo].[OSSContributionSponsors]
SET
    [Name] = @Name,
    [IsArchived] = @IsArchived
WHERE
    [Id] = @Id

END