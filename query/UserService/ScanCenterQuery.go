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
  "refSCAppointments" = ?,
  "refSCStatus" = ?,
  "refSCConsultantStatus" = ?,
  "refSCConsultantLink" = ?
WHERE
  "refSCId" = ?
`

var ScanCenterInactiveSQL = `
UPDATE
  public."Users" u
SET
  "refUserStatus" = false
FROM
  map."refScanCenterMap" rscm
WHERE
  u."refUserId" = rscm."refUserId"
  AND rscm."refSCId" = $1;
`
