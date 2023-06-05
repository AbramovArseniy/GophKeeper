CREATE TABLE IF NOT EXISTS keeper(
		login INT NOT NULL,
		data BYTEA NOT NULL,
		type SMALLINT NOT NULL,
		name VARCHAR NOT NULL,
		UNIQUE(login, type, name)
)