package query

var ListPatientConsentSQL = `
SELECT
  "refAppointmentId",
  "refAppointmentConsent"
FROM
  appointment."refAppointments"
WHERE
  "refAppointmentId" = ANY ($1)
`

var ListWGPatientBrochureSQL = `
SELECT
  *
FROM
  forms."formsData"
WHERE
  "refFTId" = $1
`

var ListSCPatientBrochureSQL = `
SELECT
  *
FROM
  forms."formsData"
WHERE
  "refFTId" = $1
  AND "refSCId" = $2
`

var InsertWGPatientBrochureSQL = `
INSERT INTO
  forms."formsData" ("refFTId", "refFDData")
VALUES
  ($1, $2)
`

var InsertscPatientBrochureSQL = `
INSERT INTO
  forms."formsData" (
    "refFTId",
    "refSCId",
    "refFDData",
    "refFDAccessData"
  )
VALUES
  ($1, $2, $3, $4);
`

var UpdateWGPatientBrochureSQL = `
UPDATE
  forms."formsData"
SET
  "refFDData" = $1
WHERE
  "refFTId" = $2;
`

var UpdatescPatientBrochureSQL = `
UPDATE
  forms."formsData"
SET
  "refFDData" = $1,
  "refFDAccessData" = $2
WHERE
  "refFTId" = $3
  AND "refSCId" = $4
`
