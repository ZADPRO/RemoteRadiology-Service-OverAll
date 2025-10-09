package query

var GetAllDoctorDataSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refDoctorDomain" rdd ON rdd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`

var UpdateDoctorDomainSQL = `
UPDATE
  userdomain."refDoctorDomain"
SET
  "refDDSocialSecurityNo" = ?,
  "refDDNPI" = ?,
  "refDDDrivingLicense" = ?,
  "refDDDigitalSignature" = ?,
  "refDDSpecialization" = ?,
  "refDDEaseQTReportAccess" = ?,
  "refDDNAsystemReportAccess" = ?
WHERE
  "refUserId" = ?
`

var GetMedicalLicenseSecuritySQL = `
SELECT
  *
FROM
  userdomain."refMedicalLicenseSecurity"
WHERE
  "refMLSId" = ?
`

var UpdateMedicalLicenseSecuritySQL = `
UPDATE
  userdomain."refMedicalLicenseSecurity"
SET
  "refMLSState" = ?,
  "refMLSNo" = ?,
  "refMLStatus" = ?
WHERE
  "refMLSId" = ?
`
