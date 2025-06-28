package query

var VerifyAppointment = `
SELECT
  COUNT(*) AS "TotalCount"
FROM
  appointment."refAppointments"
WHERE
  "refSCId" = ?
  AND "refAppointmentDate" = ?;
`

var FindScanCenterSQL = `
SELECT
  *
FROM
  public."ScanCenter"
WHERE
  "refSCCustId" = ?
`

var ViewPatientHistorySQL = `
SELECT
  *
FROM
  appointment."refAppointments" ra
  JOIN public."ScanCenter" sc ON sc."refSCId" = ra."refSCId"
WHERE
  ra."refUserId" = ?
`

var ViewTechnicianPatientQueueSQL = `
SELECT
  *
FROM
  appointment."refAppointments" ra
  JOIN public."ScanCenter" sc ON sc."refSCId" = ra."refSCId"
  JOIN public."Users" u ON u."refUserId" = ra."refUserId"
  JOIN map."refScanCenterMap" rscm ON rscm."refSCId" = ra."refSCId"
WHERE
  rscm."refUserId" = ?
`
