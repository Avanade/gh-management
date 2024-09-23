CREATE PROCEDURE [dbo].[usp_ActivityType_Select]
AS
BEGIN
    SET NOCOUNT ON

    SELECT
        [Id],
        [Name]
    FROM [dbo].[ActivityType]
END