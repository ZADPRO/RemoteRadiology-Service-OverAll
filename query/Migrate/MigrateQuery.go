package query

var CheckMigrateDicomSQL = `
SELECT
  rdf.*,
  CASE
    WHEN EXISTS (
      SELECT
        1
      FROM
        migratefile."refMigrateDicomFiles" emd
      WHERE
        emd."refDFId" = rdf."refDFId"
    )
    OR rdf."refDFFilename" LIKE 'https://easeqt-health-archi%' THEN TRUE
    ELSE FALSE
  END AS "isMigrated"
FROM
  dicom."refDicomFiles" rdf
  JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rdf."refAppointmentId"
WHERE
  ra."refAppointmentStatus" = TRUE;
`

var CheckOneMigrateDicomSQL = `
SELECT
  rdf.*,
  CASE
    WHEN EXISTS (
      SELECT
        1
      FROM
        migratefile."refMigrateDicomFiles" emd
      WHERE
        emd."refDFId" = rdf."refDFId"
    )
    OR rdf."refDFFilename" LIKE 'https://easeqt-health-archi%' THEN TRUE
    ELSE FALSE
  END AS "isMigrated"
FROM
  dicom."refDicomFiles" rdf
  JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rdf."refAppointmentId"
WHERE
  ra."refAppointmentStatus" = TRUE
ORDER BY
  "isMigrated" ASC
LIMIT
  1
`

var NewMigrateDicomSQL = `
INSERT INTO
  migratefile."refMigrateDicomFiles" (
    "refDFId",
    "refUserId",
    "refAppointmentId",
    "refDFFilename",
    "refDFFilenames3",
	"refDFMigrationStatus"
  )
VALUES
  ($1, $2, $3, $4, $5, TRUE);
`

var UpdateDicomSQL = `
UPDATE
  dicom."refDicomFiles"
SET
  "refDFFilename" = $1
WHERE
  "refDFId" = $2
`
