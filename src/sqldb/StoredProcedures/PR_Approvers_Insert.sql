CREATE PROCEDURE  [dbo].[PR_Approvers_Insert]
(
    @ApprovalTypeId INT,
    @ApproverEmail VARCHAR(100)
)
AS
BEGIN  
    SET NOCOUNT ON
    IF NOT EXISTS (
        SELECT * FROM Approvers WHERE ApprovalTypeId = @ApprovalTypeId AND ApproverEmail = @ApproverEmail
    )
    BEGIN 
        INSERT INTO [dbo].[Approvers]
            (
                [ApprovalTypeId] ,
                [ApproverEmail]
            )
        VALUES
            (
                @ApprovalTypeId,
                @ApproverEmail
            )
    END
END