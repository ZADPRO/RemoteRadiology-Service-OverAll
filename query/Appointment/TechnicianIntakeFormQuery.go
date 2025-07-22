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
      ? AS reftimeZone,
      jsonb_array_elements(?::jsonb) AS item
  )
INSERT INTO
  notes."refTechnicianIntakeForm" (
    "refUserId",
    "refAppointmentId",
    "refTITFQId",
    "refTITFAnswer",
    "refTITFCreatedAt",
    "refTITFCreatedBy"
  )
SELECT
  refUserId,
  refAppointmentId,
  (item ->> 'questionId')::int,
  item ->> 'answer',
  reftimeZone,
  refCreatedAt
FROM
  input_data;
`

var GetDicomFileSQL = `
SELECT
  *
FROM
  dicom."refDicomFiles"
WHERE
  "refDFId" = ?
`

var GetTechIntakeForm = `
SELECT
  *
FROM
  notes."refTechnicianIntakeForm"
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var ViewGetDicomFile = `
SELECT
  *
FROM
  dicom."refDicomFiles"
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
`

var GetDicomFile = `
SELECT
  *
FROM
  dicom."refDicomFiles"
WHERE
  "refAppointmentId" = ?
  AND "refUserId" = ?
  AND "refDFSide" = ?
`

var GetAllDicomSQL = `
SELECT
  *
FROM
  dicom."refDicomFiles"
WHERE
  "refAppointmentId" IN ?
ORDER BY
  "refAppointmentId" ASC
`

var ListTechnicianSQL = `
SELECT
  *
FROM
  notes."refReportsHistory"
WHERE
  "refUserId" = ?
  AND "refAppointmentId" = ?
  AND "refRHHandledUserId" = ?
`
