package query

var LoginAdminSQL = `
SELECT
  u."refUserId",
  u."refUserCustId",
  u."refRTId",
  u."refUserFirstName",
  u."refUserLastName",
  rt."refRTName",
  rad."refADPassword",
  rad."refADHashPass",
  rad."refAHPassChangeStatus",
  rcd."refCODOPhoneNo1",
  rcd."refCODOEmail"
FROM
  public."Users" u
  JOIN userdomain."refAuthDomain" rad ON rad."refUserId" = u."refUserId"
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN public."RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
  (
    rcd."refCODOEmail" = $1
  )
  AND u."refUserStatus" = true;
`

var DeleteOTPSQL = `
DELETE FROM tempotp."usersOTP" WHERE "refUserId" = ? AND "TOTPtype" = ?;
`

var CreateOTPSQL = `
INSERT INTO tempotp."usersOTP" (
  "refUserId",
  "TOTPnumber",
  "TOTPexpiresAt",
  "TOTPtype"
)
VALUES (?, ?, NOW() AT TIME ZONE 'America/Los_Angeles' + INTERVAL '10 minutes', ?);
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
  AND uotp."TOTPexpiresAt" > NOW() AT TIME ZONE 'America/Los_Angeles'
  AND uotp."TOTPtype" = ?;
`
