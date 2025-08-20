package query

var ViewNotificationQuery = `
SELECT
  *
FROM
  notification."refnotification"
WHERE
  "refUserId" = $1
ORDER BY
  "refNId" DESC
LIMIT
  10
OFFSET
  $2
`

var TotalCountNotificationQuery = `
SELECT COUNT("refNId") AS total_count
FROM notification."refnotification"
WHERE "refUserId" = $1
`

var TotalReadCountNotificationQuery = `
SELECT COUNT("refNId") AS total_count
FROM notification."refnotification"
WHERE "refUserId" = $1
  AND "refNReadStatus" = false;
`

var UpdateReadMessageSQL = `
UPDATE notification.refnotification
SET "refNReadStatus" = $1
WHERE "refNId" = $2;
`
