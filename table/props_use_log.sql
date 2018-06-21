CREATE TABLE IF NOT EXISTS [props_use_log] (
	[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
	[userid] INTEGER  NULL,
	[name] NVARCHAR(64)  NULL,
	[quality] NVARCHAR(10)  NULL,
	[usenum] INTEGER  null,
	[checkmd5] NVARCHAR(64) null,
	[isused] INTEGER  NULL,
	[createtime]  NVARCHAR(64)  NULL,
	[usedtime]  NVARCHAR(64)  NULL,
	[act]  NVARCHAR(128)  NULL
)