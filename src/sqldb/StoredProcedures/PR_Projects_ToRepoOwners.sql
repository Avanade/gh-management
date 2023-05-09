CREATE PROCEDURE [dbo].[PR_Projects_ToRepoOwners]
 AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON

	 SELECT DISTINCT
				Id  ,
				UserPrincipalName  
			 FROM 
			 (    SELECT DISTINCT  		
							Id  ,
							UserPrincipalName   AS UserPrincipalName
					FROM 
							(SELECT id  , [CreatedBy],[CoOwner] FROM Projects WHERE RepositorySource='GitHub' ) x 
							UNPIVOT  
							(UserPrincipalName FOR Owners IN (CreatedBy, CoOwner)   )AS unpvt
			   ) ProjectOwners
			   UNION
			   SELECT id  , [CreatedBy]  AS owners  FROM Projects WHERE RepositorySource='AzureDevOps'  
			     
			 
END
