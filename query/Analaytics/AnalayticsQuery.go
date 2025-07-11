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

var AdminOverallScanIndicatesAnalayticsSQL = `
SELECT
  COUNT(*) AS total_appointments,
  COUNT(
    CASE
      WHEN "refCategoryId" = 1 THEN 1
    END
  ) AS "SForm",
  COUNT(
    CASE
      WHEN "refCategoryId" = 2 THEN 1
    END
  ) AS "DaForm",
  COUNT(
    CASE
      WHEN "refCategoryId" = 3 THEN 1
    END
  ) AS "DbForm",
  COUNT(
    CASE
      WHEN "refCategoryId" = 4 THEN 1
    END
  ) AS "DcForm"
FROM
  appointment."refAppointments"
WHERE
  TO_CHAR(
    TO_DATE("refAppointmentDate", 'YYYY-MM-DD'),
    'YYYY-MM'
  ) = ?
  AND (
    ? = 0
    OR "refSCId" = ?
  )
GROUP BY
  "refSCId"
ORDER BY
  "refSCId";
`

var GetAllScanCenter = `
SELECT
  *
FROM
  public."ScanCenter"
ORDER BY
  "refSCId" DESC
`
