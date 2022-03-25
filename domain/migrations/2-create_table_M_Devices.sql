-- +migrate Up
CREATE TABLE IF NOT EXISTS "M_Devices"
(
    "id"          char(36) PRIMARY KEY DEFAULT (uuid_generate_v4()),
    "M_User_id"   char(36) NOT NULL,
    "jid"         varchar(255) NOT NULL,
    "server"      varchar(255) NOT NULL,
    "phone"       varchar(255) NOT NULL,
    "worker_id"   varchar(255) NOT NULL,
    "api_key"     varchar(255) NOT NULL,
    "created_at"  timestamp    NOT NULL,
    "updated_at"  timestamp    NOT NULL,
    "deleted_at"  timestamp
    );

-- +migrate Down
DROP TABLE IF EXISTS "M_Devices";