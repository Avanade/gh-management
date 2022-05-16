/*
Post-Deployment Script Template							
--------------------------------------------------------------------------------------
 This file contains SQL statements that will be appended to the build script.		
 Use SQLCMD syntax to include a file in the post-deployment script.			
 Example:      :r .\myfile.sql								
 Use SQLCMD syntax to reference a variable in the post-deployment script.		
 Example:      :setvar TableName MyTable							
               SELECT * FROM [$(TableName)]					
--------------------------------------------------------------------------------------
*/

/* INITIAL DATA FOR APPROVAL TYPES */
SET IDENTITY_INSERT ApprovalTypes ON

IF NOT EXISTS (SELECT Id FROM ApprovalTypes WHERE Id = 1)
INSERT INTO ApprovalTypes (Id, ApproveText, RejectText) VALUES (1, 'Approve', 'Reject')

IF NOT EXISTS (SELECT Id FROM ApprovalTypes WHERE Id = 1)
INSERT INTO ApprovalTypes (Id, ApproveText, RejectText) VALUES (2, 'Accept', 'Decline')

SET IDENTITY_INSERT ApprovalTypes OFF