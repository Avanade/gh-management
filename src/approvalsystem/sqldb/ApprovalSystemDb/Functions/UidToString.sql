-- ======================================================
-- Create Scalar Function Template for Azure SQL Database
-- ======================================================
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:      <Author, , Name>
-- Create Date: <Create Date, , >
-- Description: <Description, , >
-- =============================================
CREATE FUNCTION UidToString
(
    -- Add the parameters for the function here
    @Uid uniqueidentifier
)
RETURNS varchar(36)
AS
BEGIN
    RETURN CONVERT(VARCHAR(36),@Uid)
END
GO

