package query

var GetCategoryDataSQL = `
SELECT
  *
FROM
  impressionrecommendation."ImpressionRecommendationCategory"
`

var CheckNewImpressionRecommendationSQL = `
SELECT
  COALESCE(MAX("refIRVOrderId"), 0) + 1 AS next_order_id,
  CASE
    WHEN EXISTS (
      SELECT
        1
      FROM
        impressionrecommendation."ImpressionRecommendationVal"
      WHERE
        "refIRVCustId" = $1
    ) THEN true
    ELSE false
  END AS "CheckStatus"
FROM
  impressionrecommendation."ImpressionRecommendationVal"
`

var CheckUpdateImpressionRecommendationSQL = `
SELECT
  COALESCE(MAX("refIRVOrderId"), 0) + 1 AS next_order_id,
  CASE
    WHEN EXISTS (
      SELECT
        1
      FROM
        impressionrecommendation."ImpressionRecommendationVal"
      WHERE
        "refIRVCustId" = $1
        AND "refIRVId" != $2
    ) THEN true
    ELSE false
  END AS "CheckStatus"
FROM
  impressionrecommendation."ImpressionRecommendationVal"
`

var InsertNewImpressionRecommendationSQL = `
INSERT INTO
  impressionrecommendation."ImpressionRecommendationVal" (
    "refIRCId",
    "refIRVOrderId",
    "refIRVSystemType",
    "refIRVCustId",
    "refIRVImpressionShortDesc",
    "refIRVImpressionLongDesc",
    "refIRVImpressionTextColor",
    "refIRVRecommendationShortDesc",
    "refIRVRecommendationLongDesc",
    "refIRVRecommendationTextColor"
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
`

var GetAllImpressionRecommendationSQL = `
SELECT
  *
FROM
  impressionrecommendation."ImpressionRecommendationVal" irv
  LEFT JOIN impressionrecommendation."ImpressionRecommendationCategory" irc ON irc."refIRCId" = irv."refIRCId"
ORDER BY
  irv."refIRVOrderId" ASC
`

var UpdateOrderImpressionRecommendationSQL = `
UPDATE
  impressionrecommendation."ImpressionRecommendationVal"
SET
  "refIRVOrderId" = $1
WHERE
  "refIRVId" = $2;
`

var UpdateImpressionRecommendationSQL = `
UPDATE
  impressionrecommendation."ImpressionRecommendationVal"
SET
  "refIRCId" = $1,
  "refIRVCustId" = $2,
  "refIRVImpressionShortDesc" = $3,
  "refIRVImpressionLongDesc" = $4,
  "refIRVImpressionTextColor" = $5,
  "refIRVRecommendationShortDesc" = $6,
  "refIRVRecommendationLongDesc" = $7,
  "refIRVRecommendationTextColor" = $8
WHERE
  "refIRVId" = $9;
`

var DeleteImpressionRecommendationSQL = `
DELETE FROM
  impressionrecommendation."ImpressionRecommendationVal"
WHERE
  "refIRVId" = $1;
`
