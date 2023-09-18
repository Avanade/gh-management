CREATE PROCEDURE  [dbo].[PR_ApprovalRequestApprovers_Insert]
(
    @ApprovalRequestId INT,
    @ApproverEmail VARCHAR(100)
)
AS
BEGIN   
    INSERT INTO [dbo].[ApprovalRequestApprovers]
        (
            [ApprovalRequestId],
            [ApproverEmail]
        )
    VALUES
        (
            @ApprovalRequestId,
            @ApproverEmail
        )
END
