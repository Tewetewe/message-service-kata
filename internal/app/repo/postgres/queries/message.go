package queries

const (
	// QueryCreateMessage query to create message
	QueryCreateMessage = `
	INSERT INTO consumed_messages (message, trigger_by)
	VALUES ($1, $2)
	RETURNING id;`
)
