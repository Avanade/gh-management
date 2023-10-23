CREATE PROCEDURE [dbo].[PR_ApprovalRequests_Select_ByApproverEmail] 
(
	@ApproverEmail VARCHAR(100)
)
AS
BEGIN
	SELECT
        CONVERT(varchar(36), PA.ApprovalSystemGUID, 1) AS ApprovalSystemGUID,
        P.Name AS ProjectName,
        PA.CreatedBy AS RequestedBy,
        P.Description AS Description,
        P.Newcontribution AS NewContribution,
        OCS.Name AS OSSContributionSponsor,
        P.Avanadeofferingsassets AS IsAvanadeOfferingAssets,
        P.Willbecommercialversion AS WillBeCommercialVersion,
        P.OSSContributionInformation AS OSSContributionInformation,
        PA.ApprovalRemarks AS Remarks,
        A.Name AS Status,
        PA.RespondedBy AS RespondedBy,
        PA.ApprovalDate AS ApprovalDate,
        PA.ApprovalSystemDateSent AS ApprovalSystemDateSent
    FROM
        [dbo].[Projects] AS P
        INNER JOIN [dbo].[ProjectApprovals] AS PA ON P.Id = PA.ProjectId
        INNER JOIN [dbo].[ApprovalStatus] AS A ON PA.ApprovalStatusId = A.Id
        INNER JOIN [dbo].[OSSContributionSponsors] AS OCS ON P.OSSContributionSponsorId = OCS.Id
        INNER JOIN [dbo].[ApprovalRequestApprovers] AS ARA ON PA.Id = ARA.ApprovalRequestId
    WHERE
        ARA.ApproverEmail = @ApproverEmail
END