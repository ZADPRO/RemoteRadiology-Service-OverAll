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
  "refAppointmentComplete" = ?,
  "refAppointmentPriority" = ?
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
  public."Users"
WHERE
  "refUserId" = $1
`

var GetTextContent = `
SELECT
  *
FROM
  notes."refReportsTextContent"
WHERE
  AND "refAppointmentId" = ANY ($1)
`

var GetAppointmentConsent = `
SELECT
  *
FROM
  appointment."refAppointments"
WHERE
  "refAppointmentId" = ANY ($1)
`
