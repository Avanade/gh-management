CREATE TABLE [dbo].[RelatedCommunities](
	[ParentCommunityId] [INT] NOT NULL,
	[RelatedCommunityId] [INT] NOT NULL,
	CONSTRAINT [PK_RelatedCommunities] PRIMARY KEY ([ParentCommunityId], [RelatedCommunityId])
)