package query

var CheckAccessSQL = `
SELECT
  CASE 
    WHEN "refAppointmentAccessStatus" = true 
         AND "refAppointmentAccessId" != ? THEN false
    ELSE true
  END AS "status",
  "refAppointmentAccessId"
FROM
  appointment."refAppointments"
WHERE
  "refAppointmentId" = ?
`

var GetOneUserAppointment = `
SELECT
  sc."refSCCustId",
  sc."refSCName",
  ra.*
FROM
  appointment."refAppointments" ra
  JOIN public."ScanCenter" sc ON sc."refSCId" = ra."refSCId"
WHERE
  ra."refUserId" = ?
  AND ra."refAppointmentId" = ?
`

var GetIntakeFormSQL = `
SELECT
  *
FROM
  notes."refIntakeForm"
WHERE
  "refAppointmentId" = ?
`

var GetTechnicianIntakeFormSQL = `
SELECT
  *
FROM
  notes."refTechnicianIntakeForm"
WHERE
  "refAppointmentId" = ?
`

var GetReportIntakeFormSQL = `
SELECT
  *
FROM
  notes."refReportIntakeForm"
WHERE
  "refAppointmentId" = ?
`

var GetReporttextContent = `
SELECT
  *
FROM
  notes."refReportsTextContent"
WHERE
  "refAppointmentId" = ?
`

var GetReportHistorySQL = `
SELECT
  f."refUserFirstName" AS "HandleUserName",
  rrh.*
FROM
  notes."refReportsHistory" rrh
  JOIN public."Users" f ON f."refUserId" = rrh."refRHHandledUserId"
WHERE
  rrh."refAppointmentId" = ?
  AND (rrh."refRHHandleStatus" != 'technologistformfill' OR rrh."refRHHandleStatus" IS NULL)
ORDER BY
  rrh."refRHId" ASC;
`

var GetReportCommentsSQL = `
SELECT
  f."refUserFirstName" AS "UserForName",
  b."refUserFirstName" AS "UserByName",
  rrc.*
FROM
  notes."refReportsComments" rrc
  JOIN public."Users" f ON f."refUserId" = rrc."refRCFor"
  JOIN public."Users" b ON b."refUserId" = rrc."refRCBy"
WHERE
  rrc."refAppointmentId" = ?
`

var GetAppointmentSQL = `
SELECT
  *
FROM
  appointment."refAppointments"
WHERE
  "refAppointmentId" = ?
`

var UpdateAccessAppointment = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentAccessStatus" = ?,
  "refAppointmentAccessId" = ?
WHERE
  "refAppointmentId" = ?
`

var GetReportIntakeFormQuestionSQL = `
SELECT
  *
FROM
  notes."refReportIntakeForm"
WHERE
  "refAppointmentId" = ?
  AND "refRITFQId" = ?
`

var GetTechnicianIntakeFormQuestionSQL = `
SELECT
  *
FROM
  notes."refTechnicianIntakeForm"
WHERE
  "refAppointmentId" = ?
  AND "refTITFQId" = ?
`

var GetPatientIntakeFormQuestionSQL = `
SELECT
  *
FROM
  notes."refIntakeForm"
WHERE
  "refAppointmentId" = ?
  AND "refITFQId" = ?
`

var InsertTechnicianIntakeSQL = `
INSERT INTO
  notes."refReportIntakeForm" (
    "refUserId",
    "refAppointmentId",
    "refRITFQId",
    "refRITFAnswer",
    "refRITFCreatedAt",
    "refRITFCreatedBy"
  )
VALUES
  (?, ?, ?, ?, ?, ?);
`

var UpdateReportIntakeSQL = `
UPDATE
  notes."refReportIntakeForm"
SET
  "refRITFAnswer" = ?,
  "refRITFUpdatedAt" = ?,
  "refRITFUpdatedBy" = ?
WHERE
  "refRITFQId" = ?
  AND "refAppointmentId" = ?
`

var UpdateTechnicianIntakeSQL = `
UPDATE
  notes."refTechnicianIntakeForm"
SET
  "refTITFAnswer" = ?,
  "refTITFUpdatedAt" = ?,
  "refTITFUpdatedBy" = ?
WHERE
  "refTITFQId" = ?
  AND "refAppointmentId" = ?
`

var UpdatePatientIntakeSQL = `
UPDATE
  notes."refIntakeForm"
SET
  "refITFAnswer" = ?,
  "refITFUpdatedAt" = ?,
  "refITFUpdatedBy" = ?
WHERE
  "refITFQId" = ?
  AND "refAppointmentId" = ?
`

var GetTextContentSQL = `
SELECT
  *
FROM
  notes."refReportsTextContent"
WHERE
  "refAppointmentId" = ?
`

var InsertTextContentSQL = `
INSERT INTO
  notes."refReportsTextContent" (
    "refUserId",
    "refAppointmentId",
    "refRTCText",
    "refRTCreatedAt",
    "refRTCreatedBy",
    "refRTSyncStatus"
  )
VALUES
  (?, ?, ?, ?, ?, ?);
`

var UpdateTextContentSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTCText" = ?,
  "refRTUpdatedAt" = ?,
  "refRTUpdatedBy" = ?,
  "refRTSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
`

var InsertCommentsSQL = `
INSERT INTO
  notes."refReportsComments" (
    "refUserId",
    "refAppointmentId",
    "refRCFor",
    "refRCBy",
    "refRCStatus",
    "refRCComments",
    "refRCCreatedAt"
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?);
`

var UpdateReportRemarksSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentRemarks" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var CheckLatestReportSQL = `
SELECT
  *
FROM
  notes."refReportsHistory"
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
ORDER BY
  "refRHId" DESC
`

var TechInsertReportHistorySQL = `
INSERT INTO
  notes."refReportsHistory" (
    "refUserId",
    "refAppointmentId",
    "refRHHandledUserId",
    "refRHHandleStartTime"
  )
VALUES
  (?, ?, ?, ?);
`

var InsertReportHistorySQL = `
INSERT INTO
  notes."refReportsHistory" (
    "refUserId",
    "refAppointmentId",
    "refRHHandledUserId",
    "refRHHandleStartTime"
  )
VALUES
  (?, ?, ?, ?);
`

var CompleteReportAppointmentSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentComplete" = ?,
  "refAppointmentAccessStatus" = false,
  "refAppointmentAccessId" = NULL,
  "refAppointmentImpression" = ?,
  "refAppointmentRecommendation" = ?,
  "refAppointmentImpressionAdditional" = ?,
  "refAppointmentRecommendationAdditional" = ?,
  "refAppointmentCommonImpressionRecommendation" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var CompleteReportHistorySQL = `
UPDATE notes."refReportsHistory"
SET
  "refRHHandleEndTime" = ?,
  "refRHHandleStatus" = ?,
  "refRHHandleContentText" = ?
WHERE "refRHId" = (
  SELECT MAX("refRHId")
  FROM notes."refReportsHistory"
  WHERE
    "refAppointmentId" = ?
    AND "refRHHandledUserId" = ?
    AND "refUserId" = ?
)
`

var InsertReportTemplate = `
INSERT INTO
  notes."refReportFormate" (
    "refRFName",
    "refRFText",
    "refRFCreatedAt",
    "refRFCreatedBy",
    "refRFStatus"
  )
VALUES
  (?, ?, ?, ?, true)
  RETURNING "refRFId"
`

var GetReportFormateListSQL = `
SELECT
  *
FROM
  notes."refReportFormate"
WHERE
  "refRFStatus" = true
`

var GetOneReportFormateListSQL = `
SELECT
  *
FROM
  notes."refReportFormate"
WHERE
  "refRFStatus" = true
  AND "refRFId" = ?
`

var GetUserDetailsSQL = `
SELECT
  u."refUserCustId",
  u."refUserId",
  u."refUserFirstName",
  rcd."refCODOEmail",
  u."refRTId",
  COALESCE(
    dd."refDDSpecialization",
    rad."refRASpecialization",
    cod."refCDSpecialization"
  ) AS "specialization",
  COALESCE(sc."refSCName", 'Wellthgreen HealthCare') AS "department"
FROM
  public."Users" u
  LEFT JOIN userdomain."refTechnicianDomain" td ON u."refRTId" = 2
  AND td."refUserId" = u."refUserId"
  LEFT JOIN userdomain."refReceptionstDomain" rd ON u."refRTId" = 3
  AND rd."refUserId" = u."refUserId"
  LEFT JOIN userdomain."refDoctorDomain" dd ON u."refRTId" = 5
  AND dd."refUserId" = u."refUserId"
  LEFT JOIN userdomain."refRadiologistDomain" rad ON u."refRTId" = 6
  AND rad."refUserId" = u."refUserId"
  LEFT JOIN userdomain."refScribeDomain" sd ON u."refRTId" = 7
  AND sd."refUserId" = u."refUserId"
  LEFT JOIN userdomain."refCoDoctorDomain" cod ON u."refRTId" = 8
  AND cod."refUserId" = u."refUserId"
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  LEFT JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
  LEFT JOIN public."ScanCenter" sc ON sc."refSCId" = rscm."refSCId"
WHERE
  u."refUserId" = $1;
`

var PatientUserDetailsSQL = `
SELECT
  *
FROM
  public."Users"
WHERE
  "refUserId" = ?
`

var ListUserDataSQL = `
SELECT
  *
FROM
  notes."refReportsHistory" rrh
  JOIN public."Users" u ON u."refUserId" = rrh."refRHHandledUserId"
WHERE
  rrh."refUserId" = ?
  AND rrh."refAppointmentId" = ?
  AND u."refRTId" = ?
ORDER BY
  rrh."refRHId" DESC
`

var UpdateCorrectEditSQL = `
UPDATE
  notes."refReportsHistory"
SET
  "refRHHandleCorrect" = ?,
  "refRHHandleEdit" = ?
WHERE
  "refRHId" = ?
`

var GetPatientData = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rc ON rc."refUserId" = u."refUserId"
  JOIN appointment."refAppointments" ra ON ra."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
  AND ra."refAppointmentId" = ?
`

var GetManagerData = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rc ON rc."refUserId" = u."refUserId"
  RIGHT JOIN map."refScanCenterMap" rscmp ON rscmp."refUserId" = u."refUserId"
  RIGHT JOIN appointment."refAppointments" ra ON ra."refSCId" = rscmp."refSCId"
  JOIN public."ScanCenter" sc ON sc."refSCId" = ra."refSCId"
WHERE
  u."refRTId" = ?
  AND ra."refAppointmentId" = ?
`
