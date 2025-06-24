package query

var UpdateReceptionistSQL = `
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

var UpdateReceptionistExprienceSQL = `
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

var DeleteReceptionistExprienceSQL = `
DELETE FROM
  "userdomain"."refStaffExprience"
WHERE
  "refSEId" = ?
`

var GetAllReceptionistDataSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refReceptionstDomain" rrd ON rrd."refUserId" = u."refUserId"
WHERE
  u."refUserId" =  ?
`

var UpdateReceptionistDomainSQL = `
UPDATE
  userdomain."refReceptionstDomain"
SET
  "refRDSSId" = ?,
  "refRDDrivingLicense" = ?
WHERE
  "refUserId" = ?
`
