package query

var AdminOverallAnalayticsSQL = `
WITH months AS (
    SELECT 
        date_trunc('month', CURRENT_DATE) - interval '1 month' * n AS month_start
    FROM generate_series(0, 5) AS n
)
SELECT 
    TO_CHAR(month_start, 'YYYY-MM') AS month,
    TO_CHAR(month_start, 'Month') AS month_name,
    month_start::date AS starting_date,
    (month_start + INTERVAL '1 month - 1 day')::date AS ending_date,
  (SELECT COUNT(*) AS total_rows
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus"
    FROM
        notes."refReportsHistory" rrh
        JOIN appointment."refAppointments" ra 
            ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND (
            $1 = 0
            OR ra."refSCId" = $1
          )
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN month_start::date AND (month_start + INTERVAL '1 month - 1 day')::date
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t) AS total_appointments
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
SELECT 
  COUNT(*) AS total_appointments,
  COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END) AS "SForm",
  COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END) AS "DaForm",
  COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END) AS "DbForm",
  COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END) AS "DcForm"
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId"   -- ✅ REQUIRED for outer COUNT()
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date 
              AND $2::date
  AND (
            $3 = 0
            OR ra."refSCId" = $3
          )
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t;
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
      AND "refRHHandledUserId" = ?
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

var WellGreenUserIndicatesAnalayticsInvoiceSQL = `
SELECT
  COUNT(DISTINCT rrh."refAppointmentId") AS total_appointments,
  COUNT(
    CASE
      WHEN (
        ra."refCategoryId" = 1
        AND rrh."refRHHandleEdit" = 1
      ) THEN 1
    END
  ) AS "SFormEdit",
  COUNT(
    CASE
      WHEN (
        ra."refCategoryId" = 1
        AND rrh."refRHHandleCorrect" = 1
      ) THEN 1
    END
  ) AS "SFormCorrect",
  COUNT(
    CASE
      WHEN (
        ra."refCategoryId" = 2
        AND rrh."refRHHandleEdit" = 1
      ) THEN 1
    END
  ) AS "DaFormEdit",
  COUNT(
    CASE
      WHEN (
        ra."refCategoryId" = 2
        AND rrh."refRHHandleCorrect" = 1
      ) THEN 1
    END
  ) AS "DaFormCorrect",
  COUNT(
    CASE
      WHEN (
        ra."refCategoryId" = 3
        AND rrh."refRHHandleEdit" = 1
      ) THEN 1
    END
  ) AS "DbFormEdit",
  COUNT(
    CASE
      WHEN (
        ra."refCategoryId" = 4
        AND rrh."refRHHandleCorrect" = 1
      ) THEN 1
    END
  ) AS "DcFormCorrect",
  COUNT(
    CASE
      WHEN (
        ra."refCategoryId" = 4
        AND rrh."refRHHandleEdit" = 1
      ) THEN 1
    END
  ) AS "DcFormEdit",
  COUNT(
    CASE
      WHEN (
        ra."refCategoryId" = 3
        AND rrh."refRHHandleCorrect" = 1
      ) THEN 1
    END
  ) AS "DbFormCorrect"
FROM
  notes."refReportsHistory" rrh
  JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
WHERE
  rrh."refRHHandledUserId" = $1
  AND TO_CHAR(
    TO_DATE(
      rrh."refRHHandleStartTime",
      'YYYY-MM-DD HH24:MI:SS'
    ),
    'YYYY-MM'
  ) = $2;
`

var WellGreenUserIndicatesAnalayticsSQL = `
SELECT
  COUNT(DISTINCT (rrh."refAppointmentId", rrh."refRHHandledUserId")) AS total_appointments,
  COUNT(DISTINCT (rrh."refAppointmentId", rrh."refRHHandledUserId")) FILTER (
    WHERE ra."refCategoryId" = 1
  ) AS "SForm",
  COUNT(DISTINCT (rrh."refAppointmentId", rrh."refRHHandledUserId")) FILTER (
    WHERE ra."refCategoryId" = 2
  ) AS "DaForm",
  COUNT(DISTINCT (rrh."refAppointmentId", rrh."refRHHandledUserId")) FILTER (
    WHERE ra."refCategoryId" = 3
  ) AS "DbForm",
  COUNT(DISTINCT (rrh."refAppointmentId", rrh."refRHHandledUserId")) FILTER (
    WHERE ra."refCategoryId" = 4
  ) AS "DcForm"
FROM
  notes."refReportsHistory" rrh
  JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
WHERE
  rrh."refRHHandledUserId" = ?
  AND TO_DATE(rrh."refRHHandleStartTime", 'YYYY-MM-DD HH24:MI:SS') >= ?
  AND TO_DATE(rrh."refRHHandleStartTime", 'YYYY-MM-DD HH24:MI:SS') <= ?
;
`

var UserWorkedTimingSQL = `
SELECT
  rrh."refRHHandledUserId",
  SUM(
    EXTRACT(
      EPOCH
      FROM
        (
          TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS') - TO_TIMESTAMP(
            rrh."refRHHandleStartTime",
            'YYYY-MM-DD HH24:MI:SS'
          )
        )
    )
  ) / 60 AS total_minutes,
  ROUND(
    SUM(
      EXTRACT(
        EPOCH
        FROM
          (
            TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS') - TO_TIMESTAMP(
              rrh."refRHHandleStartTime",
              'YYYY-MM-DD HH24:MI:SS'
            )
          )
      )
    ) / 3600,
    2
  ) AS total_hours
FROM
  notes."refReportsHistory" rrh
WHERE
  rrh."refRHHandledUserId" = $1
  AND rrh."refRHHandleEndTime" IS NOT NULL
  AND rrh."refRHHandleStartTime" IS NOT NULL
  AND TO_TIMESTAMP(
    rrh."refRHHandleStartTime",
    'YYYY-MM-DD HH24:MI:SS'
  ) >= $2
  AND TO_TIMESTAMP(
    rrh."refRHHandleStartTime",
    'YYYY-MM-DD HH24:MI:SS'
  ) <= $3
GROUP BY
  rrh."refRHHandledUserId";
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
SELECT
  sc."refSCId",
  sc."refSCName",
  COUNT(DISTINCT rrh."refAppointmentId") AS total_appointments
FROM
  public."ScanCenter" sc
LEFT JOIN
  appointment."refAppointments" a ON sc."refSCId" = a."refSCId"
LEFT JOIN
  notes."refReportsHistory" rrh ON rrh."refAppointmentId" = a."refAppointmentId"
  AND rrh."refRHHandledUserId" = ?
  AND TO_TIMESTAMP(rrh."refRHHandleStartTime", 'YYYY-MM-DD HH24:MI:SS') >= ?  -- start_date param
  AND TO_TIMESTAMP(rrh."refRHHandleStartTime", 'YYYY-MM-DD HH24:MI:SS') <= ?  -- end_date param
GROUP BY
  sc."refSCId", sc."refSCName"
ORDER BY
  sc."refSCName";
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
SELECT
  SUM(COALESCE("refRHHandleCorrect", 0)) AS "totalCorrect",
  SUM(COALESCE("refRHHandleEdit", 0)) AS "totalEdit"
FROM
  notes."refReportsHistory"
WHERE
  "refRHHandledUserId" = ?
  AND "refRHHandleStartTime" != ''
  AND "refRHHandleStartTime" IS NOT NULL
  AND "refRHHandleStartTime"::timestamp >= ?  -- start_date parameter
  AND "refRHHandleStartTime"::timestamp <= ?  -- end_date parameter
;
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
WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND TO_TIMESTAMP(rrh."refRHHandleEndTime",'YYYY-MM-DD HH24:MI:SS')::date 
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
AND (
        $3 = 0
        OR ra."refSCId" = $3
      )
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    irc."refIRCName" AS group_name,
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName";
`

var RightRecommendationScancenterSQL = `
WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND TO_TIMESTAMP(rrh."refRHHandleEndTime",'YYYY-MM-DD HH24:MI:SS')::date 
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
AND (
        $3 = 0
        OR ra."refSCId" = $3
      )
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    irc."refIRCName" AS group_name,
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName";
`

var LeftRecommendationUserSQL = `
WITH
  groups AS (
    SELECT DISTINCT
      irc."refIRCName" AS group_name
    FROM
      impressionrecommendation."ImpressionRecommendationCategory" irc
      JOIN impressionrecommendation."ImpressionRecommendationVal" irv 
        ON irv."refIRCId" = irc."refIRCId"
    WHERE
      irv."refIRVSystemType" = 'WR'
  ),
  counts AS (
    SELECT
      irc."refIRCName" AS group_name,
      COUNT(*) AS total_count
    FROM (
      SELECT DISTINCT
        ON (rrh."refAppointmentId") ra."refAppointmentRecommendation"
      FROM
        notes."refReportsHistory" rrh
        JOIN appointment."refAppointments" ra 
          ON ra."refAppointmentId" = rrh."refAppointmentId"
      WHERE
        rrh."refRHHandledUserId" = $1
        AND rrh."refRHHandleStartTime" IS NOT NULL
        AND rrh."refRHHandleEndTime" IS NOT NULL
        AND rrh."refRHHandleStartTime" >= $2
        AND rrh."refRHHandleStartTime" <= $3
      ORDER BY
        rrh."refAppointmentId",
        rrh."refRHHandleStartTime" DESC
    ) t
    JOIN impressionrecommendation."ImpressionRecommendationVal" irv 
      ON irv."refIRVCustId" = t."refAppointmentRecommendation"
    JOIN impressionrecommendation."ImpressionRecommendationCategory" irc 
      ON irc."refIRCId" = irv."refIRCId"
    WHERE
      irv."refIRVSystemType" = 'WR'
    GROUP BY
      irc."refIRCName"
  )
SELECT
  g.group_name,
  COALESCE(c.total_count, 0) AS total_count
FROM
  groups g
  LEFT JOIN counts c ON g.group_name = c.group_name
ORDER BY
  g.group_name;
`

var RightRecommendationUserSQL = `
WITH
  groups AS (
    SELECT DISTINCT
      irc."refIRCName" AS group_name
    FROM
      impressionrecommendation."ImpressionRecommendationCategory" irc
      JOIN impressionrecommendation."ImpressionRecommendationVal" irv 
        ON irv."refIRCId" = irc."refIRCId"
    WHERE
      irv."refIRVSystemType" = 'WR'
  ),
  counts AS (
    SELECT
      irc."refIRCName" AS group_name,
      COUNT(*) AS total_count
    FROM (
      SELECT DISTINCT
        ON (rrh."refAppointmentId") ra."refAppointmentRecommendationRight"
      FROM
        notes."refReportsHistory" rrh
        JOIN appointment."refAppointments" ra 
          ON ra."refAppointmentId" = rrh."refAppointmentId"
      WHERE
        rrh."refRHHandledUserId" = $1
        AND rrh."refRHHandleStartTime" IS NOT NULL
        AND rrh."refRHHandleEndTime" IS NOT NULL
        AND rrh."refRHHandleStartTime" >= $2
        AND rrh."refRHHandleStartTime" <= $3
      ORDER BY
        rrh."refAppointmentId",
        rrh."refRHHandleStartTime" DESC
    ) t
    JOIN impressionrecommendation."ImpressionRecommendationVal" irv 
      ON irv."refIRVCustId" = t."refAppointmentRecommendationRight"
    JOIN impressionrecommendation."ImpressionRecommendationCategory" irc 
      ON irc."refIRCId" = irv."refIRCId"
    WHERE
      irv."refIRVSystemType" = 'WR'
    GROUP BY
      irc."refIRCName"
  )
SELECT
  g.group_name,
  COALESCE(c.total_count, 0) AS total_count
FROM
  groups g
  LEFT JOIN counts c ON g.group_name = c.group_name
ORDER BY
  g.group_name;
`

var TotalTATSQL = `
SELECT
  COUNT(*) FILTER (
    WHERE
      duration_days <= 1
  ) AS le_1_day,
  COUNT(*) FILTER (
    WHERE
      duration_days > 1
      AND duration_days <= 3
  ) AS le_3_days,
  COUNT(*) FILTER (
    WHERE
      duration_days > 3
      AND duration_days <= 7
  ) AS le_7_days,
  COUNT(*) FILTER (
    WHERE
      duration_days > 7
      AND duration_days <= 10
  ) AS le_10_days,
  COUNT(*) FILTER (
    WHERE
      duration_days > 10
  ) AS gt_10_days
FROM
  (
    SELECT
      rrh."refAppointmentId",
      (
        EXTRACT(
          EPOCH
          FROM
            (
              rrh."refRHHandleEndTime"::timestamp - rrh."refRHHandleStartTime"::timestamp
            )
        ) / 86400.0
      ) AS duration_days
    FROM
      notes."refReportsHistory" rrh
    WHERE
      rrh."refRHHandleStartTime" >= ?  -- start_date parameter
      AND rrh."refRHHandleStartTime" <= ?  -- end_date parameter
      AND rrh."refRHHandledUserId" = ?
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleEndTime" IS NOT NULL
  ) AS subquery;
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
SELECT 
 COUNT(CASE WHEN "refAppointmentTechArtifactsLeft" = TRUE 
                AND "refAppointmentTechArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
    COUNT(CASE WHEN "refAppointmentTechArtifactsLeft" = FALSE 
                AND "refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS rightartifacts,
    COUNT(CASE WHEN "refAppointmentTechArtifactsLeft" = TRUE 
                AND "refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS bothartifacts
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId",
  ra."refAppointmentTechArtifactsLeft",
  ra."refAppointmentTechArtifactsRight"
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date 
              AND $2::date
  AND (
            $3 = 0
            OR ra."refSCId" = $3
          )
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t;
`

var ReportArtificatsAll = `
SELECT 
 COUNT(CASE WHEN "refAppointmentReportArtifactsLeft" = TRUE 
                AND "refAppointmentReportArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
    COUNT(CASE WHEN "refAppointmentReportArtifactsLeft" = FALSE 
                AND "refAppointmentReportArtifactsRight" = TRUE THEN 1 END) AS rightartifacts,
    COUNT(CASE WHEN "refAppointmentReportArtifactsLeft" = TRUE 
                AND "refAppointmentReportArtifactsRight" = TRUE THEN 1 END) AS bothartifacts
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId",
        ra."refAppointmentReportArtifactsLeft",
  ra."refAppointmentReportArtifactsRight"
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date 
              AND $2::date
  AND (
            $3 = 0
            OR ra."refSCId" = $3
          )
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t;
`

var TechArtificats = `
SELECT
   COUNT(CASE WHEN ra."refAppointmentTechArtifactsLeft" = TRUE 
                AND ra."refAppointmentTechArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
    COUNT(CASE WHEN ra."refAppointmentTechArtifactsLeft" = FALSE 
                AND ra."refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS rightartifacts,
    COUNT(CASE WHEN ra."refAppointmentTechArtifactsLeft" = TRUE 
                AND ra."refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS bothartifacts
FROM (
    SELECT DISTINCT rrh."refAppointmentId"
    FROM notes."refReportsHistory" rrh
    WHERE rrh."refRHHandledUserId" = $1
	  AND rrh."refRHHandleStartTime" >= $2
    AND rrh."refRHHandleStartTime" <= $3
    AND rrh."refRHHandleStartTime" IS NOT NULL
    AND rrh."refRHHandleEndTime" IS NOT NULL
) uniq
JOIN appointment."refAppointments" ra 
  ON ra."refAppointmentId" = uniq."refAppointmentId";
`

var ReportArtificats = `
SELECT
   COUNT(CASE WHEN ra."refAppointmentReportArtifactsLeft" = TRUE 
                AND ra."refAppointmentReportArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
    COUNT(CASE WHEN ra."refAppointmentReportArtifactsLeft" = FALSE 
                AND ra."refAppointmentReportArtifactsRight" = TRUE THEN 1 END) AS rightartifacts,
    COUNT(CASE WHEN ra."refAppointmentReportArtifactsLeft" = TRUE 
                AND ra."refAppointmentReportArtifactsRight" = TRUE THEN 1 END) AS bothartifacts
FROM (
    SELECT DISTINCT rrh."refAppointmentId"
    FROM notes."refReportsHistory" rrh
    WHERE rrh."refRHHandledUserId" = $1
	  AND rrh."refRHHandleStartTime" >= $2
    AND rrh."refRHHandleStartTime" <= $3
    AND rrh."refRHHandleStartTime" IS NOT NULL
    AND rrh."refRHHandleEndTime" IS NOT NULL
) uniq
JOIN appointment."refAppointments" ra 
  ON ra."refAppointmentId" = uniq."refAppointmentId";
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
WITH
  months AS (
    SELECT
      TO_CHAR(
        date_trunc('month', CURRENT_DATE) - INTERVAL '1 month' * n,
        'YYYY-MM'
      ) AS month,
      TO_CHAR(
        date_trunc('month', CURRENT_DATE) - INTERVAL '1 month' * n,
        'Month'
      ) AS month_name,
      date_trunc('month', CURRENT_DATE) - INTERVAL '1 month' * n AS month_start,
      (
        date_trunc('month', CURRENT_DATE) - INTERVAL '1 month' * n + INTERVAL '1 month - 1 day'
      )::date AS month_end
    FROM
      generate_series(0, 5) AS n
  ),
  appointment_counts AS (
    SELECT
      TO_CHAR(
        TO_TIMESTAMP(a."refAppointmentDate", 'YYYY-MM-DD HH24:MI:SS'),
        'YYYY-MM'
      ) AS month,
      COUNT(DISTINCT a."refAppointmentId") AS total,
      COUNT(
        CASE
          WHEN a."refCategoryId" = 1 THEN 1
        END
      ) AS "SForm",
      COUNT(
        CASE
          WHEN a."refCategoryId" = 2 THEN 1
        END
      ) AS "DaForm",
      COUNT(
        CASE
          WHEN a."refCategoryId" = 3 THEN 1
        END
      ) AS "DbForm",
      COUNT(
        CASE
          WHEN a."refCategoryId" = 4 THEN 1
        END
      ) AS "DcForm"
    FROM
      appointment."refAppointments" a
    WHERE
      TO_TIMESTAMP(a."refAppointmentDate", 'YYYY-MM-DD HH24:MI:SS')
        >= date_trunc('month', CURRENT_DATE) - INTERVAL '6 months'
      AND (
        $1 = 0
        OR a."refSCId" = $1
      )
    GROUP BY
      month
  )
SELECT
  m.month,
  TRIM(m.month_name) AS month_name,
  COALESCE(a.total, 0) AS total_appointments,
  COALESCE(a."SForm", 0) AS "SForm",
  COALESCE(a."DaForm", 0) AS "DaForm",
  COALESCE(a."DbForm", 0) AS "DbForm",
  COALESCE(a."DcForm", 0) AS "DcForm"
FROM
  months m
  LEFT JOIN appointment_counts a ON m.month = a.month
ORDER BY
  m.month;
`

var TotoalUserAnalayticsSQL = `SELECT
  u."refUserId",
  u."refUserCustId",
  (
    SELECT
      COUNT(DISTINCT rrh."refAppointmentId")
    FROM
      notes."refReportsHistory" rrh
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "totalcase",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refCategoryId" = 1 THEN rrh."refAppointmentId"
        END
      ) AS category_count
    FROM
      notes."refReportsHistory" rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "totalsform",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refCategoryId" = 2 THEN rrh."refAppointmentId"
        END
      ) AS category_count
    FROM
      notes."refReportsHistory" rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "totaldaform",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refCategoryId" = 3 THEN rrh."refAppointmentId"
        END
      ) AS category_count
    FROM
      notes."refReportsHistory" rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "totaldaform",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refCategoryId" = 4 THEN rrh."refAppointmentId"
        END
      ) AS category_count
    FROM
      notes."refReportsHistory" rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "totaldcform",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refAppointmentTechArtifactsLeft" = true THEN rrh."refAppointmentId"
        END
      ) AS category_count
    FROM
      notes."refReportsHistory" rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "techartificatsleft",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refAppointmentTechArtifactsRight" = true THEN rrh."refAppointmentId"
        END
      ) AS category_count
    FROM
      notes."refReportsHistory" rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "techartificatsright",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refAppointmentReportArtifactsLeft" = true THEN rrh."refAppointmentId"
        END
      ) AS category_count
    FROM
      notes."refReportsHistory" rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "reportartificatsleft",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refAppointmentReportArtifactsRight" = true THEN rrh."refAppointmentId"
        END
      ) AS category_count
    FROM
      notes."refReportsHistory" rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "reportartificatsright",
  (
    SELECT
      ROUND(
        SUM(
          EXTRACT(
            EPOCH
            FROM
              (
                TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS') - TO_TIMESTAMP(
                  rrh."refRHHandleStartTime",
                  'YYYY-MM-DD HH24:MI:SS'
                )
              )
          )
        ) / 3600,
        2
      ) AS total_hours
    FROM
      notes."refReportsHistory" rrh
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleEndTime" IS NOT NULL
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
    GROUP BY
      rrh."refRHHandledUserId"
  ) AS "totaltiming",
  (
    SELECT
      SUM(COALESCE("refRHHandleCorrect", 0)) AS "totalCorrect"
    FROM
      notes."refReportsHistory" rrh
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "totalreportcorrect",
  (
    SELECT
      SUM(COALESCE("refRHHandleEdit", 0)) AS "totalCorrect"
    FROM
      notes."refReportsHistory" rrh
    WHERE
      rrh."refRHHandledUserId" = u."refUserId"
      AND rrh."refRHHandleStartTime" != ''
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleStartTime"::timestamp >= $1
      AND rrh."refRHHandleStartTime"::timestamp <= $2
  ) AS "totalreportedit",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = 1
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        SELECT
          "refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Annual Screening' -- 👈 category filter
      )
  ) AS "leftannualscreening",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = 1
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        SELECT
          "refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'USG/SFU' -- 👈 category filter
      )
  ) AS "leftusgsfu",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = 1
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        SELECT
          "refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Biopsy' -- 👈 category filter
      )
  ) AS "leftBiopsy",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = 1
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        SELECT
          "refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Breast Radiologist' -- 👈 category filter
      )
  ) AS "leftBreastradiologist",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = 1
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        SELECT
          "refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Clinical Correlation' -- 👈 category filter
      )
  ) AS "leftClinicalCorrelation",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = 1
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        SELECT
          "refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Onco Consult' -- 👈 category filter
      )
  ) AS "leftOncoConsult",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = 1
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        SELECT
          "refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Redo' -- 👈 category filter
      )
  ) AS "leftRedo",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        SELECT
          irv."refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Annual Screening' -- 👈 Dynamic category name
      )
  ) AS "rightannualscreening",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        SELECT
          irv."refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'USG/SFU' -- 👈 Dynamic category name
      )
  ) AS "rightusgsfu",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        SELECT
          irv."refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Biopsy' -- 👈 Dynamic category name
      )
  ) AS "rightBiopsy",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        SELECT
          irv."refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Breast Radiologist' -- 👈 Dynamic category name
      )
  ) AS "rightBreastradiologist",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        SELECT
          irv."refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Clinical Correlation' -- 👈 Dynamic category name
      )
  ) AS "rightClinicalCorrelation",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        SELECT
          irv."refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Onco Consult' -- 👈 Dynamic category name
      )
  ) AS "rightOncoConsult",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
        WHERE
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleStartTime" <> ''
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        SELECT
          irv."refIRVCustId"
        FROM
          impressionrecommendation."ImpressionRecommendationVal" irv
          JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
        WHERE
          irv."refIRVSystemType" = 'WR'
          AND irc."refIRCName" = 'Redo' -- 👈 Dynamic category name
      )
  ) AS "rightRedo"
FROM
  public."Users" u
  FULL JOIN map."refScanCenterMap" rscm ON rscm."refUserId" = u."refUserId"
WHERE
  u."refRTId" NOT IN (3, 4, 9)
  AND (
    $3 = 0
    OR rscm."refSCId" = $3
  )
ORDER BY
  u."refUserId";
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
SELECT
  sc."refSCId",
  sc."refSCCustId",
  (SELECT COUNT(*) AS total_rows
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus"
    FROM
        notes."refReportsHistory" rrh
        JOIN appointment."refAppointments" ra 
            ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND ra."refSCId" = sc."refSCId"
        AND rrh."refRHHandleStatus" = 'Signed Off'
        AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t) AS totalcase,
  (SELECT 
  COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END) AS "SForm"
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId"   -- ✅ REQUIRED for outer COUNT()
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t) AS totalsform,
  (SELECT 
  COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END) AS "SForm"
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId"   -- ✅ REQUIRED for outer COUNT()
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t) AS totaldaform,
  (SELECT 
  COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END) AS "SForm"
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId"   -- ✅ REQUIRED for outer COUNT()
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t) AS totaldbform,
  (SELECT 
  COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END) AS "SForm"
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId"   -- ✅ REQUIRED for outer COUNT()
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t) AS totaldcform,
  (
  SELECT 
 COUNT(CASE WHEN "refAppointmentTechArtifactsLeft" = TRUE THEN 1 END) AS leftartifacts
    -- COUNT(CASE WHEN "refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS rightartifacts
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId",
  ra."refAppointmentTechArtifactsLeft",
  ra."refAppointmentTechArtifactsRight"
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t
  ) AS techartificatsleft,
  (
  SELECT 
    COUNT(CASE WHEN "refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS rightartifacts
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId",
  ra."refAppointmentTechArtifactsLeft",
  ra."refAppointmentTechArtifactsRight"
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t
  ) AS techartificatsright,
  (
  SELECT 
 COUNT(CASE WHEN "refAppointmentReportArtifactsLeft" = TRUE THEN 1 END) AS leftartifacts
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId",
        ra."refAppointmentReportArtifactsLeft",
  ra."refAppointmentReportArtifactsRight"
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t
  ) AS reportartificatsleft,
  (
  SELECT 
 COUNT(CASE WHEN "refAppointmentReportArtifactsRight" = TRUE THEN 1 END) AS rightartifacts
FROM (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
        rrh."refRHHandleEndTime",
        rrh."refRHHandleStatus",
        ra."refCategoryId",
        ra."refAppointmentReportArtifactsLeft",
  ra."refAppointmentReportArtifactsRight"
    FROM notes."refReportsHistory" rrh
    JOIN appointment."refAppointments" ra 
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
) AS t
  ) AS reportartificatsright,
  --Left Annual Screening
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 1
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS leftannualscreening,
  -- Left USG/SFU
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 2
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS leftusgsfu,
-- Left Biopsy
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 3
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS leftBiopsy,
  --Left Breast Radiologist
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 4
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS leftBreastradiologist,
  --Left Clinical Correlation
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 5
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS leftClinicalCorrelation,
  --Left Onco Consult
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 6
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS leftOncoConsult,
  --Left Redo
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 7
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS leftRedo,
  --Right Annual Screening
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 1
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS rightannualscreening,
  --Right USG/SFU
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 2
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS rightusgsfu,
  --Right Biopsy
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 3
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS rightBiopsy,
  --Right Breast Radiologist
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 4
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS rightBreastradiologist,
  --Right Clinical Correlation
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 5
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS rightClinicalCorrelation,
  --Right Onco Consult
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 6
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS rightOncoConsult,
  --Right Redo
  (
  WITH latest_signed AS (
    SELECT DISTINCT ON (rrh."refAppointmentId")
        rrh."refAppointmentId",
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
    WHERE
        ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
       AND TO_TIMESTAMP(rrh."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
            BETWEEN $1::date AND $2::date
  AND irv."refIRVSystemType" = 'WR'
  AND ra."refSCId" = sc."refSCId"
    ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC
)

SELECT 
    COUNT(ls."refIRCId") AS total_count
FROM impressionrecommendation."ImpressionRecommendationCategory" irc
LEFT JOIN latest_signed ls 
    ON ls."refIRCId" = irc."refIRCId"
  WHERE irc."refIRCId" = 7
GROUP BY irc."refIRCName"
ORDER BY irc."refIRCName"
  ) AS rightRedo
  --ENd
FROM
  public."ScanCenter" sc
WHERE
  (
    $3 = 0
    OR sc."refSCId" = $3
  )
`
