CREATE TABLE [dbo].[Applications] (
    [Id]         UNIQUEIDENTIFIER CONSTRAINT [DF_Applications_Id] DEFAULT (newid()) NOT NULL,
    [Name]       VARCHAR (100)    NOT NULL,
    [IsActive]   BIT              CONSTRAINT [DF_Applications_IsActive] DEFAULT ((1)) NOT NULL,
    [Created]    DATETIME         CONSTRAINT [DF_Applications_Created] DEFAULT (getdate()) NOT NULL,
    [CreatedBy]  VARCHAR (255)    NULL,
    [Modified]   DATETIME         CONSTRAINT [DF_Applications_Modified] DEFAULT (getdate()) NOT NULL,
    [ModifiedBy] VARCHAR (255)    NULL,
    CONSTRAINT [PK_Applications] PRIMARY KEY CLUSTERED ([Id] ASC)
);

