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
  u."refUserCustId" AS "refUserCustId",
  u."refUserFirstName" AS "refUserFirstName",
  u."refUserId" AS "refUserId",
  ra.*,
  sc.*
FROM
  appointment."refAppointments" ra
  JOIN public."ScanCenter" sc ON sc."refSCId" = ra."refSCId"
  JOIN public."Users" u ON u."refUserId" = ra."refUserId"
  JOIN map."refScanCenterMap" rscm ON rscm."refSCId" = ra."refSCId"
WHERE
  rscm."refUserId" = ?
`

var InsertAdditionalFiles = `
WITH input_data AS (
  SELECT
    ?::int AS refUserId,
    ?::int AS refAppointmentId,
    ?::boolean AS refADStatus,
    NOW() AS refADCreatedAt,
    jsonb_array_elements(?::jsonb) AS file
)
INSERT INTO notes."refAddtionalDoc" (
  "refUserId",
  "refAppointmentId",
  "refADFileName",
  "refADOldFileName",
  "refADStatus",
  "refADCreatedAt"
)
SELECT
  refUserId,
  refAppointmentId,
  file ->> 'fileName',
  file ->> 'oldFileName',
  refADStatus,
  refADCreatedAt
FROM input_data;
`

var ViewAddtionalFilesSQL = `
SELECT
  *
FROM
  notes."refAddtionalDoc"
WHERE
  "refUserId" = ?
  AND "refAppointmentId" = ?
`
