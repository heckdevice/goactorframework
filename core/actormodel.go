package core

type GenericDataPipe struct {
	dataChan            chan Message `json:"dataChan"`
	closeChan           chan bool    `json:"closeChan"`
	isacceptingmessages bool         `json:"accepting_messages"`
}

// Actor - Actor model
type Actor struct {
	GenericDataPipe
	id                   string                   `json:"id"`
	ActorType            string                   `json:"actor_type"`
	handlers             map[string]func(Message) `json:"handlers"`
	internalMessageQueue messageStack             `json:"message_queue"`
	owner                *actorSystem
}

func (actor *Actor) HasMessages() bool {
	return actor.internalMessageQueue.Len() != 0
}

func (actor *Actor) ScheduleActionableMessage(am *ActionableMessage) {
	actor.internalMessageQueue.Push(*am)
}
func (actor *Actor) StopAcceptingMessages() {
	actor.isacceptingmessages = false
}

func (actor *Actor) NoOfMessagesInQueue() int {
	return actor.internalMessageQueue.Len()
}
