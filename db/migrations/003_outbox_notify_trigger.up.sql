CREATE OR REPLACE FUNCTION notify_outbox_insert()
RETURNS trigger AS $$
BEGIN
  PERFORM pg_notify('outbox_ready', NEW.id::text);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER outbox_insert_notify
  AFTER INSERT ON outbox
  FOR EACH ROW EXECUTE FUNCTION notify_outbox_insert();
