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
  "refRTId" IN (6, 7, 10)
ORDER BY
  "refRTId" DESC
`

var UpdateAmountSQL = `
UPDATE
  invoice."totalAmount"
SET
  "refTASformEdit" = $1,
  "refTASformCorrect" = $2,
  "refTADaformEdit" = $3,
  "refTADaformCorrect" = $4,
  "refTADbformEdit" = $5,
  "refTADbformCorrect" = $6,
  "refTADcformEdit" = $7,
  "refTADcformCorrect" = $8,
  "refTADScribeTotalcase" = $9
WHERE
  "refTAId" = $10
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
  -- AND ra."refAppointmentComplete" = 'Signed Off'
  AND TO_CHAR(
    TO_DATE(ra."refAppointmentDate", 'YYYY-MM-DD'),
    'YYYY-MM'
  ) = ?
WHERE
  sc."refSCId" = ?
GROUP BY
  sc."refSCId";
`

// SELECT
//   *
// FROM
//   public."Users" u
//   JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
// WHERE
//   u."refUserId" = ?

var GetOneUserSQL = `
SELECT
  u.*,
  rcd.*,
  CASE
    WHEN u."refRTId" = 7 THEN rsd."refSDPan"
    WHEN u."refRTId" = 6 THEN rradd."refRAPan"
    WHEN u."refRTId" = 10 THEN rwpp."refWGPPPan"
    ELSE NULL
  END AS "refPan"
FROM
  public."Users" u
  JOIN userdomain."refCommunicationDomain" rcd ON rcd."refUserId" = u."refUserId"
  FULL JOIN userdomain."refScribeDomain" rsd ON rsd."refUserId" = u."refUserId"
  FULL JOIN userdomain."refRadiologistDomain" rradd ON rradd."refUserId" = u."refUserId"
  FULL JOIN userdomain."refWellthgreenPerformingProvider" rwpp ON rwpp."refUserId" = u."refUserId"
WHERE
  u."refUserId" = ?;
`

var InsertInvoiceSQL = `
INSERT INTO invoice."invoiceHistory" (
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
  "refIHCreatedAt",
  "refIHCreatedBy",
  "refIHToAddress",
  "refIHSignature",
  "refIHSformEditquantity",
  "refIHSformEditamount",
  "refIHSformCorrectquantity",
  "refIHSformCorrectamount",
  "refIHDaformEditquantity",
  "refIHDaformEditamount",
  "refIHDaformCorrectquantity",
  "refIHDaformCorrectamount",
  "refIHDbformEditquantity",
  "refIHDbformEditamount",
  "refIHDbformCorrectquantity",
  "refIHDbformCorrectamount",
  "refIHDcformEditquantity",
  "refIHDcformEditamount",
  "refIHDcformCorrectquantity",
  "refIHDcformCorrectamount",
  "refIHScribeTotalcasequantity",
  "refIHScribeTotalcaseamount",
  "refIHOtherExpensiveName",
  "refIHOtherAmount",
  "refIHScanCenterTotalCase",
  "refIHScancentercaseAmount",
  "refIHTotal"
)
VALUES (
  $1,  $2,  $3,  $4,  $5,  $6,  $7,  $8,  $9,  $10,
  $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
  $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
  $31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
  $41, $42, $43, $44, $45, $46, $47
)
RETURNING "refIHId";
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
  invoice."invoiceHistory" ih
  JOIN public."Users" u ON u."refUserId" = ih."refUserId"
WHERE
  ih."refUserId" = ?
`

var GetInvoiceOverAllHistorySQL = `
SELECT
  ih.*,
  u."refUserCustId",
  u."refRTId",
  sc."refSCCustId"
FROM
  invoice."invoiceHistory" ih
  LEFT JOIN public."Users" u ON u."refUserId" = ih."refUserId"
  LEFT JOIN public."ScanCenter" sc ON sc."refSCId" = ih."refSCId"
WHERE
  (
    $1 = 1
    OR u."refRTId" IN (1, 6, 7, 10)
  )
  AND 
  (
    $2 = ''
    OR ih."refIHCreatedAt"::timestamp >= $2::timestamp
  )
  AND (
    $3 = ''
    OR ih."refIHCreatedAt"::timestamp <= $3::timestamp
  );
`

var InsertOtherInvoiceAmount = `
INSERT INTO
  invoice."otherInvoiceAmount" ("refIHId", "refOIAName", "refOIAAmount")
VALUES
  ($1, $2, $3);
`

var GetInvoiceOtherAmount = `
SELECT
  *
FROM
  invoice."otherInvoiceAmount"
WHERE
  "refIHId" = $1
`
