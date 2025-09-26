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

var GetReportHistoryFullSQL = `
SELECT
  f."refUserFirstName" AS "HandleUserName",
  f."refRTId" AS "HandlerRTId",
  rrh.*
FROM
  notes."refReportsHistory" rrh
  JOIN public."Users" f ON f."refUserId" = rrh."refRHHandledUserId"
WHERE
  rrh."refAppointmentId" = ?
ORDER BY
  rrh."refRHId" ASC;
`

var GetReportHistorySQL = `
SELECT
  f."refUserFirstName" AS "HandleUserName",
  f."refRTId" AS "HandlerRTId",
  rrh.*
FROM
  notes."refReportsHistory" rrh
  JOIN public."Users" f ON f."refUserId" = rrh."refRHHandledUserId"
WHERE
  rrh."refAppointmentId" = ?
  AND f."refRTId" IN (1, 2, 3, 5, 8, 10)
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
  notes."refReportsHistory" rrh
  JOIN public."Users" u ON u."refUserId" = rrh."refUserId"
WHERE
  rrh."refAppointmentId" = ?
  AND rrh."refUserId" = ?
ORDER BY
  rrh."refRHId" DESC
`

var TechInsertReportHistorySQL = `
INSERT INTO
  notes."refReportsHistory" (
    "refUserId",
    "refAppointmentId",
    "refRHHandledUserId",
    "refRHHandleStartTime",
    "refRHHandleStatus"
  )
VALUES
  (?, ?, ?, ?, ?);
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

var AddedumCountSQL = `
SELECT
  count(*) AS count
FROM
  notes."refAddendum"
WHERE
  "refAppointmentId" = $1
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
    "refRFStatus",
    "refRFAccessStatus"
  )
VALUES
  ($1, $2, $3, $4, true, $5)
RETURNING
  "refRFId",
  (
    SELECT
      "refUserCustId"
    FROM
      public."Users"
    WHERE
      "refUserId" = $4
  );
`

var GetReportFormateAllListSQL = `
SELECT
  u."refUserCustId",
  rrf.*
FROM
  notes."refReportFormate" rrf
  JOIN public."Users" u ON u."refUserId" = rrf."refRFCreatedBy"
WHERE
  rrf."refRFStatus" = true
  ORDER BY
  rrf."refRFId" DESC
`

var GetReportFormateListSQL = `
SELECT
  *
FROM
  notes."refReportFormate" rrf
  JOIN public."Users" u ON u."refUserId" = rrf."refRFCreatedBy"
WHERE
  (
    rrf."refRFCreatedBy" = $1
    OR (
      rrf."refRFCreatedBy" <> $1
      AND rrf."refRFAccessStatus" = 'public'
    )
  )
  AND (rrf."refRFStatus" = true)
`

var DeleteReportTemplateSQL = `
UPDATE notes."refReportFormate"
SET 
  "refRFStatus" = false
WHERE "refRFId" = $1;
`
var UpdateReportTemplateSQL = `
UPDATE notes."refReportFormate"
SET 
  "refRFAccessStatus" = $1
WHERE "refRFId" = $2;
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
  AND u."refRTId" IN ?
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

var UpdateHandlerCorrectEditIdSQL = `
UPDATE notes."refReportsHistory"
SET
  "refRHHandleCorrect" = $1,
  "refRHHandleEdit" = $2
WHERE "refRHId" = (
  SELECT "refRHId"
  FROM notes."refReportsHistory"
  WHERE "refRHHandledUserId" = $3
    AND "refAppointmentId" = $4
  ORDER BY "refRHId" DESC
  LIMIT 1
);
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
  "refDDEaseQTReportAccess",
  "refDDNAsystemReportAccess"
FROM
  userdomain."refDoctorDomain"
WHERE
  "refUserId" = $1
`

var CoDoctorReportAccessSQL = `
SELECT
  "refCDEaseQTReportAccess",
  "refCDNAsystemReportAccess"
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

var UpdateAutosavePatientHistorySyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTPatientHistorySyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveBreastImplantSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTBreastImplantSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveSymmetrySyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTSymmetrySyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveBreastDensitySyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTBreastDensityandImageRightSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveNippleAreolaSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTNippleAreolaSkinRightSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveGlandularSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTGrandularAndDuctalTissueRightSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLymphNodesSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTLymphNodesRightSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLesionsSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTLesionsRightSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveComparisonPriorSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTComparisonPriorSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveBreastDensityLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTBreastDensityandImageLeftSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveNippleAreolaLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTNippleAreolaSkinLeftSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveGlandularLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTGrandularAndDuctalTissueLeftSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLymphNodesLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTLymphNodesLeftSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLesionsLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTLesionsLeftSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveComparisonPriorLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTComparisonPriorLeftSyncStatus" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveBreastImplantReportTextSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTBreastImplantReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveSymmetryReportTextSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTSymmetryReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveBreastDensityReportTextSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTBreastDensityandImageRightReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveNippleAreolaReportTextSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTNippleAreolaSkinRightReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLesionsReportTextTextSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTLesionsRightReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveComparisonPriorReportTextSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTComparisonPriorReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveGrandularAndDuctalTissueReportTextSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTGrandularAndDuctalTissueRightReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLymphNodesReportTextSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTLymphNodesRightReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveBreastDensityReportTextLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTBreastDensityandImageLeftReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveNippleAreolaReportTextLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTNippleAreolaSkinLeftReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLesionsReportTextLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTLesionsLeftReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveComparisonPriorReportTextLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTComparisonPriorLeftReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveGrandularAndDuctalTissueReportTextLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTGrandularAndDuctalTissueLeftReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateAutosaveLymphNodesReportTextLeftSyncSQL = `
UPDATE
  notes."refReportsTextContent"
SET
  "refRTLymphNodesLeftReportText" = ?
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateReportAppointmentSQL = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentComplete" = $1
WHERE
  "refAppointmentId" = $2
`

var ListOldReportSQL = `
SELECT
  *
FROM
  notes."refOldReport" ror
WHERE
  ror."refAppointmentId" = $1
  AND ror."refUserId" = $2
  AND ror."refORCategoryId" = $3
  AND ror."refORStatus" = true
`

var AddOldReportSQL = `
INSERT INTO
  notes."refOldReport" (
    "refUserId",
    "refAppointmentId",
    "refORCategoryId",
    "refORFilename",
    "refORCreatedAt",
    "refORCreatedBy",
    "refORStatus"
  )
VALUES
  ($1, $2, $3, $4, $5, $6, true);
`

var GetParticularOldReport = `
SELECT
  *
FROM
  notes."refOldReport" ror
WHERE
  ror."refORId" = $1
`

var DeleteOldReportSQL = `
UPDATE
  notes."refOldReport"
SET
  "refORStatus" = $1
WHERE
  "refORId" = $2;
`

var GetOldReportsSQL = `
SELECT 
  r."refORCategoryId",
  COALESCE(
    json_agg(r."refORFilename") FILTER (WHERE r."refORFilename" IS NOT NULL), 
    '[]'
  ) AS files
FROM notes."refOldReport" r
WHERE r."refUserId" = $1
  AND r."refAppointmentId" = $2
  AND r."refORStatus" = true
GROUP BY r."refORCategoryId"
ORDER BY r."refORCategoryId";
`

var GetPatientPrivatePublicSQL = `
SELECT
  *
FROM
  notes."refTechnicianIntakeForm"
WHERE
  "refAppointmentId" = $1
  AND "refTITFQId" = $2;
`