package query

var GetAllCoDoctorDataSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refCoDoctorDomain"  rcdd ON rcdd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`

var UpdateCoDoctorDomainSQL = `
UPDATE
  userdomain."refCoDoctorDomain"
SET
  "refCDSocialSecurityNo" = ?,
  "refCDNPI" = ?,
  "refCDDrivingLicense" = ?,
  "refCDDigitalSignature" = ?,
  "refCDSpecialization" = ?,
  "refCDEaseQTReportAccess" = ?
WHERE
  "refUserId" = ?
`
