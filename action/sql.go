package action

var GetLatestMigrationSQL =
	`SELECT * FROM migrations
	ORDER BY id DESC
	LIMIT 1
`

var CreateMigrationsTable =
	`CREATE TABLE migrations (
	id INT(11) NOT NULL AUTO_INCREMENT,
	filename VARCHAR(255),
	created_at DATETIME,
	PRIMARY KEY(id)
)`

var CheckMigrationsTableExistsSQL =
	`SELECT TABLE_NAME AS tables
	FROM INFORMATION_SCHEMA.TABLES
	WHERE TABLE_SCHEMA = 'cw' AND TABLE_NAME = 'migrations'
`

var InsertMigrationRecordIntoTableSQL =
	`INSERT INTO migrations 
	(filename,created_at)
	VALUES
	(?,NOW())
`

