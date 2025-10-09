package query

var GetAllManagerDataSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refManagerDomain" rmd ON rmd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`

var UpdateManagerDomainSQL = `
UPDATE
  userdomain."refManagerDomain"
SET
"refMDPan" = ?,
"refMDAadhar" = ?,
"refMDDrivingLicense" = ?
WHERE
  "refUserId" = ?
`
