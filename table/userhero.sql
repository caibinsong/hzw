CREATE TABLE IF NOT EXISTS [userhero] (
	[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
	[userid] INTEGER  NULL,
	[heroid] INTEGER  NULL,
	[name] NVARCHAR(64)  NULL,
	[lv] INTEGER Null,
	[hp] REAL Null,
	[atk] REAL  null,
	[def] REAL null,
	[will] REAL null,
	[pot] REAL null,
	[nowexp] INTEGER null,
	[newselfskill] NVARCHAR(512)  NULL,
	[selfskilllv] INTEGER null
)