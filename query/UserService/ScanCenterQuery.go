package query

var GetScanCenterCountSQL = `
SELECT
  COUNT(*) as "TotalCount"
FROM
  public."ScanCenter"
`

var GetScancenterOneDataSQL = `
SELECT
  *
FROM
  public."ScanCenter"
WHERE
  "refSCId" = ?
`

var UpdateScancenterSQL = `
UPDATE
  public."ScanCenter"
SET
  "refSCProfile" = ?,
  "refSCName" = ?,
  "refSCAddress" = ?,
  "refSCPhoneNo1" = ?,
  "refSCEmail" = ?,
  "refSCWebsite" = ?,
  "refSCAppointments" = ?
WHERE
  "refSCId" = ?
`
