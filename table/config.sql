CREATE TABLE IF NOT EXISTS [config] (
	[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
	[key] NVARCHAR(32) NULL,
	[value] NVARCHAR(128) NULL,
	[time] NVARCHAR(64) Null,
	[des] NVARCHAR(128) Null
)