package db

var statements = map[string]string{
  // inserts a new row into the articles table
  "create-article": `
  INSERT INTO articles (article_id, name)
    values(UUID_TO_BIN(?), ?)
  `,
  // soft deletes a article by id
  "delete-article": `
  UPDATE
    articles
  SET
    deleted_at = NOW()
  WHERE
    article_id = UUID_TO_BIN(?)
    AND deleted_at IS NULL
  `,
  // gets a single article row by id
  "get-article": `
  SELECT
    article_id, name
  FROM
    articles
  WHERE
    article_id = UUID_TO_BIN(?)
    AND deleted_at IS NULL
  `,
  // update a single article row by ID
  "update-article": `
  UPDATE
    articles
  SET
    name = ?
  WHERE
    article_id = UUID_TO_BIN(?)
    AND deleted_at IS NULL
  `,
}
