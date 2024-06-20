CREATE PROCEDURE [dbo].[PR_Communities_Update]
(
			@Id INT,
			@Name VARCHAR(50),
			@Url VARCHAR(255),
			@Description VARCHAR(255),
			@Notes VARCHAR(255),
			@TradeAssocId VARCHAR(255),
            @CommunityType VARCHAR(10),
            @ChannelId VARCHAR(100)=NULL,
            @OnBoardingInstructions VARCHAR(MAX) = NULL,
			@CreatedBy  VARCHAR(50),
			@ModifiedBy  VARCHAR(50)
) AS
BEGIN
UPDATE [dbo].[Communities]
   SET [Name] = @Name
      ,[Url] = @Url
      ,[Description] = @Description
      ,[Notes] = @Notes
      ,[TradeAssocId] = @TradeAssocId
      ,[CommunityType] = @CommunityType
      ,[ChannelId] = @ChannelId
      ,[OnBoardingInstructions]=@OnBoardingInstructions
      ,[Created] =GETDATE()
      ,[CreatedBy] = @CreatedBy
      ,[Modified] = GETDATE()
      ,[ModifiedBy] = @ModifiedBy
 WHERE  [Id] = @Id

DELETE FROM CommunitySponsors WHERE CommunityId = @Id
DELETE FROM CommunityTags WHERE CommunityId = @Id
END