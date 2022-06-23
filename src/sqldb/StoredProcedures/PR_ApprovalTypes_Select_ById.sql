/****** Object:  StoredProcedure [dbo].[PR_ActivityTypes_Insert]    Script Date: 6/23/2022 12:30:26 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[PR_ApprovalTypes_Select_ById] 
(
	@Id INT
)
AS
BEGIN
	SELECT * FROM [dbo].[ApprovalTypes] WHERE Id = @Id
END