package scheduler

import "github.com/gofrs/uuid"

// sql to drop an already existing trigger
var postgresDropExistingTriggerSQL = `
DROP TRIGGER IF EXISTS feeds_notify_event ON feeds;
`

// sql to call notify_event whenever a new is inserted into `feeds`
// based on: https://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/
var postgresCreateTriggerSQL = `
CREATE TRIGGER feeds_notify_event
AFTER INSERT OR DELETE OR UPDATE OF feed_url ON feeds
    FOR EACH ROW EXECUTE PROCEDURE notify_event();
`

// sql to define a function that sends a notification on the feed_change channel
// based on: https://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/
var postgresCreateNotifyFunctionSQL = `
CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE
        feed_id uuid;
        notification json;
    BEGIN

        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            feed_id = OLD.id;
        ELSE
            feed_id = NEW.id;
        END IF;

        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table', TG_TABLE_NAME,
                          'action', TG_OP,
                          'feed_id', feed_id);

        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('feed_change', notification::text);

        -- Result is ignored since this is an AFTER trigger
        RETURN NULL;
    END;

$$ LANGUAGE plpgsql;
`

// this is for deserializing the notification json above
type postgresNotificationPayload struct {
	Table  string    `json:"table"`
	Action string    `json:"action"`
	FeedID uuid.UUID `json:"feed_id"`
}
