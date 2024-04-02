CREATE PROCEDURE [dbo].[PR_Projects_UpdateOrganization_AllAzureDevOps]
AS
BEGIN
    UPDATE [dbo].[Projects]
    SET Organization = SUBSTRING(
                        TFSProjectReference, 
                        CHARINDEX('://', TFSProjectReference) + 3, 
                        CHARINDEX('.', TFSProjectReference, CHARINDEX('://', TFSProjectReference) + 3) - (CHARINDEX('://', TFSProjectReference) + 3))
    FROM [dbo].[Projects]
    WHERE RepositorySource = 'AzureDevOps';
END