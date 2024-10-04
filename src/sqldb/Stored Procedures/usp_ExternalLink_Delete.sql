CREATE PROCEDURE [dbo].[usp_ExternalLink_Delete]
  @Id [INT]
AS
BEGIN
  DELETE FROM [dbo].[ExternalLink] WHERE [Id] = @Id
END