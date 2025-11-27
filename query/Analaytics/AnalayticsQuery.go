package query

var AdminOverallAnalayticsSQL = `
WITH months AS (SELECT date_trunc('month', CURRENT_DATE) - interval '1 month' * n AS month_start
                FROM generate_series(0, 5) AS n)
SELECT TO_CHAR(month_start, 'YYYY-MM')                                                                        AS month,
       TO_CHAR(month_start, 'Month')                                                                          AS month_name,
       month_start::date                                                                                      AS starting_date,
       (month_start + INTERVAL '1 month - 1 day')::date                                                       AS ending_date,
       (SELECT COUNT(*) AS total_rows
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND (
                        $1 = 0
                            OR ra."refSCId" = $1
                        )
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN month_start::date AND (month_start + INTERVAL '1 month - 1 day')::date) AS t) AS total_appointments
FROM months
ORDER BY month_start;
`

// var AdminOverallScanIndicatesAnalayticsSQL = `
// SELECT
//   COALESCE(SUM(total_appointments), 0) AS total_appointments,
//   COALESCE(SUM("SForm"), 0) AS "SForm",
//   COALESCE(SUM("DaForm"), 0) AS "DaForm",
//   COALESCE(SUM("DbForm"), 0) AS "DbForm",
//   COALESCE(SUM("DcForm"), 0) AS "DcForm"
// FROM (
//   SELECT
//     COUNT(*) AS total_appointments,
//     COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END) AS "SForm",
//     COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END) AS "DaForm",
//     COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END) AS "DbForm",
//     COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END) AS "DcForm"
//   FROM
//     appointment."refAppointments"
//   WHERE
//     TO_CHAR(TO_DATE("refAppointmentDate", 'YYYY-MM-DD'), 'YYYY-MM') = ?
//     AND (
//       ? = 0
//       OR "refSCId" = ?
//     )
// ) AS stats;
// `

var AdminOverallScanIndicatesAnalayticsSQL = `
SELECT COUNT(*)                                        AS total_appointments,
       COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END) AS "SForm",
       COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END) AS "DaForm",
       COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END) AS "DbForm",
       COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END) AS "DcForm"
FROM (SELECT *
      FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                        rrh."refRHHandleEndTime",
                                                        rrh."refRHHandleStatus",
                                                        ra."refCategoryId" -- ✅ REQUIRED for outer COUNT()
            FROM notes."refReportsHistory" rrh
                     JOIN appointment."refAppointments" ra
                          ON ra."refAppointmentId" = rrh."refAppointmentId"
            WHERE ra."refAppointmentStatus" = TRUE
              AND rrh."refRHHandleStatus" = 'Signed Off'
              AND (
                $3 = 0
                    OR ra."refSCId" = $3
                )
            ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
      WHERE TO_TIMESTAMP(s."refRHHandleEndTime"
                , 'YYYY-MM-DD HH24:MI:SS')::date
                BETWEEN $1::date
                AND $2::date) AS t;
`

var AdminOverallUserIndicatesAnalayticsSQL = `
SELECT COUNT(*)                                        AS total_appointments,
       COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END) AS "SForm",
       COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END) AS "DaForm",
       COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END) AS "DbForm",
       COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END) AS "DcForm"
FROM (SELECT *
      FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                        rrh."refRHHandleEndTime",
                                                        rrh."refRHHandleStatus",
                                                        rrh."refRHHandledUserId",
                                                        u."refRTId",
                                                        ra."refCategoryId"
            FROM notes."refReportsHistory" rrh
                     JOIN appointment."refAppointments" ra
                          ON ra."refAppointmentId" = rrh."refAppointmentId"
                     JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
            WHERE ra."refAppointmentStatus" = TRUE
              AND (
                $3 = 0
                    OR ra."refSCId" = $3
                )
              AND rrh."refRHHandleStatus" IN
                  ('Signed Off')
            ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
      WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                BETWEEN $1::date AND $2::date) AS t;
`

var FindSCIdSQL = `
SELECT
  "refSCId" as "ScancenterId"
FROM
  map."refScanCenterMap"
WHERE
  "refUserId" = $1
`

var GetAllScanCenter = `
SELECT
  *
FROM
  public."ScanCenter"
ORDER BY
  "refSCId" DESC
`

var UserListIdsSQL = `
SELECT
  *
FROM
  public."Users"
WHERE
  "refRTId" IN ?
ORDER BY
  "refRTId",
  "refUserId";
`

var ScanCenterUserListIdsSQL = `
SELECT
  *
FROM
  public."Users" u
  JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
WHERE
  u."refRTId" IN ? AND rscm."refSCId" = ?
  ORDER BY
  u."refRTId",
  u."refUserId";
`

var WellGreenUserAnalayticsSQL = `
WITH months AS (SELECT date_trunc('month', CURRENT_DATE) - interval '1 month' * n AS month_start
                FROM generate_series(0, 5) AS n)
SELECT TO_CHAR(month_start, 'YYYY-MM')                                                                        AS month,
       TO_CHAR(month_start, 'Month')                                                                          AS month_name,
       month_start::date                                                                                      AS starting_date,
       (month_start + INTERVAL '1 month - 1 day')::date                                                       AS ending_date,
       (SELECT COUNT(*) AS total_rows
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = ?
                                  AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
              WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN month_start::date AND (month_start + INTERVAL '1 month - 1 day')::date) AS t) AS total_appointments
FROM months
ORDER BY month_start;
`

var WellGreenUserIndicatesAnalayticsInvoiceSQL = `
SELECT COUNT(*)                                            AS total_appointments,
       COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END)     AS "SForm",
       COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END)     AS "DaForm",
       COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END)     AS "DbForm",
       COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END)     AS "DcForm",
       COUNT(CASE WHEN "refCategoryId" IS NULL THEN 1 END) AS "xForm",
       COUNT(CASE WHEN "refRHHandleEdit" = 1 THEN 1 END)   AS "editForm"
FROM (SELECT *
      FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                        rrh."refRHHandleEndTime",
                                                        rrh."refRHHandleStatus",
                                                        rrh."refRHHandledUserId",
                                                        u."refRTId",
                                                        rA."refAppointmentId",
                                                        rA."refCategoryId",
                                                        rrh."refRHHandleEdit"
            FROM notes."refReportsHistory" rrh
                     JOIN appointment."refAppointments" ra
                          ON ra."refAppointmentId" = rrh."refAppointmentId"
                     JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
            WHERE ra."refAppointmentStatus" = TRUE
              AND rrh."refRHHandleStatus" = 'Signed Off'
              AND EXISTS (SELECT 1
                          FROM notes."refReportsHistory" rrhi
                                   JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                          WHERE (
                              ui."refRTId" NOT IN (2)
                                  OR
                              (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                              )
                            AND (
                              ui."refRTId" NOT IN (8)
                                  OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                              )
                            AND (
                              ui."refRTId" NOT IN (7)
                                  OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                              )
                            AND (
                              ui."refRTId" NOT IN (6)
                                  OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                              )
                            AND (
                              ui."refRTId" NOT IN (1, 5, 10)
                                  OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                              )
                            AND rrhi."refRHHandledUserId" = $1
                            AND rrhi."refAppointmentId" = rrh."refAppointmentId")
            ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
      WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                BETWEEN $2::date AND ($2::date + INTERVAL '1 month - 1 day')::date) AS t
`

var WellGreenUserIndicatesAnalayticsSQL = `
SELECT COUNT(*)                                        AS total_appointments,
       COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END) AS "SForm",
       COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END) AS "DaForm",
       COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END) AS "DbForm",
       COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END) AS "DcForm"
FROM (SELECT *
      FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                        rrh."refRHHandleEndTime",
                                                        rrh."refRHHandleStatus",
                                                        rrh."refRHHandledUserId",
                                                        u."refRTId",
                                                        rA."refAppointmentId",
                                                        rA."refCategoryId"
            FROM notes."refReportsHistory" rrh
                     JOIN appointment."refAppointments" ra
                          ON ra."refAppointmentId" = rrh."refAppointmentId"
                     JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
            WHERE ra."refAppointmentStatus" = TRUE
              AND rrh."refRHHandleStatus" = 'Signed Off'
              AND EXISTS (SELECT 1
                          FROM notes."refReportsHistory" rrhi
                                   JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                          WHERE (
                              ui."refRTId" NOT IN (2)
                                  OR
                              (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                              )
                            AND (
                              ui."refRTId" NOT IN (8)
                                  OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                              )
                            AND (
                              ui."refRTId" NOT IN (7)
                                  OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                              )
                            AND (
                              ui."refRTId" NOT IN (6)
                                  OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                              )
                            AND (
                              ui."refRTId" NOT IN (1, 5, 10)
                                  OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                              )
                            AND rrhi."refRHHandledUserId" = ?
                            AND rrhi."refAppointmentId" = rrh."refAppointmentId")
            ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
      WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                BETWEEN ?::date AND ?::date) AS t
`

var UserWorkedTimingSQL = `
SELECT sub."refRHHandledUserId",
       SUM(
               EXTRACT(EPOCH FROM (sub."refRHHandleEndTime"::timestamp - sub."refStartTime"::timestamp)) / 60
       ) AS total_minutes,
       SUM(
               EXTRACT(EPOCH FROM (sub."refRHHandleEndTime"::timestamp - sub."refStartTime"::timestamp)) / 3600
       ) AS total_hours
FROM (SELECT *
      FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                        rrh."refRHHandleEndTime",
                                                        rrh."refRHHandleStatus",
                                                        rrh."refRHHandledUserId",
                                                        u."refRTId",
                                                        CASE
                                                            WHEN EXISTS (SELECT 1
                                                                         FROM notes."refReportsHistory" rh2
                                                                         WHERE rh2."refRHHandleStatus" = 'Technologist Form Fill'
                                                                           AND rh2."refAppointmentId" = rrh."refAppointmentId"
                                                                           AND rh2."refRHHandleEndTime" IS NOT NULL)
                                                                THEN (SELECT CAST(rh2."refRHHandleEndTime" AS timestamp)
                                                                      FROM notes."refReportsHistory" rh2
                                                                      WHERE rh2."refRHHandleStatus" = 'Technologist Form Fill'
                                                                        AND rh2."refAppointmentId" = rrh."refAppointmentId"
                                                                        AND rh2."refRHHandleEndTime" IS NOT NULL
                                                                      ORDER BY rh2."refRHId" DESC
                                                                      LIMIT 1)
                                                            ELSE (SELECT CAST(rdf."refDFCreatedAt" AS timestamp)
                                                                  FROM dicom."refDicomFiles" rdf
                                                                  WHERE rdf."refAppointmentId" = rrh."refAppointmentId"
                                                                  ORDER BY rdf."refDFId" DESC
                                                                  LIMIT 1)
                                                            END AS "refStartTime"
            FROM notes."refReportsHistory" rrh
                     JOIN appointment."refAppointments" ra
                          ON ra."refAppointmentId" = rrh."refAppointmentId"
                     JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
            WHERE ra."refAppointmentStatus" = TRUE
              AND rrh."refRHHandleStatus" = 'Signed Off'
              AND EXISTS (SELECT 1
                          FROM notes."refReportsHistory" rrhi
                                   JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                          WHERE (
                              ui."refRTId" NOT IN (2)
                                  OR
                              (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                              )
                            AND (
                              ui."refRTId" NOT IN (8)
                                  OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                              )
                            AND (
                              ui."refRTId" NOT IN (7)
                                  OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                              )
                            AND (
                              ui."refRTId" NOT IN (6)
                                  OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                              )
                            AND (
                              ui."refRTId" NOT IN (1, 5, 10)
                                  OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                              )
                            AND rrhi."refRHHandledUserId" = ?
                            AND rrhi."refAppointmentId" = rrh."refAppointmentId")
            ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
      WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                BETWEEN ?::date AND ?::date) sub
WHERE sub."refStartTime" IS NOT NULL
  AND sub."refRHHandleEndTime" IS NOT NULL
GROUP BY sub."refRHHandledUserId";
`

// var ListScanAppointmentCountSQL = `
// SELECT
//   sc."refSCId",
//   sc."refSCName",
//   COUNT(DISTINCT rrh."refAppointmentId") AS total_appointments
// FROM
//   public."ScanCenter" sc
// LEFT JOIN
//   appointment."refAppointments" a ON sc."refSCId" = a."refSCId"
// LEFT JOIN
//   notes."refReportsHistory" rrh ON rrh."refAppointmentId" = a."refAppointmentId"
//   AND rrh."refRHHandledUserId" = ?
//   AND TO_CHAR(TO_TIMESTAMP(rrh."refRHHandleStartTime", 'YYYY-MM-DD HH24:MI:SS'), 'YYYY-MM') = ?
// GROUP BY
//   sc."refSCId", sc."refSCName"
// ORDER BY
//   sc."refSCName";
// `

var ListScanAppointmentCountSQL = `
SELECT sc."refSCId",
       sc."refSCName",
       (SELECT COUNT(*) AS total_appointments
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                rA."refAppointmentId",
                                                                rA."refCategoryId"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND ra."refSCId" = sc."refSCId"
AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = ?
                                  AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
              WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN ?::date AND ?::date) AS t) AS total_appointments
FROM public."ScanCenter" sc
`

// var TotalCorrectEditSQL = `
// SELECT
//   SUM(COALESCE("refRHHandleCorrect", 0)) AS "totalCorrect",
//   SUM(COALESCE("refRHHandleEdit", 0)) AS "totalEdit"
// FROM
//   notes."refReportsHistory"
// WHERE
//   "refRHHandledUserId" = ?
//   AND TO_CHAR("refRHHandleStartTime"::timestamp, 'YYYY-MM') = ?;
// `

var TotalCorrectEditSQL = `
SELECT COALESCE(SUM(COALESCE(sub."totalCorrect", 0)), 0) AS "totalCorrect",
       COALESCE(SUM(COALESCE(sub."totalEdit", 0)), 0)    AS "totalEdit"
FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                  rrh."refRHHandleEndTime",
                                                  rrh."refRHHandleStatus",
                                                  rrh."refRHHandledUserId",
                                                  u."refRTId",
                                                  rrh."refRHHandleCorrect",
                                                  rrh."refRHHandleEdit",
                                                  (SELECT newCrt."refRHHandleCorrect"
                                                   FROM (SELECT DISTINCT ON (rrhi."refAppointmentId") *
                                                         FROM notes."refReportsHistory" rrhi
                                                                  JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                                         WHERE (
                                                             ui."refRTId" NOT IN (2)
                                                                 OR
                                                             (ui."refRTId" IN (2) AND
                                                              rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                                             )
                                                           AND (
                                                             ui."refRTId" NOT IN (8)
                                                                 OR
                                                             (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                                             )
                                                           AND (
                                                             ui."refRTId" NOT IN (7)
                                                                 OR
                                                             (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                                             )
                                                           AND (
                                                             ui."refRTId" NOT IN (6)
                                                                 OR
                                                             (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                                             )
                                                            AND (
                                                              ui."refRTId" NOT IN (1, 5)
                                                                  OR
                                                              (ui."refRTId" IN (1, 5) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                                              )
                                                            AND (ui."refRTId" NOT IN (10)
                                                              OR
                                                                (ui."refRTId" IN (10) AND rrhi."refRHHandleStatus" = 'Reviewed 1')
                                                            )
                                                           AND rrhi."refRHHandledUserId" = $1
                                                         ORDER BY rrhi."refAppointmentId", rrhi."refRHId" DESC) newCrt
                                                   WHERE newCrt."refAppointmentId" = rrh."refAppointmentId") AS "totalCorrect",
                                                  (SELECT newCrt."refRHHandleEdit"
                                                   FROM (SELECT DISTINCT ON (rrhi."refAppointmentId") *
                                                         FROM notes."refReportsHistory" rrhi
                                                                  JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                                         WHERE (
                                                             ui."refRTId" NOT IN (2)
                                                                 OR
                                                             (ui."refRTId" IN (2) AND
                                                              rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                                             )
                                                           AND (
                                                             ui."refRTId" NOT IN (8)
                                                                 OR
                                                             (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                                             )
                                                           AND (
                                                             ui."refRTId" NOT IN (7)
                                                                 OR
                                                             (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                                             )
                                                           AND (
                                                             ui."refRTId" NOT IN (6)
                                                                 OR
                                                             (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                                             )
                                                           AND (
                                                              ui."refRTId" NOT IN (1, 5)
                                                                  OR
                                                              (ui."refRTId" IN (1, 5) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                                              )
                                                            AND (ui."refRTId" NOT IN (10)
                                                              OR
                                                                (ui."refRTId" IN (10) AND rrhi."refRHHandleStatus" = 'Reviewed 1')
                                                            )
                                                           AND rrhi."refRHHandledUserId" = $1
                                                         ORDER BY rrhi."refAppointmentId", rrhi."refRHId" DESC) newCrt
                                                   WHERE newCrt."refAppointmentId" = rrh."refAppointmentId") AS "totalEdit"
      FROM notes."refReportsHistory" rrh
               JOIN appointment."refAppointments" ra
                    ON ra."refAppointmentId" = rrh."refAppointmentId"
               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
      WHERE ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND EXISTS (SELECT 1
                    FROM notes."refReportsHistory" rrhi
                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                    WHERE (
                        ui."refRTId" NOT IN (2)
                            OR
                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                        )
                      AND (
                        ui."refRTId" NOT IN (8)
                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                        )
                      AND (
                        ui."refRTId" NOT IN (7)
                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                        )
                      AND (
                        ui."refRTId" NOT IN (6)
                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                        )
                      AND (
                          ui."refRTId" NOT IN (1, 5)
                              OR
                          (ui."refRTId" IN (1, 5) AND rrhi."refRHHandleStatus" = 'Signed Off')
                          )
                        AND (ui."refRTId" NOT IN (10)
                          OR
                            (ui."refRTId" IN (10) AND rrhi."refRHHandleStatus" = 'Reviewed 1')
                        )
                      AND rrhi."refRHHandledUserId" = $1
                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) sub
WHERE TO_TIMESTAMP(sub."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
          BETWEEN $2::date AND $3::date;
`

var ImpressionNRecommentationScanCenterSQL = `
WITH
  latest_report_history AS (
    SELECT
      *
    FROM
      (
        SELECT
          *,
          ROW_NUMBER() OVER (
            PARTITION BY
              "refAppointmentId"
            ORDER BY
              "refRHId" DESC
          ) AS rn
        FROM
          notes."refReportsHistory"
      ) sub
    WHERE
      rn = 1
  ),
  actual_counts AS (
    SELECT
      ra."refAppointmentImpression" AS impression,
      COUNT(*) AS count
    FROM
      latest_report_history rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
      JOIN public."Users" rrhu ON rrh."refRHHandledUserId" = rrhu."refUserId"
    WHERE
      ra."refAppointmentDate" >= ?
      AND ra."refAppointmentDate" <= ?
      AND (
        ? = 0
        OR ra."refSCId" = ?
      )
      AND (
        ? = FALSE
        OR rrhu."refRTId" IN (1, 6, 7, 10)
      )
    GROUP BY
      ra."refAppointmentImpression"
  ),
  expected_impressions AS (
    SELECT DISTINCT
      "refIRVCustId" AS impression
    FROM
      impressionrecommendation."ImpressionRecommendationVal"
    WHERE
      "refIRVSystemType" = 'WR'
  )
SELECT
  ei.impression,
  COALESCE(ac.count, 0) AS count
FROM
  expected_impressions ei
  LEFT JOIN actual_counts ac ON ei.impression = ac.impression
ORDER BY
  ei.impression;
`

// var ImpressionNRecommentationSQL = `
// WITH
//   latest_report_history AS (
//     SELECT
//       *
//     FROM
//       (
//         SELECT
//           *,
//           ROW_NUMBER() OVER (
//             PARTITION BY
//               "refAppointmentId"
//             ORDER BY
//               "refRHId" DESC
//           ) AS rn
//         FROM
//           notes."refReportsHistory"
//         WHERE
//           "refRHHandledUserId" = ?
//       ) sub
//     WHERE
//       rn = 1
//   ),
//   actual_counts AS (
//     SELECT
//       ra."refAppointmentImpression" AS impression,
//       COUNT(*) AS count
//     FROM
//       latest_report_history rrh
//       JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
//     WHERE
//       ra."refAppointmentDate" >= ?  -- start_date parameter
//       AND ra."refAppointmentDate" <= ?  -- end_date parameter
//     GROUP BY
//       ra."refAppointmentImpression"
//   ),
//   expected_impressions AS (
//     SELECT
//       unnest(
//         array[
//           '1',
//           '1a',
//           '2',
//           '2a',
//           '3',
//           '3a',
//           '3b',
//           '3c',
//           '3d',
//           '3e',
//           '3f',
//           '3g',
//           '4',
//           '4a',
//           '4b',
//           '4c',
//           '4d',
//           '4e',
//           '4f',
//           '4g',
//           '4h',
//           '4i',
//           '4j',
//           '4k',
//           '4l',
//           '4m',
//           '5',
//           '6',
//           '6a',
//           '6b',
//           '6c',
//           '6d',
//           '6e',
//           '6f',
//           '6g',
//           '7a',
//           '7b',
//           '7c',
//           '7d',
//           '7e',
//           '10',
//           '10a'
//         ]
//       ) AS impression
//   )
// SELECT
//   ei.impression,
//   COALESCE(ac.count, 0) AS count
// FROM
//   expected_impressions ei
//   LEFT JOIN actual_counts ac ON ei.impression = ac.impression
// ORDER BY
//   ei.impression;
// `

// var TotalTATSQL = `
// SELECT
//   COUNT(*) FILTER (
//     WHERE
//       duration_days <= 1
//   ) AS le_1_day,
//   COUNT(*) FILTER (
//     WHERE
//       duration_days > 1
//       AND duration_days <= 3
//   ) AS le_3_days,
//   COUNT(*) FILTER (
//     WHERE
//       duration_days > 3
//       AND duration_days <= 7
//   ) AS le_7_days,
//   COUNT(*) FILTER (
//     WHERE
//       duration_days > 7
//       AND duration_days <= 10
//   ) AS le_10_days,
//   COUNT(*) FILTER (
//     WHERE
//       duration_days > 10
//   ) AS gt_10_days
// FROM
//   (
//     SELECT
//       rrh."refAppointmentId",
//       (
//         EXTRACT(
//           EPOCH
//           FROM
//             (
//               rrh."refRHHandleEndTime"::timestamp - rrh."refRHHandleStartTime"::timestamp
//             )
//         ) / 86400.0
//       ) AS duration_days
//     FROM
//       notes."refReportsHistory" rrh
//     WHERE
//       TO_CHAR(rrh."refRHHandleStartTime"::timestamp, 'YYYY-MM') = ?
//       AND rrh."refRHHandledUserId" = ?
//       AND rrh."refRHHandleStartTime" IS NOT NULL
//       AND rrh."refRHHandleEndTime" IS NOT NULL
//   ) AS subquery;
// `

// var LeftRecommendationScancenterSQL = `
// WITH
//   latest_report_history AS (
//     SELECT
//       *
//     FROM
//       (
//         SELECT
//           rrh.*,
//           ROW_NUMBER() OVER (
//             PARTITION BY
//               rrh."refAppointmentId"
//             ORDER BY
//               rrh."refRHHandleStartTime" DESC
//           ) AS rn
//         FROM
//           notes."refReportsHistory" rrh
//         WHERE
//           rrh."refRHHandleStatus" = 'Signed Off'
//           AND rrh."refRHHandleEndTime" IS NOT NULL
//       ) sub
//     WHERE
//       rn = 1
//   ),
//   groups AS (
//     SELECT DISTINCT
//       irc."refIRCName" AS group_name
//     FROM
//       impressionrecommendation."ImpressionRecommendationCategory" irc
//       JOIN impressionrecommendation."ImpressionRecommendationVal" irv ON irv."refIRCId" = irc."refIRCId"
//     WHERE
//       irv."refIRVSystemType" = 'WR'
//   ),
//   counts AS (
//     SELECT
//       irc."refIRCName" AS group_name,
//       COUNT(*) AS total_count
//     FROM
//       latest_report_history rrh
//       JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
//       JOIN public."Users" rrhu ON rrh."refRHHandledUserId" = rrhu."refUserId"
//       JOIN impressionrecommendation."ImpressionRecommendationVal" irv ON irv."refIRVCustId" = ra."refAppointmentRecommendation"
//       JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//     WHERE
//       irv."refIRVSystemType" = 'WR'
//       AND ra."refAppointmentStatus" = TRUE
//       AND rrh."refRHHandleEndTime" >= $1
//       AND rrh."refRHHandleEndTime" <= $2
//       AND (
//         $3 = 0
//         OR ra."refSCId" = $3
//       ) -- scan center filter
//       AND (
//         $4 = FALSE
//         OR rrhu."refRTId" IN (1, 6, 7, 10)
//       ) -- user role filter
//       AND rrh."refRHHandleStartTime" IS NOT NULL
//       AND rrh."refRHHandleEndTime" IS NOT NULL
//     GROUP BY
//       irc."refIRCName"
//   )
// SELECT
//   g.group_name,
//   COALESCE(c.total_count, 0) AS total_count
// FROM
//   groups g
//   LEFT JOIN counts c ON g.group_name = c.group_name
// ORDER BY
//   g.group_name;
// `

var LeftRecommendationScancenterSQL = `
WITH latest_signed AS (SELECT *
                       FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                         rrh."refRHHandleEndTime",
                                                                         rrh."refRHHandleStatus",
                                                                         ra."refAppointmentRecommendation",
                                                                         irv."refIRCId",
                                                                         irc."refIRCName"
                             FROM notes."refReportsHistory" rrh
                                      JOIN appointment."refAppointments" ra
                                           ON ra."refAppointmentId" = rrh."refAppointmentId"
                                      JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                           ON irv."refIRVCustId" = ra."refAppointmentRecommendation"
                                      JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                           ON irc."refIRCId" = irv."refIRCId"
                             WHERE ra."refAppointmentStatus" = TRUE
                               AND rrh."refRHHandleStatus" = 'Signed Off'
                               AND irv."refIRVSystemType" = 'WR'
                               AND (
                                 $3 = 0
                                     OR ra."refSCId" = $3
                                 )
                             ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                       WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                 BETWEEN $1::date AND $2::date)

SELECT irc."refIRCName"     AS group_name,
       COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
         LEFT JOIN latest_signed ls
                   ON ls."refIRCId" = irc."refIRCId"
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName";
`

var RightRecommendationScancenterSQL = `
WITH latest_signed AS (SELECT *
                       FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                         rrh."refRHHandleEndTime",
                                                                         rrh."refRHHandleStatus",
                                                                         ra."refAppointmentRecommendation",
                                                                         irv."refIRCId",
                                                                         irc."refIRCName"
                             FROM notes."refReportsHistory" rrh
                                      JOIN appointment."refAppointments" ra
                                           ON ra."refAppointmentId" = rrh."refAppointmentId"
                                      JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                           ON irv."refIRVCustId" = ra."refAppointmentRecommendationRight"
                                      JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                           ON irc."refIRCId" = irv."refIRCId"
                             WHERE ra."refAppointmentStatus" = TRUE
                               AND rrh."refRHHandleStatus" = 'Signed Off'
                               AND irv."refIRVSystemType" = 'WR'
                               AND (
                                 $3 = 0
                                     OR ra."refSCId" = $3
                                 )
                             ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                       WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                 BETWEEN $1::date AND $2::date)

SELECT irc."refIRCName"     AS group_name,
       COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
         LEFT JOIN latest_signed ls
                   ON ls."refIRCId" = irc."refIRCId"
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName";
`

var LeftRecommendationUserSQL = `
SELECT mainirc."refIRCName"                        AS group_name,
       (SELECT COUNT(irc."refIRCName")
        FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                          rrh."refRHHandleEndTime",
                                                          rrh."refRHHandleStatus",
                                                          rrh."refRHHandledUserId",
                                                          u."refRTId",
                                                          ra."refAppointmentRecommendation"
              FROM notes."refReportsHistory" rrh
                       JOIN appointment."refAppointments" ra
                            ON ra."refAppointmentId" = rrh."refAppointmentId"
                       JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
              WHERE ra."refAppointmentStatus" = TRUE
                AND rrh."refRHHandleStatus" = 'Signed Off'
                AND EXISTS (SELECT 1
                            FROM notes."refReportsHistory" rrhi
                                     JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                            WHERE (
                                ui."refRTId" NOT IN (2)
                                    OR
                                (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                )
                              AND (
                                ui."refRTId" NOT IN (8)
                                    OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                )
                              AND (
                                ui."refRTId" NOT IN (7)
                                    OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                )
                              AND (
                                ui."refRTId" NOT IN (6)
                                    OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                )
                              AND (
                                ui."refRTId" NOT IN (1, 5, 10)
                                    OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                )
                              AND rrhi."refRHHandledUserId" = ?
                              AND rrhi."refAppointmentId" = rrh."refAppointmentId")
              ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                 JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                      ON irv."refIRVCustId" = st."refAppointmentRecommendation"
                 JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN ?::date AND ?::date
          AND irv."refIRVSystemType" = 'WR'
          AND irc."refIRCId" = mainirc."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
`

var RightRecommendationUserSQL = `
SELECT mainirc."refIRCName"                        AS group_name,
       (SELECT COUNT(irc."refIRCName")
        FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                          rrh."refRHHandleEndTime",
                                                          rrh."refRHHandleStatus",
                                                          rrh."refRHHandledUserId",
                                                          u."refRTId",
                                                          ra."refAppointmentRecommendationRight"
              FROM notes."refReportsHistory" rrh
                       JOIN appointment."refAppointments" ra
                            ON ra."refAppointmentId" = rrh."refAppointmentId"
                       JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
              WHERE ra."refAppointmentStatus" = TRUE
                AND rrh."refRHHandleStatus" = 'Signed Off'
                AND EXISTS (SELECT 1
                            FROM notes."refReportsHistory" rrhi
                                     JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                            WHERE (
                                ui."refRTId" NOT IN (2)
                                    OR
                                (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                )
                              AND (
                                ui."refRTId" NOT IN (8)
                                    OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                )
                              AND (
                                ui."refRTId" NOT IN (7)
                                    OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                )
                              AND (
                                ui."refRTId" NOT IN (6)
                                    OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                )
                              AND (
                                ui."refRTId" NOT IN (1, 5, 10)
                                    OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                )
                              AND rrhi."refRHHandledUserId" = ?
                              AND rrhi."refAppointmentId" = rrh."refAppointmentId")
              ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                 JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                      ON irv."refIRVCustId" = st."refAppointmentRecommendationRight"
                 JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN ?::date AND ?::date
          AND irv."refIRVSystemType" = 'WR'
          AND irc."refIRCId" = mainirc."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
`

var TotalTATSQL = `
SELECT SUM(CASE WHEN diff_days <= 1 THEN 1 ELSE 0 END)                    AS le_1_day,
       SUM(CASE WHEN diff_days > 1 AND diff_days <= 3 THEN 1 ELSE 0 END)  AS le_3_days,
       SUM(CASE WHEN diff_days > 3 AND diff_days <= 7 THEN 1 ELSE 0 END)  AS le_7_days,
       SUM(CASE WHEN diff_days > 7 AND diff_days <= 10 THEN 1 ELSE 0 END) AS le_10_days,
       SUM(CASE WHEN diff_days > 10 THEN 1 ELSE 0 END)                    AS gt_10_days
FROM (SELECT *,
             EXTRACT(EPOCH FROM (
                 (t."refRHHandleEndTime"::timestamp) -
                 (t."refStartTime"::timestamp)
                 )) / 86400 AS diff_days
      FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                        rrh."refRHHandleStatus",
                                                        rrh."refRHHandledUserId",
                                                        u."refRTId",
                                                        ra."refAppointmentId",
                                                        ra."refCategoryId",
                                                        CASE
                                                            WHEN EXISTS (SELECT 1
                                                                         FROM notes."refReportsHistory" rh2
                                                                         WHERE rh2."refRHHandleStatus" = 'Technologist Form Fill'
                                                                           AND rh2."refAppointmentId" = rrh."refAppointmentId"
                                                                           AND rh2."refRHHandleEndTime" IS NOT NULL)
                                                                THEN (SELECT TO_CHAR(
                                                                                     rh2."refRHHandleEndTime"::timestamp,
                                                                                     'YYYY-MM-DD HH24:MI:SS')
                                                                      FROM notes."refReportsHistory" rh2
                                                                      WHERE rh2."refRHHandleStatus" = 'Technologist Form Fill'
                                                                        AND rh2."refAppointmentId" = rrh."refAppointmentId"
                                                                        AND rh2."refRHHandleEndTime" IS NOT NULL
                                                                      ORDER BY rh2."refRHId" DESC
                                                                      LIMIT 1)
                                                            ELSE (SELECT TO_CHAR(rdf."refDFCreatedAt"::timestamp, 'YYYY-MM-DD HH24:MI:SS')
                                                                  FROM dicom."refDicomFiles" rdf
                                                                  WHERE rdf."refAppointmentId" = rrh."refAppointmentId"
                                                                  ORDER BY rdf."refDFId" DESC
                                                                  LIMIT 1)
                                                            END AS "refStartTime",
                                                        rrh."refRHHandleEndTime"
            FROM notes."refReportsHistory" rrh
                     JOIN appointment."refAppointments" ra
                          ON ra."refAppointmentId" = rrh."refAppointmentId"
                     JOIN public."Users" u
                          ON u."refUserId" = rrh."refRHHandledUserId"
            WHERE ra."refAppointmentStatus" = TRUE
              AND rrh."refRHHandleStatus" = 'Signed Off'
              AND EXISTS (SELECT 1
                          FROM notes."refReportsHistory" rrhi
                                   JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                          WHERE (
                              ui."refRTId" NOT IN (2)
                                  OR
                              (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                              )
                            AND (
                              ui."refRTId" NOT IN (8)
                                  OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                              )
                            AND (
                              ui."refRTId" NOT IN (7)
                                  OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                              )
                            AND (
                              ui."refRTId" NOT IN (6)
                                  OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                              )
                            AND (
                              ui."refRTId" NOT IN (1, 5, 10)
                                  OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                              )
                            AND rrhi."refRHHandledUserId" = ?
                            AND rrhi."refAppointmentId" = rrh."refAppointmentId")
            ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
      WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
          BETWEEN ?::date AND ?::date
        AND t."refStartTime" IS NOT NULL
        AND t."refRHHandleEndTime" IS NOT NULL) x;
`

// var TechArtificatsAll = `
// SELECT
//
//	COUNT(CASE WHEN ra."refAppointmentTechArtifactsLeft" = TRUE
//	            AND ra."refAppointmentTechArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
//	COUNT(CASE WHEN ra."refAppointmentTechArtifactsLeft" = FALSE
//	            AND ra."refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS rightartifacts,
//	COUNT(CASE WHEN ra."refAppointmentTechArtifactsLeft" = TRUE
//	            AND ra."refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS bothartifacts
//
// FROM
//
//	appointment."refAppointments" ra
//
// WHERE
//
//	($1 = 0 OR ra."refSCId" = $1) AND
//	ra."refAppointmentStatus" = TRUE AND
//	ra."refAppointmentDate" >= $2 AND
//	ra."refAppointmentDate" <= $3;
//
// `
var TechArtificatsAll = `
SELECT COUNT(CASE
                 WHEN "refAppointmentTechArtifactsLeft" = TRUE
                     AND "refAppointmentTechArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
       COUNT(CASE
                 WHEN "refAppointmentTechArtifactsLeft" = FALSE
                     AND "refAppointmentTechArtifactsRight" = TRUE THEN 1 END)  AS rightartifacts,
       COUNT(CASE
                 WHEN "refAppointmentTechArtifactsLeft" = TRUE
                     AND "refAppointmentTechArtifactsRight" = TRUE THEN 1 END)  AS bothartifacts
FROM (SELECT *
      FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                        rrh."refRHHandleEndTime",
                                                        rrh."refRHHandleStatus",
                                                        ra."refCategoryId",
                                                        ra."refAppointmentTechArtifactsLeft",
                                                        ra."refAppointmentTechArtifactsRight"
            FROM notes."refReportsHistory" rrh
                     JOIN appointment."refAppointments" ra
                          ON ra."refAppointmentId" = rrh."refAppointmentId"
            WHERE ra."refAppointmentStatus" = TRUE
              AND rrh."refRHHandleStatus" = 'Signed Off'
              AND (
                $3 = 0
                    OR ra."refSCId" = $3
                )
            ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
      WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                BETWEEN $1::date
                AND $2::date) AS t;
`

var ReportArtificatsAll = `
SELECT COUNT(CASE
                 WHEN "refAppointmentReportArtifactsLeft" = TRUE
                     AND "refAppointmentReportArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
       COUNT(CASE
                 WHEN "refAppointmentReportArtifactsLeft" = FALSE
                     AND "refAppointmentReportArtifactsRight" = TRUE THEN 1 END)  AS rightartifacts,
       COUNT(CASE
                 WHEN "refAppointmentReportArtifactsLeft" = TRUE
                     AND "refAppointmentReportArtifactsRight" = TRUE THEN 1 END)  AS bothartifacts
FROM (SELECT *
      FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                        rrh."refRHHandleEndTime",
                                                        rrh."refRHHandleStatus",
                                                        ra."refCategoryId",
                                                        ra."refAppointmentReportArtifactsLeft",
                                                        ra."refAppointmentReportArtifactsRight"
            FROM notes."refReportsHistory" rrh
                     JOIN appointment."refAppointments" ra
                          ON ra."refAppointmentId" = rrh."refAppointmentId"
            WHERE ra."refAppointmentStatus" = TRUE
              AND rrh."refRHHandleStatus" = 'Signed Off'
              AND (
                $3 = 0
                    OR ra."refSCId" = $3
                )
            ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
      WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                BETWEEN $1::date
                AND $2::date) AS t;
`

var TechArtificats = `
SELECT COUNT(CASE
                 WHEN t."refAppointmentTechArtifactsLeft" = TRUE
                     AND t."refAppointmentTechArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
       COUNT(CASE
                 WHEN t."refAppointmentTechArtifactsLeft" = FALSE
                     AND t."refAppointmentTechArtifactsRight" = TRUE THEN 1 END)  AS rightartifacts,
       COUNT(CASE
                 WHEN t."refAppointmentTechArtifactsLeft" = TRUE
                     AND t."refAppointmentTechArtifactsRight" = TRUE THEN 1 END)  AS bothartifacts
FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                  rrh."refRHHandleStatus",
                                                  rrh."refRHHandledUserId",
                                                  u."refRTId",
                                                  ra."refAppointmentId",
                                                  ra."refCategoryId",
                                                  rrh."refRHHandleEndTime",
                                                  ra."refAppointmentTechArtifactsLeft",
                                                  ra."refAppointmentTechArtifactsRight"
      FROM notes."refReportsHistory" rrh
               JOIN appointment."refAppointments" ra
                    ON ra."refAppointmentId" = rrh."refAppointmentId"
               JOIN public."Users" u
                    ON u."refUserId" = rrh."refRHHandledUserId"
      WHERE ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND EXISTS (SELECT 1
                    FROM notes."refReportsHistory" rrhi
                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                    WHERE (
                        ui."refRTId" NOT IN (2)
                            OR
                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                        )
                      AND (
                        ui."refRTId" NOT IN (8)
                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                        )
                      AND (
                        ui."refRTId" NOT IN (7)
                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                        )
                      AND (
                        ui."refRTId" NOT IN (6)
                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                        )
                      AND (
                        ui."refRTId" NOT IN (1, 5, 10)
                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                        )
                      AND rrhi."refRHHandledUserId" = ?
                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
          BETWEEN ?::date AND ?::date
`

var ReportArtificats = `
SELECT COUNT(CASE
                 WHEN t."refAppointmentReportArtifactsLeft" = TRUE
                     AND t."refAppointmentReportArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
       COUNT(CASE
                 WHEN t."refAppointmentReportArtifactsLeft" = FALSE
                     AND t."refAppointmentReportArtifactsRight" = TRUE THEN 1 END)  AS rightartifacts,
       COUNT(CASE
                 WHEN t."refAppointmentReportArtifactsLeft" = TRUE
                     AND t."refAppointmentReportArtifactsRight" = TRUE THEN 1 END)  AS bothartifacts
FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                  rrh."refRHHandleStatus",
                                                  rrh."refRHHandledUserId",
                                                  u."refRTId",
                                                  ra."refAppointmentId",
                                                  ra."refCategoryId",
                                                  rrh."refRHHandleEndTime",
                                                  ra."refAppointmentReportArtifactsLeft",
                                                  ra."refAppointmentReportArtifactsRight"
      FROM notes."refReportsHistory" rrh
               JOIN appointment."refAppointments" ra
                    ON ra."refAppointmentId" = rrh."refAppointmentId"
               JOIN public."Users" u
                    ON u."refUserId" = rrh."refRHHandledUserId"
      WHERE ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND EXISTS (SELECT 1
                    FROM notes."refReportsHistory" rrhi
                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                    WHERE (
                        ui."refRTId" NOT IN (2)
                            OR
                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                        )
                      AND (
                        ui."refRTId" NOT IN (8)
                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                        )
                      AND (
                        ui."refRTId" NOT IN (7)
                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                        )
                      AND (
                        ui."refRTId" NOT IN (6)
                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                        )
                      AND (
                        ui."refRTId" NOT IN (1, 5, 10)
                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                        )
                      AND rrhi."refRHHandledUserId" = ?
                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
          BETWEEN ?::date AND ?::date
`

var TotalUserCountAnalaytics6monthSQL = `
WITH
  months AS (
    SELECT
      TO_CHAR(date_trunc('month', CURRENT_DATE) - INTERVAL '1 month' * n, 'YYYY-MM') AS month,
      TO_CHAR(date_trunc('month', CURRENT_DATE) - INTERVAL '1 month' * n, 'Month') AS month_name
    FROM
      generate_series(0, 5) AS n
  ),
  appointment_counts AS (
    SELECT
      TO_CHAR(TO_TIMESTAMP("refRHHandleStartTime", 'YYYY-MM-DD HH24:MI:SS'), 'YYYY-MM') AS month,
      COUNT(DISTINCT "refAppointmentId") AS total
    FROM
      notes."refReportsHistory"
    WHERE
      TO_TIMESTAMP("refRHHandleStartTime", 'YYYY-MM-DD HH24:MI:SS') >= date_trunc('month', CURRENT_DATE) - INTERVAL '6 months'
    GROUP BY
      month
  )
SELECT
  m.month,
  TRIM(m.month_name) AS month_name,
  COALESCE(a.total, 0) AS total_appointments
FROM
  months m
  LEFT JOIN appointment_counts a ON m.month = a.month
ORDER BY
  m.month;
`

var GetUsers6MonthTotalCountSQL = `
WITH months AS (SELECT date_trunc('month', CURRENT_DATE) - interval '1 month' * n AS month_start
                FROM generate_series(0, 5) AS n)
SELECT TO_CHAR(month_start, 'YYYY-MM')                                                                        AS month,
       TO_CHAR(month_start, 'Month')                                                                          AS month_name,
       month_start::date                                                                                      AS starting_date,
       (month_start + INTERVAL '1 month - 1 day')::date                                                       AS ending_date,
       (SELECT COUNT(*) AS total_rows
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND (
                        $1 = 0
                            OR ra."refSCId" = $1
                        )
                      AND rrh."refRHHandleStatus" IN
                          ('Signed Off')
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
              WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN month_start::date AND (month_start + INTERVAL '1 month - 1 day')::date) AS t) AS total_appointments
FROM months
ORDER BY month_start;
`

var TotoalUserAnalayticsSQL = `
SELECT uid."refUserId",
       uid."refUserCustId",
       (SELECT COUNT(*) AS total_appointments
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                rA."refAppointmentId",
                                                                rA."refCategoryId"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = uid."refUserId"
                                    AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
              WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t)                                      AS "totalcase",
       (SELECT COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END) AS "SForm"
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                rA."refAppointmentId",
                                                                rA."refCategoryId"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = uid."refUserId"
                                    AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
              WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t)                                      AS "totalsform",
       (SELECT COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END) AS "DaForm"
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                rA."refAppointmentId",
                                                                rA."refCategoryId"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = uid."refUserId"
                                    AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
              WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t)                                      AS "totaldaform",
       (SELECT COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END) AS "DbForm"
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                rA."refAppointmentId",
                                                                rA."refCategoryId"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = uid."refUserId"
                                    AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
              WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t)                                      AS "totaldbform",
       (SELECT COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END) AS "DcForm"
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                rA."refAppointmentId",
                                                                rA."refCategoryId"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = uid."refUserId"
                                    AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
              WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t)                                      AS "totaldcform",
       (SELECT COUNT(CASE
                         WHEN t."refAppointmentTechArtifactsLeft" = TRUE THEN 1 END) AS leftartifacts
        FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                          rrh."refRHHandleStatus",
                                                          rrh."refRHHandledUserId",
                                                          u."refRTId",
                                                          ra."refAppointmentId",
                                                          ra."refCategoryId",
                                                          rrh."refRHHandleEndTime",
                                                          ra."refAppointmentTechArtifactsLeft",
                                                          ra."refAppointmentTechArtifactsRight"
              FROM notes."refReportsHistory" rrh
                       JOIN appointment."refAppointments" ra
                            ON ra."refAppointmentId" = rrh."refAppointmentId"
                       JOIN public."Users" u
                            ON u."refUserId" = rrh."refRHHandledUserId"
              WHERE ra."refAppointmentStatus" = TRUE
                AND rrh."refRHHandleStatus" = 'Signed Off'
                AND EXISTS (SELECT 1
                            FROM notes."refReportsHistory" rrhi
                                     JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                            WHERE (
                                ui."refRTId" NOT IN (2)
                                    OR
                                (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                )
                              AND (
                                ui."refRTId" NOT IN (8)
                                    OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                )
                              AND (
                                ui."refRTId" NOT IN (7)
                                    OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                )
                              AND (
                                ui."refRTId" NOT IN (6)
                                    OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                )
                              AND (
                                ui."refRTId" NOT IN (1, 5, 10)
                                    OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                )
                              AND rrhi."refRHHandledUserId" = uid."refUserId"
                              AND rrhi."refAppointmentId" = rrh."refAppointmentId")
              ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
        WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                  BETWEEN $1::date AND $2::date)                                                  AS "techartificatsleft",
       (SELECT COUNT(CASE
                         WHEN t."refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS leftartifacts
        FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                          rrh."refRHHandleStatus",
                                                          rrh."refRHHandledUserId",
                                                          u."refRTId",
                                                          ra."refAppointmentId",
                                                          ra."refCategoryId",
                                                          rrh."refRHHandleEndTime",
                                                          ra."refAppointmentTechArtifactsLeft",
                                                          ra."refAppointmentTechArtifactsRight"
              FROM notes."refReportsHistory" rrh
                       JOIN appointment."refAppointments" ra
                            ON ra."refAppointmentId" = rrh."refAppointmentId"
                       JOIN public."Users" u
                            ON u."refUserId" = rrh."refRHHandledUserId"
              WHERE ra."refAppointmentStatus" = TRUE
                AND rrh."refRHHandleStatus" = 'Signed Off'
                AND EXISTS (SELECT 1
                            FROM notes."refReportsHistory" rrhi
                                     JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                            WHERE (
                                ui."refRTId" NOT IN (2)
                                    OR
                                (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                )
                              AND (
                                ui."refRTId" NOT IN (8)
                                    OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                )
                              AND (
                                ui."refRTId" NOT IN (7)
                                    OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                )
                              AND (
                                ui."refRTId" NOT IN (6)
                                    OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                )
                              AND (
                                ui."refRTId" NOT IN (1, 5, 10)
                                    OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                )
                              AND rrhi."refRHHandledUserId" = uid."refUserId"
                              AND rrhi."refAppointmentId" = rrh."refAppointmentId")
              ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
        WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                  BETWEEN $1::date AND $2::date)                                                  AS "techartificatsright",
       (SELECT COUNT(CASE
                         WHEN t."refAppointmentReportArtifactsLeft" = TRUE THEN 1 END) AS leftartifacts
        FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                          rrh."refRHHandleStatus",
                                                          rrh."refRHHandledUserId",
                                                          u."refRTId",
                                                          ra."refAppointmentId",
                                                          ra."refCategoryId",
                                                          rrh."refRHHandleEndTime",
                                                          ra."refAppointmentReportArtifactsLeft",
                                                          ra."refAppointmentReportArtifactsRight"
              FROM notes."refReportsHistory" rrh
                       JOIN appointment."refAppointments" ra
                            ON ra."refAppointmentId" = rrh."refAppointmentId"
                       JOIN public."Users" u
                            ON u."refUserId" = rrh."refRHHandledUserId"
              WHERE ra."refAppointmentStatus" = TRUE
                AND rrh."refRHHandleStatus" = 'Signed Off'
                AND EXISTS (SELECT 1
                            FROM notes."refReportsHistory" rrhi
                                     JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                            WHERE (
                                ui."refRTId" NOT IN (2)
                                    OR
                                (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                )
                              AND (
                                ui."refRTId" NOT IN (8)
                                    OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                )
                              AND (
                                ui."refRTId" NOT IN (7)
                                    OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                )
                              AND (
                                ui."refRTId" NOT IN (6)
                                    OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                )
                              AND (
                                ui."refRTId" NOT IN (1, 5, 10)
                                    OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                )
                              AND rrhi."refRHHandledUserId" = uid."refUserId"
                              AND rrhi."refAppointmentId" = rrh."refAppointmentId")
              ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
        WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                  BETWEEN $1::date AND $2::date)                                                  AS "reportartificatsleft",
       (SELECT COUNT(CASE
                         WHEN t."refAppointmentReportArtifactsRight" = TRUE THEN 1 END) AS leftartifacts
        FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                          rrh."refRHHandleStatus",
                                                          rrh."refRHHandledUserId",
                                                          u."refRTId",
                                                          ra."refAppointmentId",
                                                          ra."refCategoryId",
                                                          rrh."refRHHandleEndTime",
                                                          ra."refAppointmentReportArtifactsLeft",
                                                          ra."refAppointmentReportArtifactsRight"
              FROM notes."refReportsHistory" rrh
                       JOIN appointment."refAppointments" ra
                            ON ra."refAppointmentId" = rrh."refAppointmentId"
                       JOIN public."Users" u
                            ON u."refUserId" = rrh."refRHHandledUserId"
              WHERE ra."refAppointmentStatus" = TRUE
                AND rrh."refRHHandleStatus" = 'Signed Off'
                AND EXISTS (SELECT 1
                            FROM notes."refReportsHistory" rrhi
                                     JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                            WHERE (
                                ui."refRTId" NOT IN (2)
                                    OR
                                (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                )
                              AND (
                                ui."refRTId" NOT IN (8)
                                    OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                )
                              AND (
                                ui."refRTId" NOT IN (7)
                                    OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                )
                              AND (
                                ui."refRTId" NOT IN (6)
                                    OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                )
                              AND (
                                ui."refRTId" NOT IN (1, 5, 10)
                                    OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                )
                              AND rrhi."refRHHandledUserId" = uid."refUserId"
                              AND rrhi."refAppointmentId" = rrh."refAppointmentId")
              ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
        WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                  BETWEEN $1::date AND $2::date)                                                  AS "reportartificatsright",
       (SELECT SUM(
                       EXTRACT(EPOCH FROM (sub."refRHHandleEndTime"::timestamp - sub."refStartTime"::timestamp)) / 3600
               ) AS total_hours
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                CASE
                                                                    WHEN EXISTS (SELECT 1
                                                                                 FROM notes."refReportsHistory" rh2
                                                                                 WHERE rh2."refRHHandleStatus" = 'Technologist Form Fill'
                                                                                   AND rh2."refAppointmentId" = rrh."refAppointmentId"
                                                                                   AND rh2."refRHHandleEndTime" IS NOT NULL)
                                                                        THEN (SELECT CAST(rh2."refRHHandleEndTime" AS timestamp)
                                                                              FROM notes."refReportsHistory" rh2
                                                                              WHERE rh2."refRHHandleStatus" = 'Technologist Form Fill'
                                                                                AND rh2."refAppointmentId" = rrh."refAppointmentId"
                                                                                AND rh2."refRHHandleEndTime" IS NOT NULL
                                                                              ORDER BY rh2."refRHId" DESC
                                                                              LIMIT 1)
                                                                    ELSE (SELECT CAST(rdf."refDFCreatedAt" AS timestamp)
                                                                          FROM dicom."refDicomFiles" rdf
                                                                          WHERE rdf."refAppointmentId" = rrh."refAppointmentId"
                                                                          ORDER BY rdf."refDFId" DESC
                                                                          LIMIT 1)
                                                                    END AS "refStartTime"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = uid."refUserId"
                                    AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
              WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) sub
        WHERE sub."refStartTime" IS NOT NULL
          AND sub."refRHHandleEndTime" IS NOT NULL
        GROUP BY sub."refRHHandledUserId")                                                        AS "totaltiming",
       (SELECT COALESCE(SUM(COALESCE(sub."totalCorrect", 0)), 0) AS "totalCorrect"
        FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                          rrh."refRHHandleEndTime",
                                                          rrh."refRHHandleStatus",
                                                          rrh."refRHHandledUserId",
                                                          u."refRTId",
                                                          rrh."refRHHandleCorrect",
                                                          rrh."refRHHandleEdit",
                                                          (SELECT newCrt."refRHHandleCorrect"
                                                           FROM (SELECT DISTINCT ON (rrhi."refAppointmentId") *
                                                                 FROM notes."refReportsHistory" rrhi
                                                                          JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                                                 WHERE (
                                                                     ui."refRTId" NOT IN (2)
                                                                         OR
                                                                     (ui."refRTId" IN (2) AND
                                                                      rrhi."refRHHandleStatus" =
                                                                      'Technologist Form Fill')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (8)
                                                                         OR
                                                                     (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (7)
                                                                         OR
                                                                     (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (6)
                                                                         OR
                                                                     (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (1, 5, 10)
                                                                         OR
                                                                     (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                                                     )
                                                                   AND rrhi."refRHHandledUserId" = uid."refUserId"
                                                                 ORDER BY rrhi."refAppointmentId", rrhi."refRHId"
                                                                     DESC) newCrt
                                                           WHERE newCrt."refAppointmentId" = rrh."refAppointmentId") AS "totalCorrect",
                                                          (SELECT newCrt."refRHHandleEdit"
                                                           FROM (SELECT DISTINCT ON (rrhi."refAppointmentId") *
                                                                 FROM notes."refReportsHistory" rrhi
                                                                          JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                                                 WHERE (
                                                                     ui."refRTId" NOT IN (2)
                                                                         OR
                                                                     (ui."refRTId" IN (2) AND
                                                                      rrhi."refRHHandleStatus" =
                                                                      'Technologist Form Fill')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (8)
                                                                         OR
                                                                     (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (7)
                                                                         OR
                                                                     (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (6)
                                                                         OR
                                                                     (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (1, 5, 10)
                                                                         OR
                                                                     (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                                                     )
                                                                   AND rrhi."refRHHandledUserId" = uid."refUserId"
                                                                 ORDER BY rrhi."refAppointmentId", rrhi."refRHId"
                                                                     DESC) newCrt
                                                           WHERE newCrt."refAppointmentId" = rrh."refAppointmentId") AS "totalEdit"
              FROM notes."refReportsHistory" rrh
                       JOIN appointment."refAppointments" ra
                            ON ra."refAppointmentId" = rrh."refAppointmentId"
                       JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
              WHERE ra."refAppointmentStatus" = TRUE
                AND rrh."refRHHandleStatus" = 'Signed Off'
                AND EXISTS (SELECT 1
                            FROM notes."refReportsHistory" rrhi
                                     JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                            WHERE (
                                ui."refRTId" NOT IN (2)
                                    OR
                                (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                )
                              AND (
                                ui."refRTId" NOT IN (8)
                                    OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                )
                              AND (
                                ui."refRTId" NOT IN (7)
                                    OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                )
                              AND (
                                ui."refRTId" NOT IN (6)
                                    OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                )
                              AND (
                                ui."refRTId" NOT IN (1, 5, 10)
                                    OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                )
                              AND rrhi."refRHHandledUserId" = uid."refUserId"
                              AND rrhi."refAppointmentId" = rrh."refAppointmentId")
              ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) sub
        WHERE TO_TIMESTAMP(sub."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                  BETWEEN $1::date AND $2::date)                                                  AS "totalreportcorrect",
       (SELECT COALESCE(SUM(COALESCE(sub."totalEdit", 0)), 0) AS "totalEdit"
        FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                          rrh."refRHHandleEndTime",
                                                          rrh."refRHHandleStatus",
                                                          rrh."refRHHandledUserId",
                                                          u."refRTId",
                                                          rrh."refRHHandleCorrect",
                                                          rrh."refRHHandleEdit",
                                                          (SELECT newCrt."refRHHandleCorrect"
                                                           FROM (SELECT DISTINCT ON (rrhi."refAppointmentId") *
                                                                 FROM notes."refReportsHistory" rrhi
                                                                          JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                                                 WHERE (
                                                                     ui."refRTId" NOT IN (2)
                                                                         OR
                                                                     (ui."refRTId" IN (2) AND
                                                                      rrhi."refRHHandleStatus" =
                                                                      'Technologist Form Fill')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (8)
                                                                         OR
                                                                     (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (7)
                                                                         OR
                                                                     (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (6)
                                                                         OR
                                                                     (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                                                     )
                                                                   AND (
                                                             ui."refRTId" NOT IN (10)
                                                                 OR
                                                             (ui."refRTId" IN (10) AND rrhi."refRHHandleStatus" = 'Reviewed 1')
                                                             )
                                                           AND (
                                                             ui."refRTId" NOT IN (1, 5)
                                                                 OR
                                                             (ui."refRTId" IN (1, 5) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                                             )
                                                                   AND rrhi."refRHHandledUserId" = uid."refUserId"
                                                                 ORDER BY rrhi."refAppointmentId", rrhi."refRHId"
                                                                     DESC) newCrt
                                                           WHERE newCrt."refAppointmentId" = rrh."refAppointmentId") AS "totalCorrect",
                                                          (SELECT newCrt."refRHHandleEdit"
                                                           FROM (SELECT DISTINCT ON (rrhi."refAppointmentId") *
                                                                 FROM notes."refReportsHistory" rrhi
                                                                          JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                                                 WHERE (
                                                                     ui."refRTId" NOT IN (2)
                                                                         OR
                                                                     (ui."refRTId" IN (2) AND
                                                                      rrhi."refRHHandleStatus" =
                                                                      'Technologist Form Fill')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (8)
                                                                         OR
                                                                     (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (7)
                                                                         OR
                                                                     (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                                                     )
                                                                   AND (
                                                                     ui."refRTId" NOT IN (6)
                                                                         OR
                                                                     (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                                                     )
                                                                   AND (
                                                             ui."refRTId" NOT IN (10)
                                                                 OR
                                                             (ui."refRTId" IN (10) AND rrhi."refRHHandleStatus" = 'Reviewed 1')
                                                             )
                                                           AND (
                                                             ui."refRTId" NOT IN (1, 5)
                                                                 OR
                                                             (ui."refRTId" IN (1, 5) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                                             )
                                                                   AND rrhi."refRHHandledUserId" = uid."refUserId"
                                                                 ORDER BY rrhi."refAppointmentId", rrhi."refRHId"
                                                                     DESC) newCrt
                                                           WHERE newCrt."refAppointmentId" = rrh."refAppointmentId") AS "totalEdit"
              FROM notes."refReportsHistory" rrh
                       JOIN appointment."refAppointments" ra
                            ON ra."refAppointmentId" = rrh."refAppointmentId"
                       JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
              WHERE ra."refAppointmentStatus" = TRUE
                AND rrh."refRHHandleStatus" = 'Signed Off'
                AND EXISTS (SELECT 1
                            FROM notes."refReportsHistory" rrhi
                                     JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                            WHERE (
                                ui."refRTId" NOT IN (2)
                                    OR
                                (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                )
                              AND (
                                ui."refRTId" NOT IN (8)
                                    OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                )
                              AND (
                                ui."refRTId" NOT IN (7)
                                    OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                )
                              AND (
                                ui."refRTId" NOT IN (6)
                                    OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                )
                              AND (
                                                             ui."refRTId" NOT IN (10)
                                                                 OR
                                                             (ui."refRTId" IN (10) AND rrhi."refRHHandleStatus" = 'Reviewed 1')
                                                             )
                                                           AND (
                                                             ui."refRTId" NOT IN (1, 5)
                                                                 OR
                                                             (ui."refRTId" IN (1, 5) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                                             )
                              AND rrhi."refRHHandledUserId" = uid."refUserId"
                              AND rrhi."refAppointmentId" = rrh."refAppointmentId")
              ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) sub
        WHERE TO_TIMESTAMP(sub."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                  BETWEEN $1::date AND $2::date)                                                  AS "totalreportedit",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendation"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendation"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 1)                                                             AS "leftannualscreening",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendation"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendation"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 2)                                                             AS "leftusgsfu",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendation"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendation"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 3)                                                             AS "leftbiopsy",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendation"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendation"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 4)                                                             AS "leftbreastradiologist",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendation"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendation"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 5)                                                             AS "leftclinicalcorrelation",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendation"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendation"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 6)                                                             AS "leftoncoconsult",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendation"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND (
                          u."refRTId" NOT IN (2)
                              OR (u."refRTId" IN (2) AND rrh."refRHHandleStatus" = 'Technologist Form Fill')
                          )
                        AND (
                          u."refRTId" NOT IN (8)
                              OR (u."refRTId" IN (8) AND rrh."refRHHandleStatus" = 'Reviewed 2')
                          )
                        AND (
                          u."refRTId" NOT IN (7)
                              OR (u."refRTId" IN (7) AND rrh."refRHHandleStatus" = 'Predraft')
                          )
                        AND (
                          u."refRTId" NOT IN (6)
                              OR (u."refRTId" IN (6) AND rrh."refRHHandleStatus" = 'Draft')
                          )
                        AND (
                          u."refRTId" NOT IN (1, 5, 10)
                              OR (u."refRTId" IN (1, 5, 10) AND rrh."refRHHandleStatus" = 'Signed Off')
                          )
                        AND rrh."refRHHandledUserId" = uid."refUserId"
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendation"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 7)                                                             AS "leftredo",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendationRight"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendationRight"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 1)                                                             AS "rightannualscreening",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendationRight"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendationRight"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 2)                                                             AS "rightusgsfu",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendationRight"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendationRight"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 3)                                                             AS "rightbiopsy",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendationRight"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendationRight"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 4)                                                             AS "rightbreastradiologist",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendationRight"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendationRight"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 5)                                                             AS "rightclinicalcorrelation",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendationRight"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendationRight"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 6)                                                             AS "rightoncoconsult",
       (SELECT (SELECT COUNT(irc."refIRCName")
                FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                  rrh."refRHHandleEndTime",
                                                                  rrh."refRHHandleStatus",
                                                                  rrh."refRHHandledUserId",
                                                                  u."refRTId",
                                                                  ra."refAppointmentRecommendationRight"
                      FROM notes."refReportsHistory" rrh
                               JOIN appointment."refAppointments" ra
                                    ON ra."refAppointmentId" = rrh."refAppointmentId"
                               JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                      WHERE ra."refAppointmentStatus" = TRUE
                        AND rrh."refRHHandleStatus" = 'Signed Off'
                        AND EXISTS (SELECT 1
                                    FROM notes."refReportsHistory" rrhi
                                             JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                    WHERE (
                                        ui."refRTId" NOT IN (2)
                                            OR
                                        (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (8)
                                            OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (7)
                                            OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (6)
                                            OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                        )
                                      AND (
                                        ui."refRTId" NOT IN (1, 5, 10)
                                            OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                        )
                                      AND rrhi."refRHHandledUserId" = uid."refUserId"
                                      AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) st
                         JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                              ON irv."refIRVCustId" = st."refAppointmentRecommendationRight"
                         JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                              ON irc."refIRCId" = irv."refIRCId"
                WHERE TO_TIMESTAMP(st."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                    BETWEEN $1::date AND $2::date
                  AND irv."refIRVSystemType" = 'WR'
                  AND irc."refIRCId" = mainirc."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" mainirc
        WHERE mainirc."refIRCId" = 7)                                                             AS "rightredo",
       (SELECT COUNT(*)
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                rA."refAppointmentId",
                                                                rA."refCategoryId",
                                                                ra."refAppointmentDicomSide"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = uid."refUserId"
                                    AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
              WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) t
        WHERE t."refAppointmentDicomSide" = 'unilateralleft')                                     AS "unilateralleft",
       (SELECT COUNT(*)
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                rA."refAppointmentId",
                                                                rA."refCategoryId",
                                                                ra."refAppointmentDicomSide"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = uid."refUserId"
                                    AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
              WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) t
        WHERE t."refAppointmentDicomSide" = 'unilateralright')                                    AS "unilateralright",
       (SELECT COUNT(*)
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                rrh."refRHHandledUserId",
                                                                u."refRTId",
                                                                rA."refAppointmentId",
                                                                rA."refCategoryId",
                                                                ra."refAppointmentDicomSide"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                             JOIN public."Users" U on u."refUserId" = rrh."refRHHandledUserId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND EXISTS (SELECT 1
                                  FROM notes."refReportsHistory" rrhi
                                           JOIN public."Users" ui ON ui."refUserId" = rrhi."refRHHandledUserId"
                                  WHERE (
                                      ui."refRTId" NOT IN (2)
                                          OR
                                      (ui."refRTId" IN (2) AND rrhi."refRHHandleStatus" = 'Technologist Form Fill')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (8)
                                          OR (ui."refRTId" IN (8) AND rrhi."refRHHandleStatus" = 'Reviewed 2')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (7)
                                          OR (ui."refRTId" IN (7) AND rrhi."refRHHandleStatus" = 'Predraft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (6)
                                          OR (ui."refRTId" IN (6) AND rrhi."refRHHandleStatus" = 'Draft')
                                      )
                                    AND (
                                      ui."refRTId" NOT IN (1, 5, 10)
                                          OR (ui."refRTId" IN (1, 5, 10) AND rrhi."refRHHandleStatus" = 'Signed Off')
                                      )
                                    AND rrhi."refRHHandledUserId" = uid."refUserId"
                                    AND rrhi."refAppointmentId" = rrh."refAppointmentId")
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) t
              WHERE TO_TIMESTAMP(t."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) t
        WHERE (t."refAppointmentDicomSide" = 'bilateral' OR t."refAppointmentDicomSide" IS NULL)) AS "bilateral"
FROM public."Users" uid
         FULL JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = uid."refUserId"
WHERE uid."refRTId" NOT IN (3, 4, 9)
  AND (
    $3 = 0
        OR rscm."refSCId" = $3
    )
ORDER BY uid."refRTId";
`

// var GetOverAllScanCenterList = `
// SELECT
//   sc."refSCId",
//   sc."refSCCustId",
//   (
//     SELECT
//       COUNT(DISTINCT ra."refAppointmentId")
//     FROM
//       appointment."refAppointments" ra
//     WHERE
//       ra."refAppointmentDate"::timestamp >= $1
//       AND ra."refAppointmentDate"::timestamp <= $2
//       AND ra."refSCId" = sc."refSCId"
//       AND ra."refAppointmentStatus" = TRUE
//   ) AS "totalcase",
//   (
//     SELECT
//       COUNT(
//         DISTINCT CASE
//           WHEN ra."refCategoryId" = 1 THEN ra."refAppointmentId"
//         END
//       ) AS category_count
//     FROM
//       appointment."refAppointments" ra
//     WHERE
//       ra."refAppointmentDate"::timestamp >= $1
//       AND ra."refAppointmentDate"::timestamp <= $2
//       AND ra."refSCId" = sc."refSCId"
//       AND ra."refAppointmentStatus" = TRUE
//   ) AS "totalsform",
//   (
//     SELECT
//       COUNT(
//         DISTINCT CASE
//           WHEN ra."refCategoryId" = 2 THEN ra."refAppointmentId"
//         END
//       ) AS category_count
//     FROM
//       appointment."refAppointments" ra
//     WHERE
//       ra."refAppointmentDate"::timestamp >= $1
//       AND ra."refAppointmentDate"::timestamp <= $2
//       AND ra."refSCId" = sc."refSCId"
//       AND ra."refAppointmentStatus" = TRUE
//   ) AS "totaldaform",
//   (
//     SELECT
//       COUNT(
//         DISTINCT CASE
//           WHEN ra."refCategoryId" = 3 THEN ra."refAppointmentId"
//         END
//       ) AS category_count
//     FROM
//       appointment."refAppointments" ra
//     WHERE
//       ra."refAppointmentDate"::timestamp >= $1
//       AND ra."refAppointmentDate"::timestamp <= $2
//       AND ra."refSCId" = sc."refSCId"
//       AND ra."refAppointmentStatus" = TRUE
//   ) AS "totaldbform",
//   (
//     SELECT
//       COUNT(
//         DISTINCT CASE
//           WHEN ra."refCategoryId" = 4 THEN ra."refAppointmentId"
//         END
//       ) AS category_count
//     FROM
//       appointment."refAppointments" ra
//     WHERE
//       ra."refAppointmentDate"::timestamp >= $1
//       AND ra."refAppointmentDate"::timestamp <= $2
//       AND ra."refSCId" = sc."refSCId"
//       AND ra."refAppointmentStatus" = TRUE
//   ) AS "totaldcform",
//   (
//     SELECT
//       COUNT(
//         DISTINCT CASE
//           WHEN ra."refAppointmentTechArtifactsLeft" = true THEN ra."refAppointmentId"
//         END
//       ) AS category_count
//     FROM
//       appointment."refAppointments" ra
//     WHERE
//       ra."refAppointmentDate"::timestamp >= $1
//       AND ra."refAppointmentDate"::timestamp <= $2
//       AND ra."refSCId" = sc."refSCId"
//       AND ra."refAppointmentStatus" = TRUE
//   ) AS "techartificatsleft",
//   (
//     SELECT
//       COUNT(
//         DISTINCT CASE
//           WHEN ra."refAppointmentTechArtifactsRight" = true THEN ra."refAppointmentId"
//         END
//       ) AS category_count
//     FROM
//       appointment."refAppointments" ra
//     WHERE
//       ra."refAppointmentDate"::timestamp >= $1
//       AND ra."refAppointmentDate"::timestamp <= $2
//       AND ra."refSCId" = sc."refSCId"
//       AND ra."refAppointmentStatus" = TRUE
//   ) AS "techartificatsright",
//   (
//     SELECT
//       COUNT(
//         DISTINCT CASE
//           WHEN ra."refAppointmentReportArtifactsLeft" = true THEN ra."refAppointmentId"
//         END
//       ) AS category_count
//     FROM
//       appointment."refAppointments" ra
//     WHERE
//       ra."refAppointmentDate"::timestamp >= $1
//       AND ra."refAppointmentDate"::timestamp <= $2
//       AND ra."refSCId" = sc."refSCId"
//       AND ra."refAppointmentStatus" = TRUE
//   ) AS "reportartificatsleft",
//   (
//     SELECT
//       COUNT(
//         DISTINCT CASE
//           WHEN ra."refAppointmentReportArtifactsRight" = true THEN ra."refAppointmentId"
//         END
//       ) AS category_count
//     FROM
//       appointment."refAppointments" ra
//     WHERE
//       ra."refAppointmentDate"::timestamp >= $1
//       AND ra."refAppointmentDate"::timestamp <= $2
//       AND ra."refSCId" = sc."refSCId"
//       AND ra."refAppointmentStatus" = TRUE
//   ) AS "reportartificatsright",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendation" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Annual Screening' -- 👈 Replace this with your desired category
//       )
//   ) AS "leftannualscreening",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendation" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'USG/SFU' -- 👈 Replace this with your desired category
//       )
//   ) AS "leftusgsfu",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendation" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Biopsy' -- 👈 Replace this with your desired category
//       )
//   ) AS "leftBiopsy",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendation" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Breast Radiologist' -- 👈 Replace this with your desired category
//       )
//   ) AS "leftBreastradiologist",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendation" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Clinical Correlation' -- 👈 Replace this with your desired category
//       )
//   ) AS "leftClinicalCorrelation",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendation" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Onco Consult' -- 👈 Replace this with your desired category
//       )
//   ) AS "leftOncoConsult",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendation" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Redo' -- 👈 Replace this with your desired category
//       )
//   ) AS "leftRedo",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendationRight" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Annual Screening' -- 👈 Replace this with your desired category
//       )
//   ) AS "rightannualscreening",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendationRight" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'USG/SFU' -- 👈 Replace this with your desired category
//       )
//   ) AS "rightusgsfu",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendationRight" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Biopsy' -- 👈 Replace this with your desired category
//       )
//   ) AS "rightBiopsy",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendationRight" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Breast Radiologist' -- 👈 Replace this with your desired category
//       )
//   ) AS "rightBreastradiologist",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendationRight" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Clinical Correlation' -- 👈 Replace this with your desired category
//       )
//   ) AS "rightClinicalCorrelation",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendationRight" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Onco Consult' -- 👈 Replace this with your desired category
//       )
//   ) AS "rightOncoConsult",
//   (
//     SELECT
//       COUNT(*)
//     FROM
//       (
//         SELECT DISTINCT
//           ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
//         FROM
//           appointment."refAppointments" ra
//         WHERE
//           ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentDate"::timestamp >= $1
//           AND ra."refAppointmentDate"::timestamp <= $2
//           AND ra."refSCId" = sc."refSCId"
//           AND ra."refAppointmentStatus" = TRUE
//         ORDER BY
//           ra."refAppointmentId"
//       ) t
//     WHERE
//       t."refAppointmentRecommendationRight" IN (
//         SELECT
//           irv."refIRVCustId"
//         FROM
//           impressionrecommendation."ImpressionRecommendationVal" irv
//           JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
//         WHERE
//           irv."refIRVSystemType" = 'WR'
//           AND irc."refIRCName" = 'Redo' -- 👈 Replace this with your desired category
//       )
//   ) AS "rightRedo"
// FROM
//   public."ScanCenter" sc
// WHERE
//   (
//     $3 = 0
//     OR sc."refSCId" = $3
//   )
// `

var GetOverAllScanCenterList = `
SELECT sc."refSCId",
       sc."refSCCustId",
       (SELECT COUNT(*) AS total_rows
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND ra."refSCId" = sc."refSCId"
                      AND rrh."refRHHandleStatus" = 'Signed Off'

                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t) AS totalcase,
       (SELECT COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END) AS "SForm"
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refCategoryId" -- ✅ REQUIRED for outer COUNT()
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND ra."refSCId" = sc."refSCId"
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t) AS totalsform,
       (SELECT COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END) AS "SForm"
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refCategoryId" -- ✅ REQUIRED for outer COUNT()
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND ra."refSCId" = sc."refSCId"
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t) AS totaldaform,
       (SELECT COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END) AS "SForm"
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refCategoryId" -- ✅ REQUIRED for outer COUNT()
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND ra."refSCId" = sc."refSCId"
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t) AS totaldbform,
       (SELECT COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END) AS "SForm"
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refCategoryId" -- ✅ REQUIRED for outer COUNT()
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND ra."refSCId" = sc."refSCId"
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t) AS totaldcform,
       (SELECT COUNT(CASE WHEN "refAppointmentTechArtifactsLeft" = TRUE THEN 1 END) AS leftartifacts
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refCategoryId",
                                                                ra."refAppointmentTechArtifactsLeft",
                                                                ra."refAppointmentTechArtifactsRight"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND ra."refSCId" = sc."refSCId"
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t) AS techartificatsleft,
       (SELECT COUNT(CASE WHEN "refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS rightartifacts
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refCategoryId",
                                                                ra."refAppointmentTechArtifactsLeft",
                                                                ra."refAppointmentTechArtifactsRight"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND ra."refSCId" = sc."refSCId"
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t) AS techartificatsright,
       (SELECT COUNT(CASE WHEN "refAppointmentReportArtifactsLeft" = TRUE THEN 1 END) AS leftartifacts
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refCategoryId",
                                                                ra."refAppointmentReportArtifactsLeft",
                                                                ra."refAppointmentReportArtifactsRight"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND ra."refSCId" = sc."refSCId"
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t) AS reportartificatsleft,
       (SELECT COUNT(CASE WHEN "refAppointmentReportArtifactsRight" = TRUE THEN 1 END) AS rightartifacts
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refCategoryId",
                                                                ra."refAppointmentReportArtifactsLeft",
                                                                ra."refAppointmentReportArtifactsRight"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                      AND ra."refSCId" = sc."refSCId"
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t) AS reportartificatsright,
       --Left Annual Screening
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendation",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendation"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 1
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS leftannualscreening,
       -- Left USG/SFU
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendation",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendation"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 2
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS leftusgsfu,
       -- Left Biopsy
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendation",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendation"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 3
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS leftBiopsy,
       --Left Breast Radiologist
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendation",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendation"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 4
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS leftBreastradiologist,
       --Left Clinical Correlation
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendation",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendation"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 5
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS leftClinicalCorrelation,
       --Left Onco Consult
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendation",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendation"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 6
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS leftOncoConsult,
       --Left Redo
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendation",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendation"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 7
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS leftRedo,
       --Right Annual Screening
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendationRight",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendationRight"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 1
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS rightannualscreening,
       --Right USG/SFU
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendationRight",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendationRight"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 2
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS rightusgsfu,
       --Right Biopsy
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendationRight",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendationRight"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 3
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS rightBiopsy,
       --Right Breast Radiologist
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendationRight",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendationRight"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 4
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS rightBreastradiologist,
       --Right Clinical Correlation
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendationRight",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendationRight"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 5
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS rightClinicalCorrelation,
       --Right Onco Consult
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendationRight",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendationRight"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 6
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS rightOncoConsult,
       --Right Redo
       (WITH latest_signed AS (SELECT *
                               FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                                 rrh."refRHHandleEndTime",
                                                                                 rrh."refRHHandleStatus",
                                                                                 ra."refAppointmentRecommendationRight",
                                                                                 irv."refIRCId",
                                                                                 irc."refIRCName"
                                     FROM notes."refReportsHistory" rrh
                                              JOIN appointment."refAppointments" ra
                                                   ON ra."refAppointmentId" = rrh."refAppointmentId"
                                              JOIN impressionrecommendation."ImpressionRecommendationVal" irv
                                                   ON irv."refIRVCustId" = ra."refAppointmentRecommendationRight"
                                              JOIN impressionrecommendation."ImpressionRecommendationCategory" irc
                                                   ON irc."refIRCId" = irv."refIRCId"
                                     WHERE ra."refAppointmentStatus" = TRUE
                                       AND rrh."refRHHandleStatus" = 'Signed Off'
                                       AND irv."refIRVSystemType" = 'WR'
                                       AND ra."refSCId" = sc."refSCId"
                                     ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
                               WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                                         BETWEEN $1::date AND $2::date)

        SELECT COUNT(ls."refIRCId") AS total_count
        FROM impressionrecommendation."ImpressionRecommendationCategory" irc
                 LEFT JOIN latest_signed ls
                           ON ls."refIRCId" = irc."refIRCId"
        WHERE irc."refIRCId" = 7
        GROUP BY irc."refIRCName"
        ORDER BY irc."refIRCName")                                               AS rightRedo,
       (SELECT COUNT(*)
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refAppointmentDicomSide"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND ra."refSCId" = sc."refSCId"
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t
        WHERE t."refAppointmentDicomSide" = 'unilateralleft')                    as "unilateralleft",
       (SELECT COUNT(*)
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refAppointmentDicomSide"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND ra."refSCId" = sc."refSCId"
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t
        WHERE t."refAppointmentDicomSide" = 'unilateralright')                   as "unilateralright",
       (SELECT COUNT(*)
        FROM (SELECT *
              FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refAppointmentId",
                                                                rrh."refRHHandleEndTime",
                                                                rrh."refRHHandleStatus",
                                                                ra."refAppointmentDicomSide"
                    FROM notes."refReportsHistory" rrh
                             JOIN appointment."refAppointments" ra
                                  ON ra."refAppointmentId" = rrh."refAppointmentId"
                    WHERE ra."refAppointmentStatus" = TRUE
                      AND ra."refSCId" = sc."refSCId"
                      AND rrh."refRHHandleStatus" = 'Signed Off'
                    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
              WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
                        BETWEEN $1::date AND $2::date) AS t
        WHERE (t."refAppointmentDicomSide" = 'bilateral' OR t."refAppointmentDicomSide" IS NULL))                         as "bilateral"
--ENd
FROM public."ScanCenter" sc
WHERE (
          $3 = 0
              OR sc."refSCId" = $3
          )
`
