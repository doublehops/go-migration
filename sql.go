package go_migration

const QuerySeparator = "------------------"

var GetLatestMigrationSQL = `SELECT * FROM migrations
	ORDER BY id DESC
	LIMIT 1
`

var CreateMigrationsTable = `CREATE TABLE migrations (
	id INT(11) NOT NULL AUTO_INCREMENT,
	filename VARCHAR(255),
	created_at DATETIME,
	PRIMARY KEY(id)
)`

var CheckMigrationsTableExistsSQL = `SHOW TABLES`

var InsertMigrationRecordIntoTableSQL = `INSERT INTO migrations 
	(filename,created_at)
	VALUES
	(?,NOW())
`

var RemoveMigrationRecordFromTableSQL = `DELETE FROM migrations 
	WHERE filename = ?
`
