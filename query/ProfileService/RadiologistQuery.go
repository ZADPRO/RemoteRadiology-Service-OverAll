package query

var GetListofRadiologistSQL = `
SELECT
  *
FROM
  "Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refRadiologistDomain" rdd ON rdd."refUserId" = u."refUserId"
  JOIN "RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
  u."refRTId" = ?
ORDER BY
  u."refUserId" DESC  
`

var GetListofRadiologistOneSQL = `
SELECT
  *
FROM
  "Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refRadiologistDomain" rtd ON rtd."refUserId" = u."refUserId"
  JOIN "RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
 u."refRTId" = ?
 AND u."refUserId" = ?
`

var GetCVFilesSQL = `
SELECT
  *
FROM
  userdomain."refCV"
WHERE
  "refUserId" = ?
  AND "refCVStatus" = true
`

var GetECFilesSQL = `
SELECT
  *
FROM
  userdomain."refEducationCertificate"
WHERE
  "refUserId" = ?
  AND "refECStatus" = true
`

var GetLicenseFilesSQL = `
SELECT
  *
FROM
  userdomain."refLicences"
WHERE
  "refUserId" = ?
  AND "refLStatus" = true
`

var GetMalpracticeFilesSQL = `
SELECT
  *
FROM
userdomain."refMalpractice"
WHERE
  "refUserId" = ?
  AND "refMPStatus" = true
`
