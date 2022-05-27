CREATE TABLE [dbo].[Items] (
    [Id]                  UNIQUEIDENTIFIER CONSTRAINT [DF_Items_Id] DEFAULT (newid()) NOT NULL,
    [ApplicationModuleId] UNIQUEIDENTIFIER NOT NULL,
    [ApproverEmail]       VARCHAR (100)    NOT NULL,
    [Subject]             VARCHAR (100)    NULL,
    [Body]             VARCHAR (8000)    NULL,
    [DateSent]            DATETIME         NULL,
    [DateResponded]       DATETIME         NULL,
    [IsApproved]          BIT              NULL,
    [ApproverRemarks]     VARCHAR (255)    NULL,
    [IsCallbackFailed]    BIT               NULL,
    [LastCallbackAttemptDate] DATETIME  NULL,
    [CallbackAttemptCount] int  NOT NULL DEFAULT 0,
    [Created]             DATETIME         CONSTRAINT [DF_Items_Created] DEFAULT (getdate()) NOT NULL,
    [CreatedBy]           VARCHAR (255)    NULL,
    [Modified]            DATETIME         CONSTRAINT [DF_Items_Modified] DEFAULT (getdate()) NOT NULL,
    [ModifiedBy]          VARCHAR (255)    NULL,
    CONSTRAINT [PK_Items] PRIMARY KEY CLUSTERED ([Id] ASC),
    CONSTRAINT [FK_Items_ApplicationModules] FOREIGN KEY ([ApplicationModuleId]) REFERENCES [dbo].[ApplicationModules] ([Id])
);

