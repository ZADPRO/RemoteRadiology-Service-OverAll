package query

var VerifyData = `
SELECT
  *
FROM
  userdomain."refCommunicationDomain" rc
WHERE
  rc."refCODOPhoneNo1" IN (?, ?)
  OR rc."refCODOPhoneNo2" IN (?, ?)
  OR rc."refCODOEmail" = ?
`

var VerifyDataSQL = `
SELECT
  *
FROM
  userdomain."refCommunicationDomain" rc
WHERE
  rc."refCODOPhoneNo1" = ?
  OR rc."refCODOEmail" = ?
`

var ScanCenterVerifyDataSQL = `
SELECT
  *
FROM
  public."ScanCenter"
WHERE
  "refSCPhoneNo1" = ?
  OR "refSCEmail" = ?
  OR "refSCCustId" = ?
`

var UpdateScanCenterVerifyDataSQL = `
SELECT
  *
FROM
  public."ScanCenter"
WHERE
  "refSCPhoneNo1" = ?
  OR "refSCEmail" = ?
`

var GetUsersCountSQL = `
SELECT COUNT(*) as "TotalCount" FROM public."Users" WHERE "refRTId" = ?
`

var GetUsersScanCountSQL = `
SELECT
  COUNT(u."refUserId") AS "TotalCount",
  sc."refSCCustId"
FROM
  public."ScanCenter" sc
  LEFT JOIN map."refScanCenterMap" rscm ON sc."refSCId" = rscm."refSCId"
  LEFT JOIN public."Users" u ON u."refUserId" = rscm."refUserId"
  AND u."refRTId" = ?
WHERE
  sc."refSCId" = ?
GROUP BY
  sc."refSCCustId";
`

var UpdateTechnicianSQL = `
UPDATE
  "public"."Users"
SET
  "refUserFirstName" = ?,
  "refUserLastName" = ?,
  "refUserProfileImg" = ?,
  "refUserDOB" = ?,
  "refUserGender" = ?,
  "refUserStatus" = ?
WHERE
  "refUserId" = ?;
`

var UpdateTechnicianExprienceSQL = `
UPDATE
  "userdomain"."refStaffExprience"
SET
  "refSEHospitalName" = ?,
  "refSEDesignation" = ?,
  "refSESpecialization" = ?,
  "refSEAddress" = ?,
  "refSEFrom" = ?,
  "refSETo" = ?
WHERE
  "refSEId" = ?;
`

var DeleteTechnicianExprienceSQL = `
DELETE FROM
  "userdomain"."refStaffExprience"
WHERE
  "refSEId" = ?
`

var GetAllTechnicianDataSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refTechnicianDomain" trd ON trd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`
var GetStaffExprienceDataSQl = `
SELECT
  *
FROM
  userdomain."refStaffExprience"
WHERE
  "refSEId" = ?
`
var GetMapDataSQL = `
SELECT
  *
FROM
  map."refScanCenterMap"
WHERE
  "refUserId" = ?
  AND "refSCId" = ?
  `
var UpdateMapSQL = `
  UPDATE
  "map"."refScanCenterMap"
SET
  "refRTId" = ?,
  "refSCMStatus" = ?
WHERE
  "refUserId" = ?
  AND "refSCId" = ?
  `

var GetScanCenterNameSQL = `
SELECT
  "refSCName"
FROM
  "ScanCenter"
WHERE
  "refSCId" = ?
  `

var UpdateTechnicianDomainSQL = `
UPDATE
  userdomain."refTechnicianDomain"
SET
  "refTDTrainedEaseQTStatus" = ?,
  "refTDSSNo" = ?,
  "refTDDrivingLicense" = ?,
  "refTDDigitalSignature" = ?
WHERE
  "refUserId" = ?
`
