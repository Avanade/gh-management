CREATE PROCEDURE PR_ExternalLinks_Select
	@IconSVG VARCHAR(100),
	@SVGName VARCHAR(100),
	@UserPrincipalName VARCHAR(100)

AS
BEGIN

SET NOCOUNT ON;

INSERT INTO [dbo].[ExternalLinks]
           ( 
            [UserPrincipalName],
			[IconSVG],
			[SVGName],
            [Created],
            [CreatedBy],
            [Modified],
            [ModifiedBy]
            )
     VALUES
           ( 
            @UserPrincipalName,
			@IconSVG,
			@SVGName,
            GETDATE(),
            @UserPrincipalName,
            GETDATE(),
            @UserPrincipalName
            )
END
