CREATE PROCEDURE [dbo].[usp_Community_Insert]
  @Name [VARCHAR](50),
  @Url [VARCHAR](255),
  @Description [VARCHAR](255),
  @Notes [VARCHAR](255),
  @TradeAssocId [VARCHAR](255),
  @CommunityType [VARCHAR](10),
  @ChannelId [VARCHAR](100) = NULL,
  @OnBoardingInstructions [VARCHAR](MAX) = NULL,
  @CreatedBy [VARCHAR](50),
  @ModifiedBy [VARCHAR](50),
  @Id [INT] = 0
AS
BEGIN
  IF NOT EXISTS (SELECT * FROM [dbo].[Community] WHERE [Id] = @Id)
  BEGIN
    INSERT INTO [dbo].[Community]
    (
      [Name]
      ,[Url]
      ,[Description]
      ,[Notes]
      ,[TradeAssocId]
      ,[CommunityType]
      ,[ChannelId]
      ,[OnBoardingInstructions]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
    )
    VALUES
    (
      @Name
      ,@Url
      ,@Description
      ,@Notes
      ,@TradeAssocId
      ,@CommunityType
      ,@ChannelId
      ,@OnBoardingInstructions
      ,GETDATE()
      ,@CreatedBy
      ,GETDATE()
      ,@ModifiedBy
    )

    SET @Id = SCOPE_IDENTITY()
  END
	ELSE
  BEGIN
    EXEC [dbo].[usp_Community_Update]
      @Id, @Name, @Url, @Description, 
      @Notes, @TradeAssocId, @CommunityType, 
      @ChannelId, @OnBoardingInstructions, @CreatedBy, @ModifiedBy
  END

  SELECT @Id AS Id
END