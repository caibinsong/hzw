CREATE TABLE IF NOT EXISTS [adventurehero] (
	[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
	[adventureindex] INTEGER NULL,
	[name] NVARCHAR(100) NULL,
	[hp] REAL NULL,
	[atk] REAL NULL,
	[def] REAL NULL,
	[will] REAL NULL,
	[skill1_name] NVARCHAR(100) NULL,
	[skill1_num] REAL NULL,
	[skill1typename] NVARCHAR(100) NULL
)

