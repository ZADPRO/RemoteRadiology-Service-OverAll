package query

var AdminOverallAnalayticsSQL = `
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
      ) AS month_name
    FROM
      generate_series(0, 5) AS n
  ),
  appointment_counts AS (
    SELECT
      TO_CHAR(TO_DATE("refAppointmentDate", 'YYYY-MM-DD'), 'YYYY-MM') AS month,
      COUNT(*) AS total
    FROM
      appointment."refAppointments"
    WHERE
      TO_DATE("refAppointmentDate", 'YYYY-MM-DD') >= date_trunc('month', CURRENT_DATE) - INTERVAL '6 months'
      AND (? = 0 OR "refSCId" = ?)
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
  COALESCE(SUM(total_appointments), 0) AS total_appointments,
  COALESCE(SUM("SForm"), 0) AS "SForm",
  COALESCE(SUM("DaForm"), 0) AS "DaForm",
  COALESCE(SUM("DbForm"), 0) AS "DbForm",
  COALESCE(SUM("DcForm"), 0) AS "DcForm"
FROM (
  SELECT
    COUNT(*) AS total_appointments,
    COUNT(CASE WHEN "refCategoryId" = 1 THEN 1 END) AS "SForm",
    COUNT(CASE WHEN "refCategoryId" = 2 THEN 1 END) AS "DaForm",
    COUNT(CASE WHEN "refCategoryId" = 3 THEN 1 END) AS "DbForm",
    COUNT(CASE WHEN "refCategoryId" = 4 THEN 1 END) AS "DcForm"
  FROM
    appointment."refAppointments"
  WHERE
    TO_DATE("refAppointmentDate", 'YYYY-MM-DD') BETWEEN TO_DATE(?, 'YYYY-MM-DD') AND TO_DATE(?, 'YYYY-MM-DD')
    AND (
      ? = 0
      OR "refSCId" = ?
    )
) AS stats;
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

// var WellGreenUserIndicatesAnalayticsSQL = `
// SELECT
//   COUNT(DISTINCT rrh."refAppointmentId") AS total_appointments,
//   COUNT(
//     CASE
//       WHEN ra."refCategoryId" = 1 THEN 1
//     END
//   ) AS "SForm",
//   COUNT(
//     CASE
//       WHEN ra."refCategoryId" = 2 THEN 1
//     END
//   ) AS "DaForm",
//   COUNT(
//     CASE
//       WHEN ra."refCategoryId" = 3 THEN 1
//     END
//   ) AS "DbForm",
//   COUNT(
//     CASE
//       WHEN ra."refCategoryId" = 4 THEN 1
//     END
//   ) AS "DcForm"
// FROM
//   notes."refReportsHistory" rrh
//   JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
// WHERE
//   rrh."refRHHandledUserId" = ?
//   AND TO_CHAR(
//     TO_DATE(
//       rrh."refRHHandleStartTime",
//       'YYYY-MM-DD HH24:MI:SS'
//     ),
//     'YYYY-MM'
//   ) = ?;
// `

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
  rrh."refRHHandledUserId" = ?
  AND rrh."refRHHandleEndTime" IS NOT NULL
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
  AND "refRHHandleStartTime"::timestamp >= ?  -- start_date parameter
  AND "refRHHandleStartTime"::timestamp <= ?  -- end_date parameter
;
`

// WITH impressions AS (
//   SELECT '1' AS impression UNION ALL
//   SELECT '1a' UNION ALL
//   SELECT '2' UNION ALL
//   SELECT '3' UNION ALL
//   SELECT '3a' UNION ALL
//   SELECT '3b' UNION ALL
//   SELECT '3c' UNION ALL
//   SELECT '3d' UNION ALL
//   SELECT '3e' UNION ALL
//   SELECT '3f' UNION ALL
//   SELECT '3g' UNION ALL
//   SELECT '4' UNION ALL
//   SELECT '4a' UNION ALL
//   SELECT '4b' UNION ALL
//   SELECT '4c' UNION ALL
//   SELECT '4d' UNION ALL
//   SELECT '4e' UNION ALL
//   SELECT '4f' UNION ALL
//   SELECT '4g' UNION ALL
//   SELECT '5' UNION ALL
//   SELECT '6' UNION ALL
//   SELECT '6a' UNION ALL
//   SELECT '6b' UNION ALL
//   SELECT '6c' UNION ALL
//   SELECT '6d' UNION ALL
//   SELECT '6e' UNION ALL
//   SELECT '6f' UNION ALL
//   SELECT '7a' UNION ALL
//   SELECT '7b' UNION ALL
//   SELECT '7c' UNION ALL
//   SELECT '7d' UNION ALL
//   SELECT '7e'
// ),
// actual_counts AS (
//   SELECT
//     ra."refAppointmentImpression" AS impression,
//     COUNT(*) AS count
//   FROM
//     notes."refReportsHistory" rrh
//     JOIN appointment."refAppointments" ra
//       ON ra."refAppointmentId" = rrh."refAppointmentId"
//   WHERE
//     rrh."refRHHandledUserId" = ?
//     AND TO_CHAR(ra."refAppointmentDate"::timestamp, 'YYYY-MM') = ?
//   GROUP BY
//     ra."refAppointmentImpression"
// )
// SELECT
//   i.impression,
//   COALESCE(a.count, 0) AS count
// FROM
//   impressions i
//   LEFT JOIN actual_counts a ON i.impression = a.impression
// ORDER BY
//   i.impression;

// var ImpressionNRecommentationScanCenterSQL = `
// WITH
//   latest_report_history AS (
//     SELECT *
//     FROM (
//       SELECT *,
//         ROW_NUMBER() OVER (
//           PARTITION BY "refAppointmentId"
//           ORDER BY "refRHId" DESC
//         ) AS rn
//       FROM notes."refReportsHistory"
//     ) sub
//     WHERE rn = 1
//   ),
//   actual_counts AS (
//     SELECT
//       ra."refAppointmentImpression" AS impression,
//       COUNT(*) AS count
//     FROM
//       latest_report_history rrh
//       JOIN appointment."refAppointments" ra
//         ON ra."refAppointmentId" = rrh."refAppointmentId"
//     WHERE
//       TO_CHAR(ra."refAppointmentDate"::timestamp, 'YYYY-MM') = ?
//       AND (
//         ? = 0 OR ra."refSCId" = ?
//       )
//     GROUP BY
//       ra."refAppointmentImpression"
//   ),
//   expected_impressions AS (
//     SELECT unnest(ARRAY[
//       '1','1a','2','3','3a','3b','3c','3d','3e','3f','3g',
//       '4','4a','4b','4c','4d','4e','4f','4g',
//       '5','6','6a','6b','6c','6d','6e','6f',
//       '7a','7b','7c','7d','7e'
//     ]) AS impression
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

var ImpressionNRecommentationScanCenterSQL = `
WITH
  latest_report_history AS (
    SELECT *
    FROM (
      SELECT *,
        ROW_NUMBER() OVER (
          PARTITION BY "refAppointmentId"
          ORDER BY "refRHId" DESC
        ) AS rn
      FROM notes."refReportsHistory"
    ) sub
    WHERE rn = 1
  ),
  actual_counts AS (
    SELECT
      ra."refAppointmentImpression" AS impression,
      COUNT(*) AS count
    FROM
      latest_report_history rrh
      JOIN appointment."refAppointments" ra
        ON ra."refAppointmentId" = rrh."refAppointmentId"
    WHERE
      ra."refAppointmentDate" >= ?
      AND ra."refAppointmentDate" <= ? 
      AND (
        ? = 0 OR ra."refSCId" = ?
      )
    GROUP BY
      ra."refAppointmentImpression"
  ),
  expected_impressions AS (
    SELECT unnest(ARRAY[
      '1','1a','2','3','3a','3b','3c','3d','3e','3f','3g',
      '4','4a','4b','4c','4d','4e','4f','4g',
      '5','6','6a','6b','6c','6d','6e','6f',
      '7a','7b','7c','7d','7e','10'
    ]) AS impression
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
//       TO_CHAR(ra."refAppointmentDate"::timestamp, 'YYYY-MM') = ?
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
//           '5',
//           '6',
//           '6a',
//           '6b',
//           '6c',
//           '6d',
//           '6e',
//           '6f',
//           '7a',
//           '7b',
//           '7c',
//           '7d',
//           '7e'
//         ]
//       ) AS impression
//   )
// SELECT
//   ei.impression,
//   COALESCE(ac.count, 0) AS count
// FROM
//   expected_impressions ei
//   LEFT JOIN actual_counts ac ON ei.impression = ac.impression;
// `

var ImpressionNRecommentationSQL = `
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
        WHERE
          "refRHHandledUserId" = ?
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
    WHERE
      ra."refAppointmentDate" >= ?  -- start_date parameter
      AND ra."refAppointmentDate" <= ?  -- end_date parameter
    GROUP BY
      ra."refAppointmentImpression"
  ),
  expected_impressions AS (
    SELECT
      unnest(
        array[
          '1',
          '1a',
          '2',
          '3',
          '3a',
          '3b',
          '3c',
          '3d',
          '3e',
          '3f',
          '3g',
          '4',
          '4a',
          '4b',
          '4c',
          '4d',
          '4e',
          '4f',
          '4g',
          '5',
          '6',
          '6a',
          '6b',
          '6c',
          '6d',
          '6e',
          '6f',
          '7a',
          '7b',
          '7c',
          '7d',
          '7e',
          '10'
        ]
      ) AS impression
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
