package query

var VerifyAppointment = `
SELECT
  COUNT(*) AS "TotalCount"
FROM
  appointment."refAppointments"
WHERE
  "refSCId" = ?
  AND "refAppointmentDate" = ?
  AND (
    (
      "refAppointmentStartTime" >= ?
      AND "refAppointmentEndTime" <= ?
    )
  );
`
