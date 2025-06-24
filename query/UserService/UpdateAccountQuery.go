package query

var GetOneScanCenterSQL = `
SELECT
  *
FROM
  public."ScanCenter"
WHERE
  "refSCId" = ?
`

var UpdateScanCenterSQL = `
UPDATE
  "public"."ScanCenter"
SET
  "reSCCity" = ?,
  "refSCDoorNo" = ?,
  "refSCEmail" = ?,
  "refSCName" = ?,
  "refSCPhoneNo1CountryCode" = ?,
  "refSCPhoneNo1" = ?,
  "refSCPhoneNo2CountryCode" = NULL,
  "refSCPhoneNo2" = NULL,
  "refSCPincode" = ?,
  "refSCProfile" = ?,
  "refSCRegNo" = ?,
  "refSCState" = ?,
  "refSCStreet" = ?,
  "refSCWebsite" = ?,
  "refSCAppointments" = ?
WHERE
  "refSCId" = ?;
`

var UpdateWorkingHoursSQL = `
UPDATE
  "centerdomain"."refWorkingHours"
SET
  "refWHDay" = ?,
  "refWHFromTime" = ?,
  "refWHToTime" = ?,
  "refWHLeave" = ?
WHERE
  "refWHId" = ?
`

var DeleteWorkingHoursSQL = `
DELETE FROM
  "centerdomain"."refWorkingHours"
WHERE
  "refWHId" = ?
`

var UpdateHolidaysSQL = `
UPDATE
  "centerdomain"."refHolidays"
SET
  "refHODate" = ?,
  "refHOFromTime" = ?,
  "refHOToTime" = ?,
  "refHODesc" = ?
WHERE
  "refHOId" = ?
`

var DeleteHolidaysSQL = `
DELETE FROM
  "centerdomain"."refHolidays"
WHERE
  "refHOId" = ?
`
var DeleteScanCenterSQL = `
DELETE FROM
  "scan"."refScanCenterScanMap"
WHERE
  "refSCSMId" = ?
`

var CheckOneRegNoSQL = `
SELECT * FROM public."ScanCenter" WHERE "refSCRegNo" = ? AND "refSCId" NOT IN (?) 
`
var CheckOnePhoneNoSQL = `
SELECT * FROM public."ScanCenter" WHERE "refSCPhoneNo1" = ? AND "refSCId" NOT IN (?)
`
var ChecOneEmailSQL = `
SELECT * FROM public."ScanCenter" WHERE "refSCEmail" = ? AND "refSCId" NOT IN (?)
`

var InsertTransactionDataSQL = `
WITH input_data AS (
  SELECT
    ?::integer AS transTypeId,
    ?::integer AS refUserId,
    ?::integer AS refTHActionBy,
    jsonb_array_elements(
      ?::jsonb
    ) AS refTHData
)
INSERT INTO "aduit"."refTransHistory" (
  "transTypeId", "refTHData", "refUserId", "refTHActionBy"
)
SELECT
  transTypeId, refTHData, refUserId, refTHActionBy
FROM input_data;
`
var GetOneWorkingHoursSQL = `
SELECT
  *
FROM
  "centerdomain"."refWorkingHours"
WHERE
  "refWHId" = ?
`

var GetOneHolidaysSQL = `
SELECT
  *
FROM
  "centerdomain"."refHolidays"
WHERE
  "refHOId" = ?
`

var GetOneScanSQL = `
SELECT
  *
FROM
  scan."refScanCenterScanMap"
WHERE
  "refSCSMId" = ?
`
