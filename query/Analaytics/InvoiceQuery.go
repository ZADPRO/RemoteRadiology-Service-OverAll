package query

var GetAmountSQL = `
SELECT
  *
FROM
  invoice."totalAmount"
WHERE
  "refTAId" = 1
`

var ListAllScanCenter = `
SELECT
  *
FROM
  public."ScanCenter"
ORDER BY
  "refSCId"
`

var ListUserSQL = `
SELECT
  *
FROM
  public."Users"
WHERE
  "refRTId" IN (6, 7)
ORDER BY
  "refRTId" DESC
`

var UpdateAmountSQL = `
UPDATE
  invoice."totalAmount"
SET
  "refTAAmountScanCenter" = ?,
  "refTAAmountUser" = ?
WHERE
  "refTAId" = 1
`

var ListOneScanCenter = `
SELECT
  *
FROM
  public."ScanCenter"
WHERE
  "refSCId" = ?
`

var GetScanCenterCountPerMonthSQL = `
SELECT
  sc."refSCId",
  COUNT(ra."refAppointmentId") AS total_appointments
FROM
  public."ScanCenter" sc
  LEFT JOIN appointment."refAppointments" ra ON ra."refSCId" = sc."refSCId"
  AND ra."refAppointmentComplete" = 'Signed Off'
  AND TO_CHAR(
    TO_DATE(ra."refAppointmentDate", 'YYYY-MM-DD'),
    'YYYY-MM'
  ) = ?
WHERE
  sc."refSCId" = ?
GROUP BY
  sc."refSCId";
`

var GetOneUserSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?
`

var InsertInvoiceSQL = `
INSERT INTO
  invoice."invoiceHistory" (
    "refSCId",
    "refUserId",
    "refIHFromId",
    "refIHFromName",
    "refIHFromPhoneNo",
    "refIHFromEmail",
    "refIHFromPan",
    "refIHFromGST",
    "refIHFromAddress",
    "refIHToId",
    "refIHToName",
    "refIHFromDate",
    "refIHToDate",
    "refIHModePayment",
    "refIHUPIId",
    "refIHAccountHolderName",
    "refIHAccountNo",
    "refIHAccountBank",
    "refIHAccountBranch",
    "refIHAccountIFSC",
    "refIHQuantity",
    "refIHAmount",
    "refIHTotal",
    "refIHCreatedAt",
    "refIHCreatedBy",
    "refIHToAddress"
  )
VALUES
  (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    NOW(),
    ?,
    ?
  );
`

var GetInvoiceHistoryScancenterSQL = `
SELECT
  *
FROM
  invoice."invoiceHistory"
WHERE
  "refSCId" = ?
`

var GetInvoiceHistoryUserSQL = `
SELECT
  *
FROM
  invoice."invoiceHistory"
WHERE
  "refUserId" = ?
`
