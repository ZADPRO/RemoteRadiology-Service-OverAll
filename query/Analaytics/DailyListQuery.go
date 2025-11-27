package query

var GetDailyListSQL = `
SELECT *
FROM (SELECT DISTINCT ON (rrh."refAppointmentId") rrh."refRHHandleEndTime",
                                                  rrh."refAppointmentId",
                                                  rrh."refRHHandleEndTime" AS "AppointmentDate",
                                                  u."refUserCustId",
                                                  u."refUserFirstName",
                                                  CASE
                                                      WHEN ra."refCategoryId" = 1 THEN 'Sform'
                                                      WHEN ra."refCategoryId" = 2 THEN 'Daform'
                                                      WHEN ra."refCategoryId" = 3 THEN 'Dbform'
                                                      WHEN ra."refCategoryId" = 4 THEN 'Dcform'
                                                      ELSE '-'
                                                      END                        AS "refCategoryId",
                                                  CASE
                                                      WHEN ra."refAppointmentDicomSide" = 'bilateral' THEN 'Bilateral'
                                                      WHEN ra."refAppointmentDicomSide" = 'unilateralright'
                                                          THEN 'Unilateral Right'
                                                      WHEN ra."refAppointmentDicomSide" = 'unilateralleft'
                                                          THEN 'Unilateral Left'
                                                      ELSE 'Bilateral'
                                                      END                        AS "scanSide",
                                                  COALESCE(
                                                          (SELECT u."refUserCustId"
                                                           FROM notes."refReportsHistory" rrha
                                                                    JOIN public."Users" U on rrha."refRHHandledUserId" = U."refUserId"
                                                           WHERE rrha."refRHHandleStatus" = 'Draft'
                                                             AND rrha."refAppointmentId" = rrh."refAppointmentId"
                                                           ORDER BY rrh."refRHId" DESC
                                                           LIMIT 1), '-'
                                                  )                              AS "handlerName",
                                                  CASE
                                                      WHEN ra."refAppointmentImpression" = '' THEN '-'
                                                      ELSE ra."refAppointmentImpression"
                                                      END,
                                                  CASE
                                                      WHEN ra."refAppointmentRecommendation" = '' THEN '-'
                                                      ELSE ra."refAppointmentRecommendation"
                                                      END,
                                                  CASE
                                                      WHEN ra."refAppointmentImpressionRight" = '' THEN '-'
                                                      ELSE ra."refAppointmentImpressionRight"
                                                      END,
                                                  CASE
                                                      WHEN ra."refAppointmentRecommendationRight" = '' THEN '-'
                                                      ELSE ra."refAppointmentRecommendationRight"
                                                      END
      FROM notes."refReportsHistory" rrh
               JOIN appointment."refAppointments" ra
                    ON ra."refAppointmentId" = rrh."refAppointmentId"
               JOIN public."Users" U on ra."refUserId" = U."refUserId"
      WHERE ra."refAppointmentStatus" = TRUE
        AND rrh."refRHHandleStatus" = 'Signed Off'
      ORDER BY rrh."refAppointmentId", rrh."refRHId" DESC) s
WHERE TO_TIMESTAMP(s."refRHHandleEndTime", 'YYYY-MM-DD HH24:MI:SS')::date
          BETWEEN $1::date
          AND $2::date;
`
