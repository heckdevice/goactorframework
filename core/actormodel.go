package core

type GenericDataPipe struct {
	dataChan            chan Message
	closeChan           chan bool
	isAcceptingMessages bool
}

// Actor - Actor model
type Actor struct {
	GenericDataPipe
	id                   string
	ActorType            string `json:"actor_type"`
	handlers             map[string]func(Message)
	internalMessageQueue messageStack
	owner                *actorSystem
}
