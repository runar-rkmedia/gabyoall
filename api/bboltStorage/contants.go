package bboltStorage

type PubType string
type PubVerb string

const (
	PubTypeEndpoint PubType = "endpoint"
	PubTypeRequest  PubType = "request"
	PubTypeSchedule PubType = "schedule"
	PubTypeStat     PubType = "stat"

	PubVerbCreate PubVerb = "create"
	PubVerbUpdate PubVerb = "update"
	// Marks the item as deleted in the database, but does not delete it
	PubVerbSoftDelete PubVerb = "soft-delete"
	// Removes all items permanently
	PubVerbClean PubVerb = "clean"
)
