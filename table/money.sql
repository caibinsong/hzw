CREATE TABLE IF NOT EXISTS [money] (
	[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
	[user_id] NVARCHAR(32)  NULL,
	[gold] INTEGER  NULL,
	[diamond] INTEGER Null,
	[addrmbsum] INTEGER Null
)