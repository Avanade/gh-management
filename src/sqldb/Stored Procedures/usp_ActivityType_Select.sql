CREATE PROCEDURE [dbo].[usp_ActivityType_Select]
AS
BEGIN
    SET NOCOUNT ON

    SELECT
        [Id],
        [Name],
        [Created],
        [CreatedBy],
        [Modified],
        [ModifiedBy]
    FROM [dbo].[ActivityType]
END