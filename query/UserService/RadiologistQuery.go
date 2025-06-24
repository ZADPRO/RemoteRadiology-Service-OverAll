package query

var GetAllRadiologistDataSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refRadiologistDomain" rrd ON rrd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`

var UpdateUserSQL = `
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

var UpdateCommunicationSQL = `
UPDATE
  "userdomain"."refCommunicationDomain"
SET
  "refCODOPhoneNo1CountryCode" = ?,
  "refCODOPhoneNo1" = ?,
  "refCODOEmail" = ?
WHERE
  "refUserId" = ?;
`

var UpdateRadiologistDomainSQL = `
UPDATE
  "userdomain"."refRadiologistDomain"
SET
  "refRAMBBSRegNo" = ?,
  "refRAMDRegNo" = ?,
  "refRASpecialization" = ?,
  "refRAPan" = ?,
  "refRAAadhar" = ?,
  "refRADrivingLicense" = ?,
  "refRADigitalSignature" = ?
WHERE
  "refUserId" = ?;
`

var GetCVFilesSQL = `
SELECT
  *
FROM
  "userdomain"."refCV"
  WHERE "refCVID" = ?
  AND "refCVStatus" = true
`

var GetEducationCertificateFilesSQL = `
SELECT
  *
FROM
  userdomain."refEducationCertificate"
WHERE
  "refECId" = ?
  AND "refECStatus" = true
`

var GetLicenseFilesSQL = `
SELECT
  *
FROM
  "userdomain"."refLicences"
WHERE
  "refUserId" = ?
  AND "refLStatus" = true
`

var GetMalpracticeFilesSQL = `
SELECT
  *
FROM
  "userdomain"."refMalpractice"
WHERE
  "refUserId" = ?
  AND "refMPStatus" = true
`

var UpdateCVFilesSQL = `
UPDATE
  "userdomain"."refCV"
SET
  "refCVFileName" = ?,
  "refCVOldFileName" = ?,
  "refCVStatus" = ?
WHERE
  "refCVID" = ?;
`

var UpdateEductionCertificateFilesSQL = `
UPDATE
  userdomain."refEducationCertificate"
SET
  "refECFileName" = ?,
  "refECOldFileName" = ?,
  "refECStatus" = ?
WHERE
  "refECId" = ?
`

var UpdateLicenseFilesSQL = `
UPDATE
  "userdomain"."refLicences"
SET
  "refLFileName" = ?,
  "refLOldFileName" = ?,
  "refLStatus" = ?
WHERE
  "refLId" = ?;
`

var UpdateMalpracticeFilesSQL = `
UPDATE
  userdomain."refMalpractice"
SET
  "refMPFileName" = ?,
  rm."refMPOldFileName" = ?,
  "refMPStatus" = ?
WHERE
  "refMPId" = ?
`
