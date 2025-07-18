package query

var RegisterUserVerifyData = `
SELECT
  rcd."refUserId",
  rcd."refCODOPhoneNo1",
  rcd."refCODOPhoneNo2",
  rcd."refCODOEmail"
FROM
  userdomain."refCommunicationDomain" rcd
WHERE
  rcd."refCODOPhoneNo1" = ?
  OR rcd."refCODOPhoneNo2" = ?
  OR rcd."refCODOEmail" = ?
`

var CreateOTPSQL = `
INSERT INTO tempotp."usersOTP" (
  "refUserId",
  "TOTPnumber",
  "TOTPexpiresAt",
  "TOTPtype"
)
VALUES (?, ?, NOW() + INTERVAL '10 minutes', ?);
`

var VerifyOTPSQL = `
SELECT
  CASE
    WHEN COUNT(*) = 1 THEN true
    ELSE false
  END AS result
FROM
  tempotp."usersOTP" uotp
WHERE
  uotp."refUserId" = ?
  AND uotp."TOTPnumber" = ?
  AND uotp."TOTPexpiresAt" > NOW()
  AND uotp."TOTPtype" = ?;
`

var DeleteOTPSQL = `
DELETE FROM tempotp."usersOTP" WHERE "refUserId" = ? AND "TOTPtype" = ?;
`
