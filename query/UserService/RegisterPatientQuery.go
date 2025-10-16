package query

var CheckpatientExits = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
WHERE
  rcd."refCODOEmail" = $1
  OR rcd."refCODOPhoneNo1" = $2
  OR lower(u."refUserCustId") = lower($3)
`

var GetAllPatientDataQuery = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`

var UpdatePatientQuery = `
UPDATE
  public."Users"
SET
  "refUserCustId" = ?,
  "refUserFirstName" = ?,
  "refUserProfileImg" = ?,
  "refUserDOB" = ?,
  "refUserGender" = ?,
  "refUserStatus" = ?
WHERE
  "refUserId" = ?
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
