CREATE PROCEDURE [dbo].[PR_Projects_UpdateOrganization_All]
(
	@PrivateOrg  VARCHAR(100), -- VISIBILITY = 1
	@InternalOrg  VARCHAR(100), -- VISIBILITY = 2
	@PublicOrg VARCHAR(100), -- VISIBILITY = 3
) AS
BEGIN

	UPDATE dbo.Projects
	SET Organization = @PrivateOrg,
	WHERE  [Visibility] = 1;

	UPDATE dbo.Projects
	SET Organization = @InternalOrg,
	WHERE  [Visibility] = 2;

	UPDATE dbo.Projects
	SET Organization = @PublicOrg,
	WHERE  [Visibility] = 3;

END
