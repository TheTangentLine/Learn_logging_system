DROP TRIGGER IF EXISTS outbox_insert_notify ON outbox;
DROP FUNCTION IF EXISTS notify_outbox_insert();
