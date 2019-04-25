package core

// GenericDataPipe - Basic structure to facilitate a data and close channel
type GenericDataPipe struct {
	dataChan            chan Message
	closeChan           chan bool
	isAcceptingMessages bool
}

// Actor - Actor model with embedded data pipeline and message stack
type Actor struct {
	GenericDataPipe
	id                   string
	ActorType            string `json:"actor_type"`
	handlers             map[string]func(Message)
	internalMessageQueue messageStack
	owner                *actorSystem
}
