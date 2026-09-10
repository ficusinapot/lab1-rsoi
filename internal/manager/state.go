package manager

type State string

const (
	Created              State = "created"
	Prepared             State = "prepared"
	Running              State = "running"
	CommunicationStopped State = "communication_stopped"
	Shutdown             State = "shutdown"
	Disposed             State = "disposed"
)

func (s State) ValidateTransition(next State) error {
	if s == next {
		return nil
	}

	valid := map[State][]State{
		Created:              {Prepared, Disposed},
		Prepared:             {Running, Shutdown, Disposed},
		Running:              {CommunicationStopped, Shutdown, Disposed},
		CommunicationStopped: {Shutdown, Disposed},
		Shutdown:             {Disposed},
		Disposed:             {},
	}

	for _, candidate := range valid[s] {
		if candidate == next {
			return nil
		}
	}

	return InvalidStateError.New("invalid manager state transition: %s -> %s", s, next)
}
