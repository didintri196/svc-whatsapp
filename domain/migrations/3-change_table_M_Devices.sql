-- +migrate Up
ALTER TABLE public."M_Devices" RENAME TO devices;
ALTER TABLE public.devices RENAME COLUMN "M_User_id" TO m_user_id;

-- +migrate Down
DROP TABLE IF EXISTS "M_Devices";