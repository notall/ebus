package ebus

import "fmt"

var (
	ErrorDuplicatedHandler = fmt.Errorf("duplicated handler subscribed to same topic")
	ErrorTopicNotExist     = fmt.Errorf("topic doesn't exist")
)
