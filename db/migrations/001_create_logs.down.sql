-- CASCADE drops the outbox FK constraint automatically if migrations are rolled back out of order
DROP TABLE IF EXISTS logs CASCADE;
