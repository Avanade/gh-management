CREATE TABLE [dbo].[ApprovalTypes] (
    [Id]          INT          IDENTITY (1, 1) NOT NULL,
    [ApproveText] VARCHAR (50) NOT NULL,
    [RejectText]  VARCHAR (50) NOT NULL,
    CONSTRAINT [PK_ApprovalTypes] PRIMARY KEY CLUSTERED ([Id] ASC)
);