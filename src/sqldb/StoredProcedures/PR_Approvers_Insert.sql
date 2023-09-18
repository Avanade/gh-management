CREATE PROCEDURE  [dbo].[PR_Approvers_Insert]
(
    @ApprovalTypeId INT,
    @ApproverEmail VARCHAR(100)
)
AS
BEGIN   
    INSERT INTO [dbo].[Approvers]
        (
            [ApprovalTypeId],
            [ApproverEmail]
        )
    VALUES
        (
            @ApprovalTypeId,
            @ApproverEmail
        )
END
