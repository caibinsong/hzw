CREATE TABLE IF NOT EXISTS [money_log] (
	[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
	[user_id] NVARCHAR(32)  NULL,
	[add_gold] INTEGER  NULL,
	[add_diamond] INTEGER Null,
	[add_rmb] INTEGER Null,
	[add_time] datetime null
)