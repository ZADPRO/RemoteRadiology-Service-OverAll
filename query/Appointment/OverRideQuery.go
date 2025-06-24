package query

var ListAllOverRideSQL = `
SELECT
  *
FROM
  appointment."refAppointments" ra
  JOIN notes."refOverRide" ror ON ror."refAppointmentId" = ra."refAppointmentId"
WHERE
  ra."refAppointmentId" = ?
`

var UpdateOverRideSQL = `
UPDATE
  notes."refOverRide"
SET
  "refApprovedStatus" = ?
WHERE
  "refOVId" = ?
`
