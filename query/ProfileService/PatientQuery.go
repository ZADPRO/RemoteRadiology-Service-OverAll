package query

var GetAllPatientList = `
SELECT
  *
FROM
  map."refScanCenterMapPatient" rscmp
  JOIN public."Users" u ON u."refUserId" = rscmp."refUserId"
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
WHERE
  rscmp."refSCId" = $1
ORDER BY u."refUserId" DESC
`

// var GetOnePatientList = `
// SELECT
//   u.*,
//   rcd.*,
//   COALESCE(ra.appointments, '[]') AS appointments
// FROM
//   public."Users" u
//   JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
//   LEFT JOIN LATERAL (
//     SELECT
//       json_agg(
//         json_build_object(
//           'refAppointmentId',
//           ra."refAppointmentId",
//           'refUserId',
//           ra."refUserId",
//           'refAppointmentDate',
//           ra."refAppointmentDate",
//           'refAppointmentStatus',
//           ra."refAppointmentStatus",
//           'refSCId',
//           sc."refSCId",
//           'refSCustId',
//           sc."refSCCustId"
//         )
//       ) AS appointments
//     FROM
//       appointment."refAppointments" ra
//       JOIN public."ScanCenter" sc ON sc."refSCId" = ra."refSCId"
//     WHERE
//       ra."refUserId" = u."refUserId"
//       AND ra."refAppointmentStatus" = true
//   ) ra ON true
// WHERE
//   u."refUserId" = $1;
// `

var GetOnePatientList = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN map."refScanCenterMapPatient" rscmp ON rscmp."refUserId" = u."refUserId"
  JOIN public."ScanCenter" sc ON sc."refSCId" = rscmp."refSCId"
WHERE
  u."refUserId" = $1
`

var GetAppointmentListSQL = `
SELECT
  CASE
    WHEN (
      ra."refAppointmentComplete" = 'fillform'
      OR ra."refAppointmentComplete" = 'technologistformfill'
    )
    AND (
      SELECT
        COUNT(*)
      FROM
        dicom."refDicomFiles" rdf
      WHERE
        rdf."refAppointmentId" = ra."refAppointmentId"
    ) = 0 THEN true
    ELSE false
  END AS "allowCancelResh",
  ra.*,
  sc.*
FROM
  appointment."refAppointments" ra
  JOIN public."ScanCenter" sc ON sc."refSCId" = ra."refSCId"
WHERE
  ra."refUserId" = $1
  AND ra."refAppointmentStatus" = true
ORDER BY
  ra."refAppointmentId" ASC;
`
