
CREATE PROCEDURE [dbo].[PR_Organizations_Insert]
(
    @Region INT,
    @ClientName VARCHAR(100),
    @ProjectName VARCHAR(100),
    @WBS VARCHAR(50),
    @CreatedBy VARCHAR(100)
) AS
BEGIN
	DECLARE @returnID AS INT
 
	INSERT INTO [dbo].[Organizations]
        ([Region]
        ,[ClientName]
        ,[ProjectName]
        ,[WBS]
        ,[CreatedBy]
        ,[Created])
    VALUES
        (@Region
        ,@ClientName
        ,@ProjectName
        ,@WBS
        ,@CreatedBy
        ,GETDATE())
    SET @returnID = SCOPE_IDENTITY()

    SELECT @returnID Id
END
