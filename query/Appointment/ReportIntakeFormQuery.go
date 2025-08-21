package query

var CheckAccessSQL = `
SELECT
  CASE
    WHEN "refAppointmentAccessStatus" = true
    AND "refAppointmentAccessId" != ? THEN false
    ELSE true
  END AS "status",
  "refAppointmentAccessId",
  (
    SELECT
      "refUserCustId"
    FROM
      public."Users"
    WHERE
      "refUserId" = "refAppointmentAccessId"
    ) AS "userCustId"
FROM
  appointment."refAppointments"
WHERE
  "refAppointmentId" = ?
`

var ScribeCheckAccessSQL = `
SELECT
  CASE
    WHEN "refAppointmentScribeAccessStatus" = true
    AND "refAppointmentScribeAccessId" != ? THEN false
    ELSE true
  END AS "status",
  "refAppointmentScribeAccessId" AS "refAppointmentAccessId",
  (
    SELECT
      "refUserCustId"
    FROM
      public."Users"
    WHERE
      "refUserId" = "refAppointmentScribeAccessId"
    ) AS "userCustId"
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

var ScribeUpdateAccessAppointment = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentScribeAccessStatus" = ?,
  "refAppointmentScribeAccessId" = ?
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

var AutoCompleteReportAppointmentSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentImpression" = ?,
  "refAppointmentRecommendation" = ?,
  "refAppointmentImpressionAdditional" = ?,
  "refAppointmentRecommendationAdditional" = ?,
  "refAppointmentCommonImpressionRecommendation" = ?,
  "refAppointmentImpressionRight" = ?,
  "refAppointmentRecommendationRight" = ?,
  "refAppointmentImpressionAdditionalRight" = ?,
  "refAppointmentRecommendationAdditionalRight" = ?,
  "refAppointmentCommonImpressionRecommendationRight" = ?,
  "refAppointmentReportArtifactsLeft" = ?,
  "refAppointmentReportArtifactsRight" = ?,
  "refAppointmentPatietHistory" = ?,
  "refAppointmentBreastImplantImageText" = ?,
  "refAppointmentSymmetryImageText" = ?,
  "refAppointmentBreastdensityImageText" = ?,
  "refAppointmentNippleAreolaImageText" = ?,
  "refAppointmentGlandularImageText" = ?,
  "refAppointmentLymphnodeImageText" = ?,
  "refAppointmentBreastdensityImageTextLeft" = ?,
  "refAppointmentNippleAreolaImageTextLeft" = ?,
  "refAppointmentGlandularImageTextLeft" = ?,
  "refAppointmentLymphnodeImageTextLeft" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
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
  "refAppointmentCommonImpressionRecommendation" = ?,
  "refAppointmentImpressionRight" = ?,
  "refAppointmentRecommendationRight" = ?,
  "refAppointmentImpressionAdditionalRight" = ?,
  "refAppointmentRecommendationAdditionalRight" = ?,
  "refAppointmentCommonImpressionRecommendationRight" = ?,
  "refAppointmentReportArtifactsLeft" = ?,
  "refAppointmentReportArtifactsRight" = ?,
  "refAppointmentPatietHistory" = ?,
  "refAppointmentBreastImplantImageText" = ?,
  "refAppointmentSymmetryImageText" = ?,
  "refAppointmentBreastdensityImageText" = ?,
  "refAppointmentNippleAreolaImageText" = ?,
  "refAppointmentGlandularImageText" = ?,
  "refAppointmentLymphnodeImageText" = ?,
  "refAppointmentBreastdensityImageTextLeft" = ?,
  "refAppointmentNippleAreolaImageTextLeft" = ?,
  "refAppointmentGlandularImageTextLeft" = ?,
  "refAppointmentLymphnodeImageTextLeft" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var ScribeCompleteReportAppointmentSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentComplete" = ?,
  "refAppointmentScribeAccessStatus" = false,
  "refAppointmentScribeAccessId" = NULL,
  "refAppointmentImpression" = ?,
  "refAppointmentRecommendation" = ?,
  "refAppointmentImpressionAdditional" = ?,
  "refAppointmentRecommendationAdditional" = ?,
  "refAppointmentCommonImpressionRecommendation" = ?,
  "refAppointmentImpressionRight" = ?,
  "refAppointmentRecommendationRight" = ?,
  "refAppointmentImpressionAdditionalRight" = ?,
  "refAppointmentRecommendationAdditionalRight" = ?,
  "refAppointmentCommonImpressionRecommendationRight" = ?,
  "refAppointmentReportArtifactsLeft" = ?,
  "refAppointmentReportArtifactsRight" = ?,
  "refAppointmentPatietHistory" = ?,
  "refAppointmentBreastImplantImageText" = ?,
  "refAppointmentSymmetryImageText" = ?,
  "refAppointmentBreastdensityImageText" = ?,
  "refAppointmentNippleAreolaImageText" = ?,
  "refAppointmentGlandularImageText" = ?,
  "refAppointmentLymphnodeImageText" = ?,
  "refAppointmentBreastdensityImageTextLeft" = ?,
  "refAppointmentNippleAreolaImageTextLeft" = ?,
  "refAppointmentGlandularImageTextLeft" = ?,
  "refAppointmentLymphnodeImageTextLeft" = ?
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

var InsertRemark = `
INSERT INTO notes."Remarks" (
  "refAppointmentId",
  "refUserId",
  "refRemarksMessage",
  "refRCreatedAt"
) VALUES (
  $1,
  $2,
  $3,
  $4
);
`
var ListRemarkSQL = `
SELECT
  r.*,
  u."refUserCustId"
FROM
  notes."Remarks" r
  JOIN public."Users" u ON u."refUserId" = r."refUserId"
WHERE
  "refAppointmentId" = $1
`

var UpdateMailStatusSQL = `UPDATE
  appointment."refAppointments"
SET
  "refAppointmentMailSendStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var DoctorReportAccessSQL = `
SELECT
  "refDDEaseQTReportAccess"
FROM
  userdomain."refDoctorDomain"
WHERE
  "refUserId" = $1
`

var CoDoctorReportAccessSQL = `
SELECT
  "refCDEaseQTReportAccess"
FROM
  userdomain."refCoDoctorDomain"
WHERE
  "refUserId" = $1
`

var ScanCenterSQL = `
SELECT
  *
FROM
  public."ScanCenter"
WHERE
  "refSCId" = $1
`

var DownloadReportSQL = `
SELECT
  *
FROM
  notes."refIntakeForm"
WHERE
  "refITFId" = $1
`

var InsertAddedumSQL = `
INSERT INTO notes."refAddendum"(
	"refAppointmentId", "refUserId", "refADText", "refADCreatedAt")
	VALUES ($1, $2, $3, $4);
`

var ListAddendumSQL = `
SELECT
  ra.*,
  u."refUserCustId"
FROM
  notes."refAddendum" ra
  JOIN public."Users" u ON u."refUserId" = ra."refUserId"
WHERE
  ra."refAppointmentId" = $1
`

var UpdateAutosaveTextContentSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTCText" = ?,
  "refRTUpdatedAt" = ?,
  "refRTUpdatedBy" = ?
WHERE
  "refAppointmentId" = ?
`

var UpdateAutosaveSyncStatusSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTSyncStatus" = ?,
  "refRTUpdatedAt" = ?,
  "refRTUpdatedBy" = ?
WHERE
  "refAppointmentId" = ?
`

var UpdateAutosaveImpressionSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentImpression" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveRecommendationSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentRecommendation" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveImpressionAddtionalSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentImpressionAdditional" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveRecommendationAddtionalSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentRecommendationAdditional" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveCommonImpressionRecommendationSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentCommonImpressionRecommendation" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveImpressionRightSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentImpressionRight" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveRecommendationRightSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentRecommendationRight" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveImpressionAddtionalRightSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentImpressionAdditionalRight" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveRecommendationAddtionalRightSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentRecommendationAdditionalRight" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveCommonImpressionRecommendationRightRightSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentCommonImpressionRecommendationRight" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`
var UpdateAutosaveArtificatsLeftSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentReportArtifactsLeft" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveArtificatsRightSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentReportArtifactsRight" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosavePatientHistorySQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentPatietHistory" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveBreastImplantsImagetextSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentBreastImplantImageText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveSymmetryImageTextSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentSymmetryImageText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveBreastDensityImageTextSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentBreastdensityImageText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`
var UpdateAutosaveNippleAreolaImageTextSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentNippleAreolaImageText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveGlandularImageTextSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentGlandularImageText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLymphnodesImageTextSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentLymphnodeImageText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveBreastDensityImageTextLeftSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentBreastdensityImageTextLeft" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveNippleAreolaImageTextLeftSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentNippleAreolaImageTextLeft" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveGlandularImageTextLeftSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentGlandularImageTextLeft" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLymphNodesImageTextLeftSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentLymphnodeImageTextLeft" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`
