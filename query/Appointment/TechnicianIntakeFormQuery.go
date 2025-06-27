package query

var GetCategoryId = `
SELECT
  *
FROM
  appointment."refAppointments"
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var UpdateCategoryId = `
UPDATE
  appointment."refAppointments"
SET
  "refCategoryId" = ?
WHERE
  "refAppointmentId" = ?
`

var TechnicianInsertAnswerSQL = `
WITH
  input_data AS (
    SELECT
      ?::int AS refUserId,
      ?::int AS refAppointmentId,
      ?::int AS refCreatedAt,
      jsonb_array_elements(?::jsonb) AS item
  )
INSERT INTO
  notes."refIntakeForm" (
    "refUserId",
    "refAppointmentId",
    "refITFQId",
    "refITFAnswer",
    "refITFCreatedAt",
    "refITFCreatedBy"
  )
SELECT
  refUserId,
  refAppointmentId,
  (item ->> 'questionId')::int,
  item ->> 'answer',
  NOW(),
  refCreatedAt
FROM
  input_data;
`
