create PROCEDURE [dbo].[PR_RelatedCommunities_Select]
(
 @ParentCommunityId int
 )
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.

SELECT   
[ParentCommunityId],
[RelatedCommunityId],
c.IsExternal,
c.Name
FROM
[dbo].[RelatedCommunities] RC
inner  join 
Communities c
on rc.RelatedCommunityId = c.id
where
ParentCommunityId = @ParentCommunityId

END
