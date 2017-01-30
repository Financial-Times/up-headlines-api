package headlines

import (
	"encoding/json"
	"net/http"
)

type handlers struct {
	service service
}

func NewHeadlineHandler(connStr string) handlers {
	return handlers{
		service: newHeadlineService(connStr),
	}
}

func (h *handlers) GetHeadlines(resp http.ResponseWriter, req *http.Request) {
	UUIDs := HeadlineInput{}
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&UUIDs)
	if err != nil {
		panic(err)
	}

	output := h.service.getHeadlines(UUIDs.UUIDs)

	if len(output) > 0 {

		resp.Header().Add("Content-Type", "application/json")

		enc := json.NewEncoder(resp)
		err = enc.Encode(output)
		if err != nil {
			panic(err)
		}
	}
}
