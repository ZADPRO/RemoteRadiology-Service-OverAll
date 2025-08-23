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
  u."refUserCustId",
  u."refUserFirstName",
  u."refUserId",
  ra.*,
  sc.*,
  COALESCE(d.dicomFiles, '[]') AS "dicomFiles"
FROM
  appointment."refAppointments" ra
  JOIN public."ScanCenter" sc ON sc."refSCId" = ra."refSCId"
  JOIN public."Users" u ON u."refUserId" = ra."refUserId"
  JOIN map."refScanCenterMap" rscm ON rscm."refSCId" = ra."refSCId"
  LEFT JOIN (
    SELECT
      rdf."refAppointmentId",
      json_agg(rdf.*) AS dicomFiles
    FROM dicom."refDicomFiles" rdf
    WHERE rdf."refDFId" IS NOT NULL
    GROUP BY rdf."refAppointmentId"
  ) d ON d."refAppointmentId" = ra."refAppointmentId"
WHERE
  rscm."refUserId" = ?
  AND rscm."refSCId" = ?
  AND ra."refAppointmentStatus" = true;
`

var ViewAllPatientQueueSQL = `
SELECT
  u."refUserCustId",
  u."refUserFirstName",
  u."refUserId",
  ra.*,
  sc.*,
  COALESCE(dicom_data."dicomFiles", '[]') AS "dicomFiles"
FROM
  appointment."refAppointments" ra
  JOIN public."ScanCenter" sc ON sc."refSCId" = ra."refSCId"
  JOIN public."Users" u ON u."refUserId" = ra."refUserId"
  LEFT JOIN (
    SELECT
      "refAppointmentId",
      json_agg(rdf.*) FILTER (
        WHERE
          rdf."refDFId" IS NOT NULL
      ) AS "dicomFiles"
    FROM
      dicom."refDicomFiles" rdf
    GROUP BY
      rdf."refAppointmentId"
  ) dicom_data ON dicom_data."refAppointmentId" = ra."refAppointmentId"
  WHERE ra."refAppointmentStatus" = true;
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

// var GetUserWithScanDetails = `
// SELECT
//   *
// FROM
//   public."Users" u
//   JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
// WHERE
//   u."refRTId" = ?
//   AND rscm."refSCId" = ?
//   AND (
//     u."refRTId" = 2
//   )
// `

var GetUserWithScanDetails = `
SELECT
  *
FROM
  public."Users" u
  FULL JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
WHERE
  u."refRTId" IN (1, 2, 3, 5, 8, 10)
ORDER BY
  u."refRTId"
`

var GetUserDetails = `
SELECT
  *
FROM
  public."Users" u
  FULL JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
  WHERE u."refRTId" != 4
ORDER BY
  u."refRTId"
`

var IdentifyScanCenterWithUser = `
SELECT
  *
FROM
  map."refScanCenterMap"
WHERE
  "refUserId" = ?
`

var GetAllUserDetailsSQL = `
SELECT 
    uu."refUserCustId" AS "User_Id", 
    ua."refUserCustId" AS "Assigned_Id",
    up."refUserCustId" AS "Patient_Id",
	ra."refAppointmentDate" AS "appointment_date"
FROM public."Users" uu 
JOIN public."Users" ua ON ua."refUserId" = $1
FULL JOIN public."Users" up ON up."refUserId" = $2
FULL JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = $3
WHERE uu."refUserId" = $4
`

var UpdateAssignUser = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentAssignedUserId" = ?
WHERE
  "refAppointmentId" = ?
`

var InsertNotificationSQL = `
INSERT INTO notification.refnotification(
	"refUserId", "refNMessage", "refAppointmentId", "refNAssignedBy", "refNCreatedAt", "refNReadStatus", "refNstatus")
	VALUES ( $1, $2, $3, $4, $5, $6, $7);
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

var GetAllUserSQL = `
SELECT
  *
FROM
  public."Users"
WHERE
  "refUserId" != ?
ORDER BY
  "refRTId"
`

var NotificationSQL = `
SELECT
  *
FROM
  notification."notification"
WHERE
  "refUserId" = $1
`

var GetReportStatusSQL = `
SELECT
  *
FROM
  notes."refTechnicianIntakeForm"
WHERE
  "refAppointmentId" = $1
  AND "refTITFQId" = 1
  `
