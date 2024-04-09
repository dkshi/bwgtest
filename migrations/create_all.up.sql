CREATE TABLE IF NOT EXISTS quotations (
    update_id SERIAL PRIMARY KEY,
    code_from VARCHAR(3),
    code_to VARCHAR(3),
    rate REAL,
    update_time TIMESTAMP WITH TIME ZONE,
    success BOOLEAN
);

CREATE OR REPLACE FUNCTION notify_quotation_update() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM pg_notify('new_quotation_update', row_to_json(NEW)::text);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER quotation_trigger
AFTER INSERT ON quotations
FOR EACH ROW
EXECUTE FUNCTION notify_quotation_update();