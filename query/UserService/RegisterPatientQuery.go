package query

var CheckpatientExits = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
WHERE
  rcd."refCODOEmail" = $1
  OR rcd."refCODOPhoneNo1" = $2
  OR lower(u."refUserCustId") = lower($3)
`

var GetAllPatientDataQuery = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`

var UpdatePatientQuery = `
UPDATE
  public."Users"
SET
  "refUserFirstName" = ?,
  "refUserProfileImg" = ?,
  "refUserDOB" = ?,
  "refUserGender" = ?,
  "refUserStatus" = ?
WHERE
  "refUserId" = ?
`
