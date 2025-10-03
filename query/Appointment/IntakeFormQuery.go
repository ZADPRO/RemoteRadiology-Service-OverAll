package query

var InsertAnswerSQL = `
WITH
  input_data AS (
    SELECT
      ?::int AS refUserId,
      ?::int AS refAppointmentId,
      ?::int AS refCreatedAt,
      ? AS refTimezone,
      jsonb_array_elements(?::jsonb) AS item
  )
INSERT INTO
  notes."refIntakeForm" (
    "refUserId",
    "refAppointmentId",
    "refITFQId",
    "refITFAnswer",
    "refITFCreatedAt",
    "refITFCreatedBy",
    "refITFUpdatedBy"
  )
SELECT
  refUserId,
  refAppointmentId,
  (item ->> 'questionId')::int,
  item ->> 'answer',
  refTimezone,
  refCreatedAt,
  refCreatedAt
FROM
  input_data;
`

var ViewIntakeFormQuery = `
SELECT
  *
FROM
  notes."refIntakeForm"
WHERE
  "refUserId" = ?
  AND "refAppointmentId" = ?
`

var GetVerifyIntakeFormQuery = `
SELECT
  *
FROM
  notes."refOverRide"
WHERE
  "refAppointmentId" = ?
`

var UpdateAppointment = `
UPDATE
  appointment."refAppointments"
SET
  "refCategoryId" = ?,
  "refAppointmentConsent" = ?
WHERE
  "refAppointmentId" = ?
`

var UpdateAppointmentStatus = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentComplete" = ?
WHERE
  "refAppointmentId" = ?
`

var UpdateTechnicianAppointmentStatus = `
UPDATE
  appointment."refAppointments"
SET
  "refAppointmentDate" = ?,
  "refAppointmentComplete" = ?,
  "refAppointmentPriority" = ?,
  "refAppointmentTechArtifactsLeft" = ?,
  "refAppointmentTechArtifactsRight" = ?
WHERE
  "refAppointmentId" = ?
`

var GetIntakeDataSQL = `
SELECT
  *
FROM
  notes."refIntakeForm"
WHERE
  "refITFId" = ?
`

var GetIntakeAppointmentDataSQL = `
SELECT
  *
FROM
  notes."refIntakeForm"
WHERE
  "refAppointmentId" = $1
  AND "refITFQId" = $2
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

var InsertReportIntakeAllSQL = `
INSERT INTO notes."refReportIntakeForm"
    ("refUserId", "refAppointmentId", "refRITFQId", "refRITFCreatedAt", "refRITFCreatedBy")
SELECT
    $1 AS refUserId,          -- constant
    $2 AS refAppointmentId, -- constant
    gs AS refRITFQId,        -- series 1 to 137
    $3 AS refRITFCreatedAt, -- constant
    $4 AS refRITFCreatedBy   -- constant
FROM generate_series(1, 137) gs;
`

var InsertNewReportTextContentSQL = `
INSERT INTO notes."refReportsTextContent" (
    "refUserId",
    "refAppointmentId",
    "refRTCreatedAt",
    "refRTCreatedBy",
    "refRTSyncStatus",
    "refRTPatientHistorySyncStatus",
    "refRTBreastImplantSyncStatus",
    "refRTSymmetrySyncStatus",
    "refRTBreastDensityandImageRightSyncStatus",
    "refRTNippleAreolaSkinRightSyncStatus",
    "refRTLesionsRightSyncStatus",
    "refRTComparisonPriorSyncStatus",
    "refRTGrandularAndDuctalTissueRightSyncStatus",
    "refRTLymphNodesRightSyncStatus",
    "refRTBreastDensityandImageLeftSyncStatus",
    "refRTNippleAreolaSkinLeftSyncStatus",
    "refRTLesionsLeftSyncStatus",
    "refRTComparisonPriorLeftSyncStatus",
    "refRTGrandularAndDuctalTissueLeftSyncStatus",
    "refRTLymphNodesLeftSyncStatus"
) VALUES (
    $1,
    $2,
    $3,
    $4,
    TRUE, -- refRTSyncStatus
    TRUE, -- refRTPatientHistorySyncStatus
    TRUE, -- refRTBreastImplantSyncStatus
    TRUE, -- refRTSymmetrySyncStatus
    TRUE, -- refRTBreastDensityandImageRightSyncStatus
    TRUE, -- refRTNippleAreolaSkinRightSyncStatus
    TRUE, -- refRTLesionsRightSyncStatus
    TRUE, -- refRTComparisonPriorSyncStatus
    TRUE, -- refRTGrandularAndDuctalTissueRightSyncStatus
    TRUE, -- refRTLymphNodesRightSyncStatus
    TRUE, -- refRTBreastDensityandImageLeftSyncStatus
    TRUE, -- refRTNippleAreolaSkinLeftSyncStatus
    TRUE, -- refRTLesionsLeftSyncStatus
    TRUE, -- refRTComparisonPriorLeftSyncStatus
    TRUE, -- refRTGrandularAndDuctalTissueLeftSyncStatus
    TRUE -- refRTLymphNodesLeftSyncStatus
)
`

var UpdateCreateIntakeDataSQL = `
UPDATE
  notes."refIntakeForm"
SET
  "refITFAnswer" = $1,
  "refITFUpdatedBy" = $2,                
  "refITFUpdatedAt" = $3        
WHERE
"refAppointmentId" = $4
  AND "refITFQId" = $5                
`

var InsertIntakeSQL = `
INSERT INTO
  notes."refIntakeForm" (
    "refUserId",
    "refAppointmentId",
    "refITFQId",
    "refITFAnswer",
    "refITFCreatedAt",
    "refITFCreatedBy"
  )
VALUES
  ($1, $2, $3, $4, $5, $6)
`

var UpdateIntakeDataSQL = `
UPDATE
  notes."refIntakeForm"
SET
  "refITFAnswer" = ?,
  "refITFUpdatedBy" = ?,                
  "refITFUpdatedAt" = ?,           
  "refITFVerifiedTechnician" = ?      
WHERE
  "refITFId" = ?                     
`

var GetAuditforIntakeForm = `
SELECT
  *
FROM
  aduit."refTransHistory"
WHERE
  "transTypeId" IN (23, 24);
`

var TechnicianUserSQL = `
SELECT
  *
FROM
  notes."refReportsHistory" rrh
  JOIN public."Users" u ON u."refUserId" = rrh."refRHHandledUserId"
WHERE
  rrh."refRHHandleStatus" = 'Technologist Form Fill'
  AND rrh."refAppointmentId" = $1;
`

var GetTextContent = `
SELECT
  *
FROM
  notes."refReportsTextContent"
WHERE
  "refAppointmentId" = ANY ($1)
`

var GetAppointmentConsent = `
SELECT
  *
FROM
  appointment."refAppointments"
WHERE
  "refAppointmentId" = ANY ($1)
`

var UpdateOverrideSQL = `
UPDATE
  notes."refOverRide"
SET
  "refApprovedStatus" = $1,
  "refApprovedBy" = $2,
  "refApprovedAt" = $3
WHERE
  "refAppointmentId" = $4
`

var CheckOverride = `
SELECT
  *
FROM
  notes."refOverRide"
WHERE
  "refAppointmentId" = $1
`

var ReportHistorySQL = `
INSERT INTO notes."refReportsHistory" (
  "refUserId",
  "refAppointmentId",
  "refRHHandledUserId",
  "refRHHandleStartTime",
  "refRHHandleEndTime",
  "refRHHandleStatus"
) VALUES (
  $1, $2, $3, $4, $5, $6
);
`
