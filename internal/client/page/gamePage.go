package page

type GamePage struct {
	roomId string
}

func NewGamePage(roomId string) GamePage {
	return GamePage{
		roomId: roomId,
	}
}

func (m GamePage) Run() Page {
	return m
}
