UPDATE todos
SET priority = 'med'
WHERE priority IS NULL;

ALTER TABLE todos
ALTER COLUMN priority SET NOT NULL,
ALTER COLUMN priority SET DEFAULT 'med';
