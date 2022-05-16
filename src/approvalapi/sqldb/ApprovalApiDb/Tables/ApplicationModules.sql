CREATE TABLE [dbo].[ApplicationModules] (
    [Id]             UNIQUEIDENTIFIER CONSTRAINT [DF_ApplicationModules_Id] DEFAULT (newid()) NOT NULL,
    [ApplicationId]  UNIQUEIDENTIFIER NOT NULL,
    [Name]           VARCHAR (100)    NOT NULL,
    [IsActive]       BIT              CONSTRAINT [DF_ApplicationModules_IsActive] DEFAULT ((1)) NOT NULL,
    [CallbackUrl]    VARCHAR (100)    NULL,
    [RequireRemarks] BIT              CONSTRAINT [DF_ApplicationModules_RequireRemarks] DEFAULT ((1)) NOT NULL,
    [ApprovalTypeId] INT              NULL,
    [Created]        DATETIME         CONSTRAINT [DF_ApplicationModules_Created] DEFAULT (getdate()) NOT NULL,
    [CreatedBy]      VARCHAR (255)    NULL,
    [Modified]       DATETIME         CONSTRAINT [DF_ApplicationModules_Modified] DEFAULT (getdate()) NOT NULL,
    [ModifiedBy]     VARCHAR (255)    NULL,
    CONSTRAINT [PK_ApplicationModules] PRIMARY KEY CLUSTERED ([Id] ASC),
    CONSTRAINT [FK_ApplicationModules_Applications] FOREIGN KEY ([ApplicationId]) REFERENCES [dbo].[Applications] ([Id]),
    CONSTRAINT [FK_ApplicationModules_ApprovalTypes] FOREIGN KEY ([ApprovalTypeId]) REFERENCES [dbo].[ApprovalTypes] ([Id])
);

