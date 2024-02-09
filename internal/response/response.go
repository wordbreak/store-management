package response

const (
	ResponseMessageOK = "ok"
)

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func New(code int, message string, data interface{}) Response {
	return Response{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
}
