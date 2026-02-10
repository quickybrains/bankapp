package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/http/common"
)

func respondWithStatus(w http.ResponseWriter, status int, data []byte) error {
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("http response writer write, err=%w", err)
	}

	return nil
}

func respondWithStatusAndError(w http.ResponseWriter, status int, resErr error) error {
	w.WriteHeader(status)

	errResp := &common.ErrorResp{
		Message: resErr.Error(),
	}

	data, err := json.Marshal(errResp)
	if err != nil {
		return fmt.Errorf("json marshal, err=%w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("http response writer write, err=%w", err)
	}

	return nil
}

func handleError(err error, w http.ResponseWriter, data []byte) error {
	if domain.IsDomainError(err) {
		return respondWithStatusAndError(w, http.StatusBadRequest, err)
	}

	return respondWithStatus(w, http.StatusInternalServerError, data)
}
