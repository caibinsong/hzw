CREATE TABLE IF NOT EXISTS [skillinfo] (
	[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
	[name] NVARCHAR(64)  NULL,
	[quality] NVARCHAR(10)  NULL,
	[typename] NVARCHAR(10) Null,
	[initnum] REAL  null,
	[growthnum] REAL null,
	[des] NVARCHAR(512)  NULL
)