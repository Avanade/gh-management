CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Update_ApprovalSystemGUID]
    @Id [INT],
    @ApprovalSystemGUID [UNIQUEIDENTIFIER]
AS
BEGIN
    UPDATE 
        [dbo].[ApprovalRequest]
    SET
        [ApprovalStatusId] = 2,
        [ApprovalSystemGUID] = @ApprovalSystemGUID,
        [ApprovalSystemDateSent] = GETDATE(),
        [Modified] = GETDATE()
    WHERE [Id] = @Id
END