package page

import (
	"github.com/eiannone/keyboard"
	"github.com/nsf/termbox-go"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
)

type PageManager struct {
	currentPage Page
	msg         chan pageMsg.PageMsg
}

func NewPageManager() *PageManager {
	return &PageManager{
		currentPage: nil,
		msg:         make(chan pageMsg.PageMsg),
	}
}

func (pm *PageManager) Init() {
	termbox.Init()
	//termbox.SetInputMode(termbox.InputEsc)
	go pm.listenForKeyboardInput()
}

func (pm *PageManager) listenForKeyboardInput() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
		termbox.Close()
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		pm.msg <- pageMsg.NewKeyMsg(char, key)
	}
}

func (pm *PageManager) Run() {
	for {
		pm.currentPage.Render()
		select {
		case msg := <-pm.msg:
			pageCmd := pm.currentPage.Update(msg)
			switch pageCmd {
			case QuitCommand:
				return
			default:
				continue
			}
		}
	}
}

func (pm *PageManager) ToMainMenu() {
	pm.currentPage = NewMainMenu(pm)
	pm.currentPage.Init()
}

func (pm *PageManager) ToGameRoom(roomId string) {
	pm.currentPage = NewGameRoom(pm, roomId)
	pm.currentPage.Init()
}

func (pm *PageManager) ToCreateRoomForm() {
	pm.currentPage = NewCreateRoomForm(pm)
	pm.currentPage.Init()
}

func (pm *PageManager) ToJoinRoomForm() {
	pm.currentPage = NewJoinRoomForm(pm)
	pm.currentPage.Init()
}
