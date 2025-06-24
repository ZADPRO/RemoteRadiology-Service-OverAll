package query

var GetOneTechnicianSQL = `
SELECT
  *
FROM
  "Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refTechnicianDomain" rtd ON rtd."refUserId" = u."refUserId"
  JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
  JOIN "RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
  u."refUserId" = ?
  AND rscm."refSCId" = ?
`


var GetListofTechnicianSQL = `
SELECT
  *
FROM
  "Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refTechnicianDomain" rtd ON rtd."refUserId" = u."refUserId"
  JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
  JOIN "RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
  u."refRTId" = ?
  AND rscm."refSCId" = ?
ORDER BY
  u."refUserId" DESC
`

var GetListofTechnicianOneSQL = `
SELECT
  *
FROM
  "Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refTechnicianDomain" rtd ON rtd."refUserId" = u."refUserId"
  JOIN "RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
 u."refRTId" = 2
 AND u."refUserId" = ?
`

var GetStaffExperienceSQL = `
SELECT
  *
FROM
  userdomain."refStaffExprience"
WHERE
  "refUserId" = ?
`
var GetTechnicianMapListSQL = `
SELECT
  rscm.*,
  sc."refSCName",
  rt."refRTName"
FROM
  map."refScanCenterMap" rscm
  JOIN "ScanCenter" sc ON sc."refSCId" = rscm."refSCId"
  JOIN "RoleType" rt ON rt."refRTId" = rscm."refRTId"
WHERE
  rscm."refUserId" = ?
`

var GetTechnicianAuditSQL = `
SELECT
  rth.*,
  tt."transTypeName",
  u."refUserFirstName" || ' ' || u."refUserLastName" AS "Username"
FROM
  aduit."refTransHistory" rth
  JOIN aduit."transType" tt ON tt."transTypeId" = rth."transTypeId"
  JOIN "Users" u ON u."refUserId" = rth."refTHActionBy"
WHERE
  rth."transTypeId" IN (4, 5, 6)
  AND rth."refUserId" = ?
`
