package query

var ListTrainingFilesSQL = `
SELECT
  *
FROM
  public."TrainingMaterial"
WHERE
  "refTMStatus" = true
`

var DeleteTrainingSQL = `
UPDATE
  public."TrainingMaterial"
SET
  "refTMStatus" = false
WHERE
  "refTMId" = ?
`

var OneListTrainingFilesSQL = `
SELECT
  *
FROM
  public."TrainingMaterial"
WHERE
  "refTMStatus" = true
  AND "refTMId" = ?
`
