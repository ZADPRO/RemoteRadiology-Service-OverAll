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
      WHEN ra."refCategoryId" = 1 THEN 1
    END
  ) AS "SForm",
  COUNT(
    CASE
      WHEN ra."refCategoryId" = 2 THEN 1
    END
  ) AS "DaForm",
  COUNT(
    CASE
      WHEN ra."refCategoryId" = 3 THEN 1
    END
  ) AS "DbForm",
  COUNT(
    CASE
      WHEN ra."refCategoryId" = 4 THEN 1
    END
  ) AS "DcForm"
FROM
  notes."refReportsHistory" rrh
  JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
WHERE
  rrh."refRHHandledUserId" = ?
  AND TO_CHAR(
    TO_DATE(
      rrh."refRHHandleStartTime",
      'YYYY-MM-DD HH24:MI:SS'
    ),
    'YYYY-MM'
  ) = ?;
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
      JOIN public."Users" rrhu 
        ON rrh."refRHHandledUserId" = rrhu."refUserId"
    WHERE
      ra."refAppointmentDate" >= ?
      AND ra."refAppointmentDate" <= ?
      AND (
        ? = 0 OR ra."refSCId" = ?
      )
      AND ( ? = FALSE OR rrhu."refRTId" IN (1, 6, 7, 10) )
    GROUP BY
      ra."refAppointmentImpression"
  ),
  expected_impressions AS (
    SELECT unnest(ARRAY[
      '1','1a',
      '2','2a',
      '3','3a','3b','3c','3d','3e','3f','3g',
      '4','4a','4b','4c','4d','4e','4f','4g','4h','4i','4j','4k','4l','4m',
      '5',
      '6','6a','6b','6c','6d','6e','6f','6g',
      '7a','7b','7c','7d','7e',
      '10','10a'
    ]) AS impression
  )
SELECT
  ei.impression,
  COALESCE(ac.count, 0) AS count
FROM
  expected_impressions ei
  LEFT JOIN actual_counts ac 
    ON ei.impression = ac.impression
ORDER BY
  ei.impression;
`
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
          '2a',
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
          '4h',
          '4i',
          '4j',
          '4k',
          '4l',
          '4m',
          '5',
          '6',
          '6a',
          '6b',
          '6c',
          '6d',
          '6e',
          '6f',
          '6g',
          '7a',
          '7b',
          '7c',
          '7d',
          '7e',
          '10',
          '10a'
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

var LeftRecommendationScancenterSQL = `
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
  groups AS (
    SELECT
      'Annual Screening' AS group_name
    UNION ALL
    SELECT
      'USG/ SFU'
    UNION ALL
    SELECT
      'Biopsy'
    UNION ALL
    SELECT
      'Breast radiologist'
    UNION ALL
    SELECT
      'Clinical Correlation'
    UNION ALL
    SELECT
      'Onco Consult'
    UNION ALL
    SELECT
      'Redo'
  ),
  counts AS (
    SELECT
      CASE
        WHEN ra."refAppointmentRecommendation" IN ('1', '1a', '7', '10') THEN 'Annual Screening'
        WHEN ra."refAppointmentRecommendation" IN (
          '2',
          '2a',
          '3',
          '3a',
          '3b',
          '3g',
          '4h',
          '4i1',
          '4i2',
          '4k'
        ) THEN 'USG/ SFU'
        WHEN ra."refAppointmentRecommendation" IN ('4g', '4n', '5', '5a', '6e') THEN 'Biopsy'
        WHEN ra."refAppointmentRecommendation" IN ('4a', '4c', '4d', '4e', '4j', '6d', '6g', '10a') THEN 'Breast radiologist'
        WHEN ra."refAppointmentRecommendation" IN (
          '3c',
          '3d',
          '3e',
          '3f',
          '4',
          '4b',
          '4f',
          '4l',
          '4m',
          '6',
          '6a',
          '6f',
          '6h',
          '7a',
          '7b',
          '7c',
          '7d',
          '7e',
          '8',
          '8a'
        ) THEN 'Clinical Correlation'
        WHEN ra."refAppointmentRecommendation" IN ('6b', '6c') THEN 'Onco Consult'
        WHEN ra."refAppointmentRecommendation" = '0' THEN 'Redo'
      END AS group_name,
      COUNT(*) AS total_count
    FROM
      latest_report_history rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
      JOIN public."Users" rrhu ON rrh."refRHHandledUserId" = rrhu."refUserId"
    WHERE
      rrh."refRHHandleStartTime" >= $1
      AND rrh."refRHHandleStartTime" <= $2
      AND (
        $3 = 0
        OR ra."refSCId" = $3
      ) -- scan center filter
      AND (
        $4 = FALSE
        OR rrhu."refRTId" IN (1, 6, 7, 10)
      ) -- user role filter
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleEndTime" IS NOT NULL
    GROUP BY
      group_name
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

var RightRecommendationScancenterSQL = `
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
  groups AS (
    SELECT
      'Annual Screening' AS group_name
    UNION ALL
    SELECT
      'USG/ SFU'
    UNION ALL
    SELECT
      'Biopsy'
    UNION ALL
    SELECT
      'Breast radiologist'
    UNION ALL
    SELECT
      'Clinical Correlation'
    UNION ALL
    SELECT
      'Onco Consult'
    UNION ALL
    SELECT
      'Redo'
  ),
  counts AS (
    SELECT
      CASE
        WHEN ra."refAppointmentRecommendationRight" IN ('1', '1a', '7', '10') THEN 'Annual Screening'
        WHEN ra."refAppointmentRecommendationRight" IN (
          '2',
          '2a',
          '3',
          '3a',
          '3b',
          '3g',
          '4h',
          '4i1',
          '4i2',
          '4k'
        ) THEN 'USG/ SFU'
        WHEN ra."refAppointmentRecommendationRight" IN ('4g', '4n', '5', '5a', '6e') THEN 'Biopsy'
        WHEN ra."refAppointmentRecommendationRight" IN ('4a', '4c', '4d', '4e', '4j', '6d', '6g', '10a') THEN 'Breast radiologist'
        WHEN ra."refAppointmentRecommendationRight" IN (
          '3c',
          '3d',
          '3e',
          '3f',
          '4',
          '4b',
          '4f',
          '4l',
          '4m',
          '6',
          '6a',
          '6f',
          '6h',
          '7a',
          '7b',
          '7c',
          '7d',
          '7e',
          '8',
          '8a'
        ) THEN 'Clinical Correlation'
        WHEN ra."refAppointmentRecommendationRight" IN ('6b', '6c') THEN 'Onco Consult'
        WHEN ra."refAppointmentRecommendationRight" = '0' THEN 'Redo'
      END AS group_name,
      COUNT(*) AS total_count
    FROM
      latest_report_history rrh
      JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
      JOIN public."Users" rrhu ON rrh."refRHHandledUserId" = rrhu."refUserId"
    WHERE
      rrh."refRHHandleStartTime" >= $1
      AND rrh."refRHHandleStartTime" <= $2
      AND (
        $3 = 0
        OR ra."refSCId" = $3
      ) -- scan center filter
      AND (
        $4 = FALSE
        OR rrhu."refRTId" IN (1, 6, 7, 10)
      ) -- user role filter
      AND rrh."refRHHandleStartTime" IS NOT NULL
      AND rrh."refRHHandleEndTime" IS NOT NULL
    GROUP BY
      group_name
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

var LeftRecommendationUserSQL = `
WITH
  groups AS (
    SELECT
      'Annual Screening' AS group_name
    UNION ALL
    SELECT
      'USG/ SFU'
    UNION ALL
    SELECT
      'Biopsy'
    UNION ALL
    SELECT
      'Breast radiologist'
    UNION ALL
    SELECT
      'Clinical Correlation'
    UNION ALL
    SELECT
      'Onco Consult'
    UNION ALL
    SELECT
      'Redo'
  ),
  counts AS (
    SELECT
      CASE
        WHEN t."refAppointmentRecommendation" IN ('1', '1a', '7', '10') THEN 'Annual Screening'
        WHEN t."refAppointmentRecommendation" IN (
          '2',
          '2a',
          '3',
          '3a',
          '3b',
          '3g',
          '4h',
          '4i1',
          '4i2',
          '4k'
        ) THEN 'USG/ SFU'
        WHEN t."refAppointmentRecommendation" IN ('4g', '4n', '5', '5a', '6e') THEN 'Biopsy'
        WHEN t."refAppointmentRecommendation" IN ('4a', '4c', '4d', '4e', '4j', '6d', '6g', '10a') THEN 'Breast radiologist'
        WHEN t."refAppointmentRecommendation" IN (
          '3c',
          '3d',
          '3e',
          '3f',
          '4',
          '4b',
          '4f',
          '4l',
          '4m',
          '6',
          '6a',
          '6f',
          '6h',
          '7a',
          '7b',
          '7c',
          '7d',
          '7e',
          '8',
          '8a'
        ) THEN 'Clinical Correlation'
        WHEN t."refAppointmentRecommendation" IN ('6b', '6c') THEN 'Onco Consult'
        WHEN t."refAppointmentRecommendation" = '0' THEN 'Redo'
        -- ELSE 'Other'
      END AS group_name,
      COUNT(*) AS total_count
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
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
    GROUP BY
      group_name
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
    SELECT
      'Annual Screening' AS group_name
    UNION ALL
    SELECT
      'USG/ SFU'
    UNION ALL
    SELECT
      'Biopsy'
    UNION ALL
    SELECT
      'Breast radiologist'
    UNION ALL
    SELECT
      'Clinical Correlation'
    UNION ALL
    SELECT
      'Onco Consult'
    UNION ALL
    SELECT
      'Redo'
  ),
  counts AS (
    SELECT
      CASE
        WHEN t."refAppointmentRecommendationRight" IN ('1', '1a', '7', '10') THEN 'Annual Screening'
        WHEN t."refAppointmentRecommendationRight" IN (
          '2',
          '2a',
          '3',
          '3a',
          '3b',
          '3g',
          '4h',
          '4i1',
          '4i2',
          '4k'
        ) THEN 'USG/ SFU'
        WHEN t."refAppointmentRecommendationRight" IN ('4g', '4n', '5', '5a', '6e') THEN 'Biopsy'
        WHEN t."refAppointmentRecommendationRight" IN ('4a', '4c', '4d', '4e', '4j', '6d', '6g', '10a') THEN 'Breast radiologist'
        WHEN t."refAppointmentRecommendationRight" IN (
          '3c',
          '3d',
          '3e',
          '3f',
          '4',
          '4b',
          '4f',
          '4l',
          '4m',
          '6',
          '6a',
          '6f',
          '6h',
          '7a',
          '7b',
          '7c',
          '7d',
          '7e',
          '8',
          '8a'
        ) THEN 'Clinical Correlation'
        WHEN t."refAppointmentRecommendationRight" IN ('6b', '6c') THEN 'Onco Consult'
        WHEN t."refAppointmentRecommendationRight" = '0' THEN 'Redo'
        -- ELSE 'Other'
      END AS group_name,
      COUNT(*) AS total_count
    FROM
      (
        SELECT DISTINCT
          ON (rrh."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          notes."refReportsHistory" rrh
          JOIN appointment."refAppointments" ra ON ra."refAppointmentId" = rrh."refAppointmentId"
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
    GROUP BY
      group_name
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

var TechArtificatsAll = `
SELECT
    COUNT(CASE WHEN ra."refAppointmentTechArtifactsLeft" = TRUE 
                AND ra."refAppointmentTechArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
    COUNT(CASE WHEN ra."refAppointmentTechArtifactsLeft" = FALSE 
                AND ra."refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS rightartifacts,
    COUNT(CASE WHEN ra."refAppointmentTechArtifactsLeft" = TRUE 
                AND ra."refAppointmentTechArtifactsRight" = TRUE THEN 1 END) AS bothartifacts
FROM
    appointment."refAppointments" ra
WHERE
    ($1 = 0 OR ra."refSCId" = $1) AND
    ra."refAppointmentDate" >= $2 AND
    ra."refAppointmentDate" <= $3;
`

var ReportArtificatsAll = `
SELECT
    COUNT(CASE WHEN ra."refAppointmentReportArtifactsLeft" = TRUE 
                AND ra."refAppointmentReportArtifactsRight" = FALSE THEN 1 END) AS leftartifacts,
    COUNT(CASE WHEN ra."refAppointmentReportArtifactsLeft" = FALSE 
                AND ra."refAppointmentReportArtifactsRight" = TRUE THEN 1 END) AS rightartifacts,
    COUNT(CASE WHEN ra."refAppointmentReportArtifactsLeft" = TRUE 
                AND ra."refAppointmentReportArtifactsRight" = TRUE THEN 1 END) AS bothartifacts
FROM
    appointment."refAppointments" ra
WHERE
    ($1 = 0 OR ra."refSCId" = $1) AND
    ra."refAppointmentDate" >= $2 AND
    ra."refAppointmentDate" <= $3;
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

var TotoalUserAnalayticsSQL = `
SELECT
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
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('1', '1a', '7', '10')
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
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        '2',
        '2a',
        '3',
        '3a',
        '3b',
        '3g',
        '4h',
        '4i1',
        '4i2',
        '4k'
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
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('4g', '4n', '5', '5a', '6e')
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
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('4a', '4c', '4d', '4e', '4j', '6d', '6g', '10a')
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
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        '3c',
        '3d',
        '3e',
        '3f',
        '4',
        '4b',
        '4f',
        '4l',
        '4m',
        '6',
        '6a',
        '6f',
        '6h',
        '7a',
        '7b',
        '7c',
        '7d',
        '7e',
        '8',
        '8a'
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
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('6b', '6c')
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
          rrh."refRHHandledUserId" = u."refUserId"
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('0')
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
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('1', '1a', '7', '10')
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
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        '2',
        '2a',
        '3',
        '3a',
        '3b',
        '3g',
        '4h',
        '4i1',
        '4i2',
        '4k'
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
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('4g', '4n', '5', '5a', '6e')
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
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('4a', '4c', '4d', '4e', '4j', '6d', '6g', '10a')
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
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        '3c',
        '3d',
        '3e',
        '3f',
        '4',
        '4b',
        '4f',
        '4l',
        '4m',
        '6',
        '6a',
        '6f',
        '6h',
        '7a',
        '7b',
        '7c',
        '7d',
        '7e',
        '8',
        '8a'
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
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('6b', '6c')
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
          AND rrh."refRHHandleStartTime" != ''
          AND rrh."refRHHandleStartTime" IS NOT NULL
          AND rrh."refRHHandleEndTime" IS NOT NULL
          AND rrh."refRHHandleStartTime"::timestamp >= $1
          AND rrh."refRHHandleStartTime"::timestamp <= $2
        ORDER BY
          rrh."refAppointmentId",
          rrh."refRHHandleStartTime" DESC
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('0')
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

var GetOverAllScanCenterList = `
 SELECT
  sc."refSCId",
  sc."refSCCustId",
  (
    SELECT
      COUNT(DISTINCT ra."refAppointmentId")
    FROM
      appointment."refAppointments" ra
    WHERE
      ra."refAppointmentDate"::timestamp >= $1
      AND ra."refAppointmentDate"::timestamp <= $2
  ) AS "totalcase",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refCategoryId" = 1 THEN ra."refAppointmentId"
        END
      ) AS category_count
    FROM
      appointment."refAppointments" ra
    WHERE
      ra."refAppointmentDate"::timestamp >= $1
      AND ra."refAppointmentDate"::timestamp <= $2
  ) AS "totalsform",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refCategoryId" = 2 THEN ra."refAppointmentId"
        END
      ) AS category_count
    FROM
      appointment."refAppointments" ra
    WHERE
      ra."refAppointmentDate"::timestamp >= $1
      AND ra."refAppointmentDate"::timestamp <= $2
  ) AS "totaldaform",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refCategoryId" = 3 THEN ra."refAppointmentId"
        END
      ) AS category_count
    FROM
      appointment."refAppointments" ra
    WHERE
      ra."refAppointmentDate"::timestamp >= $1
      AND ra."refAppointmentDate"::timestamp <= $2
  ) AS "totaldbform",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refCategoryId" = 4 THEN ra."refAppointmentId"
        END
      ) AS category_count
    FROM
      appointment."refAppointments" ra
    WHERE
      ra."refAppointmentDate"::timestamp >= $1
      AND ra."refAppointmentDate"::timestamp <= $2
  ) AS "totaldcform",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refAppointmentTechArtifactsLeft" = true THEN ra."refAppointmentId"
        END
      ) AS category_count
    FROM
      appointment."refAppointments" ra
    WHERE
      ra."refAppointmentDate"::timestamp >= $1
      AND ra."refAppointmentDate"::timestamp <= $2
  ) AS "techartificatsleft",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refAppointmentTechArtifactsRight" = true THEN ra."refAppointmentId"
        END
      ) AS category_count
    FROM
      appointment."refAppointments" ra
    WHERE
      ra."refAppointmentDate"::timestamp >= $1
      AND ra."refAppointmentDate"::timestamp <= $2
  ) AS "techartificatsright",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refAppointmentReportArtifactsLeft" = true THEN ra."refAppointmentId"
        END
      ) AS category_count
    FROM
      appointment."refAppointments" ra
    WHERE
      ra."refAppointmentDate"::timestamp >= $1
      AND ra."refAppointmentDate"::timestamp <= $2
  ) AS "reportartificatsleft",
  (
    SELECT
      COUNT(
        DISTINCT CASE
          WHEN ra."refAppointmentReportArtifactsRight" = true THEN ra."refAppointmentId"
        END
      ) AS category_count
    FROM
      appointment."refAppointments" ra
    WHERE
      ra."refAppointmentDate"::timestamp >= $1
      AND ra."refAppointmentDate"::timestamp <= $2
  ) AS "reportartificatsright",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('1', '1a', '7', '10')
  ) AS "leftannualscreening",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        '2',
        '2a',
        '3',
        '3a',
        '3b',
        '3g',
        '4h',
        '4i1',
        '4i2',
        '4k'
      )
  ) AS "leftusgsfu",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('4g', '4n', '5', '5a', '6e')
  ) AS "leftBiopsy",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('4a', '4c', '4d', '4e', '4j', '6d', '6g', '10a')
  ) AS "leftBreastradiologist",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendation" IN (
        '3c',
        '3d',
        '3e',
        '3f',
        '4',
        '4b',
        '4f',
        '4l',
        '4m',
        '6',
        '6a',
        '6f',
        '6h',
        '7a',
        '7b',
        '7c',
        '7d',
        '7e',
        '8',
        '8a'
      )
  ) AS "leftClinicalCorrelation",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('6b', '6c')
  ) AS "leftOncoConsult",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendation"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendation" IN ('0')
  ) AS "leftRedo",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('1', '1a', '7', '10')
  ) AS "rightannualscreening",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        '2',
        '2a',
        '3',
        '3a',
        '3b',
        '3g',
        '4h',
        '4i1',
        '4i2',
        '4k'
      )
  ) AS "rightusgsfu",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('4g', '4n', '5', '5a', '6e')
  ) AS "rightBiopsy",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('4a', '4c', '4d', '4e', '4j', '6d', '6g', '10a')
  ) AS "rightBreastradiologist",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN (
        '3c',
        '3d',
        '3e',
        '3f',
        '4',
        '4b',
        '4f',
        '4l',
        '4m',
        '6',
        '6a',
        '6f',
        '6h',
        '7a',
        '7b',
        '7c',
        '7d',
        '7e',
        '8',
        '8a'
      )
  ) AS "rightClinicalCorrelation",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('6b', '6c')
  ) AS "rightOncoConsult",
  (
    SELECT
      COUNT(*)
    FROM
      (
        SELECT DISTINCT
          ON (ra."refAppointmentId") ra."refAppointmentRecommendationRight"
        FROM
          appointment."refAppointments" ra
        WHERE
          ra."refSCId" = sc."refSCId"
          AND ra."refAppointmentDate"::timestamp >= $1
          AND ra."refAppointmentDate"::timestamp <= $2
        ORDER BY
          ra."refAppointmentId"
      ) t
    WHERE
      t."refAppointmentRecommendationRight" IN ('0')
  ) AS "rightOncoConsult"
FROM
  public."ScanCenter" sc
WHERE
  (
    $3 = 0
    OR sc."refSCId" = $3
  ) 
`