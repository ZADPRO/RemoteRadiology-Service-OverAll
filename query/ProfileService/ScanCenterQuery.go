package query

var GetAllScanCenter = `
SELECT
  *
FROM
  public."ScanCenter"
ORDER BY
  "refSCId" DESC
`

var GetScanCenter = `
SELECT
  *
FROM
  public."ScanCenter"
WHERE
  "refSCId" = ?
ORDER BY
  "refSCId" DESC
`

var IdentifyScanCenterMapping = `
SELECT
  *
FROM
  map."refScanCenterMap" rscm
  JOIN public."ScanCenter" sc ON sc."refSCId" = rscm."refSCId"
WHERE
  rscm."refUserId" = ?
`
