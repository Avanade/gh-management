-- This file contains SQL statements that will be executed after the build script.
/* INITIAL DATA FOR APPROVAL TYPES */
SET IDENTITY_INSERT ApprovalTypes ON

IF NOT EXISTS (SELECT Id FROM ApprovalTypes WHERE Id = 1)
INSERT INTO ApprovalTypes (Id, ApproveText, RejectText) VALUES (1, 'Approve', 'Reject')

IF NOT EXISTS (SELECT Id FROM ApprovalTypes WHERE Id = 1)
INSERT INTO ApprovalTypes (Id, ApproveText, RejectText) VALUES (2, 'Accept', 'Decline')

SET IDENTITY_INSERT ApprovalTypes OFF