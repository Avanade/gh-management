/****** Object:  StoredProcedure [dbo].[PR_Communities_select]    Script Date: 6/3/2022 2:24:06 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
create PROCEDURE [dbo].[PR_Communities_select]
as 
begin

SELECT [Id]
      ,[Name]
      ,[Url]
      ,[Description]
      ,[Notes]
      ,[TradeAssocId]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[Communities]



end
