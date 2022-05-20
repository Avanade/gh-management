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

