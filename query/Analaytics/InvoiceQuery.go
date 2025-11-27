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
  "refTASform" = $1,
  "refTADaform" = $2,
  "refTADbform" = $3,
  "refTADcform" = $4,
  "refTAXform" = $5,
  "refTAEditform" = $6,
  "refTADScribeTotalcase" = $7
WHERE
  "refTAId" = $8
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
SELECT $1 AS "refSCId", COUNT(*) AS total_appointments
FROM (SELECT *
      FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                        rrh."refRHHandleEndTime",
                                                        rrh."refRHHandleStatus",
                                                        ra."refCategoryId",
                                                        ra."refSCId"
            FROM notes."refReportsHistory" rrh
                     JOIN appointment."refAppointments" ra
                          ON ra."refAppointmentId" = rrh."refAppointmentId"
            WHERE ra."refAppointmentStatus" = TRUE
              AND rrh."refRHHandleStatus" = 'Signed Off'
              AND ra."refSCId" = $1
            ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
      WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                BETWEEN date_trunc('month', to_date($2, 'YYYY-MM'))
                AND (
              date_trunc('month', to_date($2, 'YYYY-MM'))
                  + INTERVAL '1 month - 1 day'
              )) t;
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
  "refIHSFormquantity",
  "refIHSFormamount",
  "refIHDaFormquantity",
  "refIHDaFormamount",
  "refIHDbFormquantity",
  "refIHDbFormamount",
  "refIHDcFormquantity",
  "refIHDcFormamount",
  "refIHxFormquantity",
  "refIHxFormamount",
  "refIHEditquantity",
  "refIHEditFormamount",
  "refIHScribeTotalcasequantity",
  "refIHScribeTotalcaseamount",
  "refIHAddtionalAmount",
  "refIHDeductibleAmount",
  "refIHScanCenterTotalCase",
  "refIHScancentercaseAmount",
  "refIHTotal"
)
VALUES (
  $1,  $2,  $3,  $4,  $5,  $6,  $7,  $8,  $9,  $10,
  $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
  $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
  $31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
  $41, $42, $43
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
  invoice."otherInvoiceAmount" ("refIHId", "refOIAName", "refOIAAmount", "refOIAAmountType")
VALUES
  ($1, $2, $3, $4);
`

var GetInvoiceOtherAmount = `
SELECT
  *
FROM
  invoice."otherInvoiceAmount"
WHERE
  "refIHId" = $1
`
