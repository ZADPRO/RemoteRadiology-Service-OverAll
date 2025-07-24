package query

var WGPPGetAllRadiologistDataSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refRadiologistDomain" rrd ON rrd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`

var WGPPUpdateUserSQL = `
UPDATE
  public."Users"
SET
  "refUserFirstName" = ?,
  "refUserLastName" = ?,
  "refUserProfileImg" = ?,
  "refUserDOB" = ?,
  "refUserStatus" = ?
WHERE
  "refUserId" = ?
`

var WGPPUpdateCommunicationSQL = `
UPDATE
  "userdomain"."refCommunicationDomain"
SET
  "refCODOPhoneNo1CountryCode" = ?,
  "refCODOPhoneNo1" = ?,
  "refCODOEmail" = ?
WHERE
  "refUserId" = ?;
`

var WGPPUpdateRadiologistDomainSQL = `
UPDATE
  "userdomain"."refWellthgreenPerformingProvider"
SET
  "refWGPPMBBSRegNo" = ?,
  "refWGPPMDRegNo" = ?,
  "refWGPPSpecialization" = ?,
  "refWGPPPan" = ?,
  "refWGPPAadhar" = ?,
  "refWGPPDrivingLicense" = ?,
  "refWGPPDigitalSignature" = ?
WHERE
  "refUserId" = ?;
`

var WGPPGetCVFilesSQL = `
SELECT
  *
FROM
  "userdomain"."refCV"
  WHERE "refCVID" = ?
  AND "refCVStatus" = true
`

var WGPPGetEducationCertificateFilesSQL = `
SELECT
  *
FROM
  userdomain."refEducationCertificate"
WHERE
  "refECId" = ?
  AND "refECStatus" = true
`

var WGPPGetLicenseFilesSQL = `
SELECT
  *
FROM
  "userdomain"."refLicences"
WHERE
  "refUserId" = ?
  AND "refLStatus" = true
`

var WGPPGetMalpracticeFilesSQL = `
SELECT
  *
FROM
  "userdomain"."refMalpractice"
WHERE
  "refUserId" = ?
  AND "refMPStatus" = true
`

var WGPPUpdateCVFilesSQL = `
UPDATE
  "userdomain"."refCV"
SET
  "refCVFileName" = ?,
  "refCVOldFileName" = ?,
  "refCVStatus" = ?
WHERE
  "refCVID" = ?;
`

var WGPPUpdateEductionCertificateFilesSQL = `
UPDATE
  userdomain."refEducationCertificate"
SET
  "refECFileName" = ?,
  "refECOldFileName" = ?,
  "refECStatus" = ?
WHERE
  "refECId" = ?
`

var WGPPUpdateLicenseFilesSQL = `
UPDATE
  "userdomain"."refLicences"
SET
  "refLFileName" = ?,
  "refLOldFileName" = ?,
  "refLStatus" = ?
WHERE
  "refLId" = ?;
`

var WGPPUpdateMalpracticeFilesSQL = `
UPDATE
  userdomain."refMalpractice"
SET
  "refMPFileName" = ?,
  rm."refMPOldFileName" = ?,
  "refMPStatus" = ?
WHERE
  "refMPId" = ?
`
