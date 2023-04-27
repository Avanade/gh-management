create PROCEDURE [dbo].[PR_Projects_ToRepoOwners]
 as
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON

	 select distinct
				Id  ,
				UserPrincipalName  
			 from 
			 (    select distinct 
							Id  ,
							UserPrincipalName   as UserPrincipalName
					from 
							(select id  , [CreatedBy],[CoOwner] from Projects where RepositorySource='GitHub' ) x 
							UNPIVOT  
							(UserPrincipalName FOR Owners IN (CreatedBy, CoOwner)   )AS unpvt
			   ) ProjectOwners
			   union
			   select id  , [CreatedBy]  as owners  from Projects where RepositorySource='AzureDevOps'  
			     
			 
END
