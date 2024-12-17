package response

type (
	// Message is a type alias for translation message.
	Message map[string]string

	// Error holds the structure for any error of response.
	Error struct {
		UserMessage     string `json:"userMessage"`
		InternalMessage string `json:"internalMessage"`
		Code            int    `json:"code"`
		MoreInfo        string `json:"moreInfo"`
	}

	// HTTPResponse holds the default HTTP response structures.
	HTTPResponse struct {
		Status  int         `json:"status"`
		Message Message     `json:"message"`
		Data    interface{} `json:"data,omitempty"`
		Meta    interface{} `json:"meta,omitempty"`
		Errors  []Error     `json:"errors,omitempty"`
	}
)

var (
	// DefaultMessage http status: 200 - status ok.
	DefaultMessage = map[string]string{
		"id": "Berhasil",
		"en": "Success",
	}

	// DefaultErrorMessage http status: 400 - bad request.
	DefaultErrorMessage = map[string]string{
		"id": "Gagal",
		"en": "Failed",
	}

	// ResponseMessageNotFound http status: 404 - data not found.
	ResponseMessageNotFound = map[string]string{
		"id": "Data tidak ditemukan",
		"en": "Data not found",
	}

	// ResponseMessageMethodNotAllowed http status: 405 - method not allowed.
	ResponseMessageMethodNotAllowed = map[string]string{
		"id": "Permintaan tidak diizinkan",
		"en": "Method Not Allowed",
	}

	// ResponseMessageInternalServerError http status: 500 - internal server error.
	ResponseMessageInternalServerError = map[string]string{
		"id": "Terjadi kesalahan tak terduga. Silahkan coba lagi nanti",
		"en": "An unexpected error occurred. Please try again later",
	}

	// ResponseMessageServiceUnavailable http status: 503 - service unavailable.
	ResponseMessageServiceUnavailable = map[string]string{
		"id": "Server tidak tersedia. Silahkan coba lagi nanti",
		"en": "Service Unavailable. Please try again later",
	}

	// ResponseMessageUnprocessableEntity http status 422 - Unprocessable Entity
	ResponseMessageUnprocessableEntity = map[string]string{
		"id": "Entitas tidak dapat diproses",
		"en": "Unprocessable Entity",
	}
)
