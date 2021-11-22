-- ALTER TABLE  dbo.t_users ALTER COLUMN user_country_id DROP DEFAULT;
ALTER TABLE dbo.t_users ALTER COLUMN user_country_id DROP NOT NULL;
ALTER TABLE dbo.t_users ALTER COLUMN user_middle_name DROP NOT NULL;