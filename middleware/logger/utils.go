package logger

import "github.com/ximispot/woody"

func methodColor(method string, colors woody.Colors) string {
	switch method {
	case woody.MethodGet:
		return colors.Cyan
	case woody.MethodPost:
		return colors.Green
	case woody.MethodPut:
		return colors.Yellow
	case woody.MethodDelete:
		return colors.Red
	case woody.MethodPatch:
		return colors.White
	case woody.MethodHead:
		return colors.Magenta
	case woody.MethodOptions:
		return colors.Blue
	default:
		return colors.Reset
	}
}

func statusColor(code int, colors woody.Colors) string {
	switch {
	case code >= woody.StatusOK && code < woody.StatusMultipleChoices:
		return colors.Green
	case code >= woody.StatusMultipleChoices && code < woody.StatusBadRequest:
		return colors.Blue
	case code >= woody.StatusBadRequest && code < woody.StatusInternalServerError:
		return colors.Yellow
	default:
		return colors.Red
	}
}
