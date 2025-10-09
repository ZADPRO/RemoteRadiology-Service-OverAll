package query

var UpdatePasswordSQL = `
UPDATE userdomain."refAuthDomain"
SET
  "refADPassword" = ?,
  "refADHashPass" = ?,
  "refAHPassChangeStatus" = ?
WHERE
  "refUserId" = ?;
`

var UpdateUserDataSQL = `
UPDATE
  public."Users"
SET
  "refUserAgreementStatus" = ?,
  "refUserConsent" = ?
WHERE
  "refUserId" = ?
`
