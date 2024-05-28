CREATE TABLE [dbo].[RelatedCommunity](
	[ParentCommunityId] [INT] NOT NULL,
	[RelatedCommunityId] [INT] NOT NULL,
	CONSTRAINT [PK_RelatedCommunity] PRIMARY KEY ([ParentCommunityId], [RelatedCommunityId])
)