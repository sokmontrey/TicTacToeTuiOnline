package page

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
)

type JoinRoomForm struct {
	roomId          string
	maxRoomIdDigits int
	pageManager     *PageManager
	msg             string
}

func NewJoinRoomForm(pm *PageManager) *JoinRoomForm {
	return &JoinRoomForm{
		maxRoomIdDigits: 4,
		roomId:          "",
		pageManager:     pm,
		msg:             "",
	}
}

func (m *JoinRoomForm) Init() {
}

func (m *JoinRoomForm) Update(msg pageMsg.PageMsg) Command {
	m.msg = ""
	switch msg := msg.(type) {
	case pageMsg.KeyMsg:
		switch msg.Key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return QuitCommand
		case keyboard.KeyBackspace:
			m.deleteChar()
		case keyboard.KeyEnter:
			if len(m.roomId) != m.maxRoomIdDigits {
				m.msg = "Room id must be 4 digits long"
			} else {
				m.pageManager.ToGameRoom(m.roomId)
			}
		}
		if msg.Char >= '0' && msg.Char <= '9' {
			m.writeChar(msg.Char)
		}
	}
	return NoneCommand
}

func (m *JoinRoomForm) writeChar(char rune) {
	if len(m.roomId) >= m.maxRoomIdDigits {
		return
	}
	m.roomId += string(char)
}

func (m *JoinRoomForm) deleteChar() {
	if len(m.roomId) == 0 {
		return
	}
	m.roomId = m.roomId[:len(m.roomId)-1]
}

func (m *JoinRoomForm) Render() {
	s := "AddClient room\n\n"
	s += fmt.Sprintf("Room id (4 digits number): %s_ \n", m.roomId)
	s += fmt.Sprintf("%s                                \n", m.msg)
}
