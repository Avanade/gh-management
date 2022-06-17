/****** Object:  StoredProcedure [dbo].[PR_Communities_select_byID]    Script Date: 6/3/2022 2:22:04 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
create PROCEDURE [dbo].[PR_Communities_select_byID]
@Id int
as 
begin

SELECT [Id]
      ,[Name]
      ,[Url]
      ,[Description]
      ,[Notes]
      ,[TradeAssocId]
      ,[IsExternal]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[Communities]
  where [Id] = @Id


end
