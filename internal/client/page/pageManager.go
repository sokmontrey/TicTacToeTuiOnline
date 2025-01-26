package page

import (
	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"
)

type PageManager struct {
	currentPage Page
	msg         chan PageMsg
}

func NewPageManager() *PageManager {
	return &PageManager{
		currentPage: nil,
		msg:         make(chan PageMsg),
	}
}

func (pm *PageManager) Init() {
	go pm.listenForKeyboardInput()
}

func (pm *PageManager) listenForKeyboardInput() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		pm.msg <- KeyMsg{char, key}
	}
}

func (pm *PageManager) Run() {
	for {
		tm.MoveCursor(1, 1)
		tm.Println(pm.currentPage.View())
		tm.Flush()
		select {
		case msg := <-pm.msg:
			pageCmd := pm.currentPage.Update(msg)
			switch pageCmd {
			case ProgramQuit:
				return
			}
		}
	}
}

func (pm *PageManager) ToMainMenu() {
	tm.Clear()
	pm.currentPage = NewMainMenu(pm)
	pm.currentPage.Init()
}

func (pm *PageManager) ToGameRoom(roomId string) {
	//pm.currentPage = NewGamePage(pm, roomId)
	//pm.currentPage.Init()
}

func (pm *PageManager) ToCreateRoomForm() {
	tm.Clear()
	pm.currentPage = NewCreateRoomForm(pm)
	pm.currentPage.Init()
}

func (pm *PageManager) ToJoinRoomForm() {
	tm.Clear()
	pm.currentPage = NewJoinRoomForm(pm)
	pm.currentPage.Init()
}
