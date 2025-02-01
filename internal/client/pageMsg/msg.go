package pageMsg

import (
	"github.com/eiannone/keyboard"
	"github.com/sokmontrey/TicTacToeTuiOnline/payload"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type PageMsg any

type KeyMsg struct {
	Char rune
	Key  keyboard.Key
}

type OkMsg struct {
	Data any
}

type ErrMsg struct {
	Data any
}

// TODO reuse payload type here

type PlayerPositionMsg struct {
	PlayerId int
	Position pkg.Vec2
}

type TerminationMsg struct {
	WinnerId       int
	ConnectedCells map[pkg.Vec2]struct{}
}

type SyncMsg struct {
	PlayerPositions []payload.PlayerUpdate
	CellPositions   []payload.CellUpdate
	CurrentTurn     int
	CurrentPlayerId int
}

type BoardUpdateMsg struct {
	CellPos  pkg.Vec2
	CellId   int
	NextTurn int
}

func (k KeyMsg) ToMoveCode() payload.MoveCode {
	moveCode := payload.KeyToMoveCode(k.Key)
	if moveCode == payload.MoveCodeNone {
		moveCode = payload.CharToMoveCode(k.Char)
	}
	return moveCode
}

func NewKeyMsg(char rune, key keyboard.Key) KeyMsg {
	return KeyMsg{char, key}
}

func NewOkMsg(data any) OkMsg {
	return OkMsg{data}
}

func NewErrMsg(data any) ErrMsg {
	return ErrMsg{data}
}

func NewPositionMsg(playerId int, position pkg.Vec2) PlayerPositionMsg {
	return PlayerPositionMsg{playerId, position}
}

func NewBoardUpdateMsg(cellPos pkg.Vec2, cellId int, nextTurn int) BoardUpdateMsg {
	return BoardUpdateMsg{
		CellPos:  cellPos,
		CellId:   cellId,
		NextTurn: nextTurn,
	}
}

func NewSyncMsg(
	playerPositions []payload.PlayerUpdate,
	cellPositions []payload.CellUpdate,
	currentTurn int,
	currentPlayerId int,
) SyncMsg {
	return SyncMsg{
		PlayerPositions: playerPositions,
		CellPositions:   cellPositions,
		CurrentTurn:     currentTurn,
		CurrentPlayerId: currentPlayerId,
	}
}

func NewTerminationMsg(
	winnerId int,
	connectedCellsArr []pkg.Vec2,
) TerminationMsg {
	connectedCells := make(map[pkg.Vec2]struct{})
	for _, v := range connectedCellsArr {
		connectedCells[v] = struct{}{}
	}
	return TerminationMsg{
		WinnerId:       winnerId,
		ConnectedCells: connectedCells,
	}
}
