CREATE TABLE IF NOT EXISTS [user] (
	[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
	[user] NVARCHAR(32)  NULL,
	[pwd] NVARCHAR(100)  NULL,
	[name] NVARCHAR(100) Null,
	[lv] NVARCHAR(100) null,
	[strength] INTEGER null,
	[maxstrength] INTEGER null,
	[maxstrengthtime] datetime null,
	[createtime] datetime null,
	[logintime] datetime null
)