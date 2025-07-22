package query

var VerifyAppointment = `
SELECT
  COUNT(*) AS "TotalCount"
FROM
  appointment."refAppointments"
WHERE
  "refSCId" = ?
  AND "refAppointmentDate" = ?
  AND "refUserId" = ?
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
  AND rscm."refSCId" = ?
`

var ViewAllPatientQueueSQL = `
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
`

var InsertAdditionalFiles = `
WITH input_data AS (
  SELECT
    ?::int AS refUserId,
    ?::int AS refAppointmentId,
    ?::boolean AS refADStatus,
    ? AS refADCreatedAt,
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

var GetUserWithScanDetails = `
SELECT
  *
FROM
  public."Users" u
  JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
WHERE
  u."refRTId" = ?
  AND rscm."refSCId" = ?
  AND (
    u."refRTId" = 2
  )
`

var GetUserDetails = `
SELECT
  *
FROM
  public."Users"
WHERE
  "refRTId" IN (?)
`

var IdentifyScanCenterWithUser = `
SELECT
  *
FROM
  map."refScanCenterMap"
WHERE
  "refUserId" = ?
`

var UpdateAssignUser = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentAssignedUserId" = ?
WHERE
  "refAppointmentId" = ?
`

var CorrectEditStatusSQL = `
SELECT
  (rrh."refRHHandleCorrect" = 1) AS "isHandleCorrect",
  (rrh."refRHHandleEdit" = 1) AS "isHandleEdited"
FROM
  notes."refReportsHistory" rrh
WHERE
  rrh."refUserId" = ?
  AND rrh."refAppointmentId" = ?
  AND rrh."refRHHandledUserId" = ?
`
