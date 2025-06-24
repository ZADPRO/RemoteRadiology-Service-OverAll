package query

var GetUserModel = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rc ON rc."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`
