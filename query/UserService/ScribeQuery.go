package query

var GetAllScribeDataSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refScribeDomain" rsd ON rsd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`

var UpdateScribeDomainSQL = `
UPDATE
  userdomain."refScribeDomain"
SET
  "refSDPan" = ?,
  "refSDAadhar" = ?,
  "refSDDrivingLicense" = ?
WHERE
  "refUserId" = ?
`
