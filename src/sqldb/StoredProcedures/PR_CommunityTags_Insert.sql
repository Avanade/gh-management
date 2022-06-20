-- =======================================================
-- Create Stored Procedure Template for Azure SQL Database
-- =======================================================
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:      <Author, , Name>
-- Create Date: <Create Date, , >
-- Description: <Description, , >
-- =============================================
CREATE PROCEDURE PR_CommunityTags_Insert
(
    -- Add the parameters for the stored procedure here
			@CommunityId int,
			@Tag varchar(20)
            
)
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON

    -- Insert statements for procedure here
  INSERT INTO [dbo].[CommunityTags]
           ([CommunityId]
           ,[Tag])
     VALUES
           (@CommunityId
           ,@Tag)
END
GO
