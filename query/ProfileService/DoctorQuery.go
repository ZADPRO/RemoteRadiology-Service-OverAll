package query

var GetListofDoctorSQL = `
SELECT
  *
FROM
  "Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refDoctorDomain" rdd ON rdd."refUserId" = u."refUserId"
  JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
  JOIN "RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
  u."refRTId" = ?
  AND rscm."refSCId" = ?
  ORDER BY
  u."refUserId" DESC
`

var GetListofDoctorOneSQL = `
SELECT
  *
FROM
  "Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refDoctorDomain" rdd ON rdd."refUserId" = u."refUserId"
  JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
  JOIN "RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
  u."refUserId" = ?
  AND rscm."refSCId" = ?
`

var GetMedicalLicenseSecuritySQL = `
SELECT
  *
FROM
  userdomain."refMedicalLicenseSecurity"
WHERE
  "refUserId" = ?
`
