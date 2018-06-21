CREATE TABLE IF NOT EXISTS [usersign] (
	[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
	[user_id] NVARCHAR(32)  NULL,
	[sign_type] INTEGER  NULL,
	[sign_time] datetime Null
)