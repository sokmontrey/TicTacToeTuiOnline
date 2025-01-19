package network

type ResponseType int

const (
	ResponseTypeError ResponseType = iota
	ResponseTypeSuccess
)

type Response struct {
	Type ResponseType `json:"type"`
	Data any          `json:"data"`
}

func NewSuccessResponse(data any) Response {
	return Response{
		Type: ResponseTypeSuccess,
		Data: data,
	}
}

func NewErrorResponse(data any) Response {
	return Response{
		Type: ResponseTypeError,
		Data: data,
	}
}
