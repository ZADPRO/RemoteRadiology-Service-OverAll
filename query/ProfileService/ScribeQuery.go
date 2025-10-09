package query

var GetListofScribeSQL = `
SELECT
  *
FROM
  "Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refScribeDomain" rsd ON rsd."refUserId" = u."refUserId"
  JOIN "RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
  u."refRTId" = ?
  ORDER BY
  u."refUserId" DESC  
`

var GetListofScribeOneSQL = `
SELECT
  *
FROM
  "Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  JOIN userdomain."refScribeDomain" rsd ON rsd."refUserId" = u."refUserId"
  JOIN "RoleType" rt ON rt."refRTId" = u."refRTId"
WHERE
  u."refRTId" = ?
  AND u."refUserId" = ?
`
