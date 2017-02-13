package headlines

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type handlers struct {
	service service
}

func NewHeadlineHandler(headlineService service) handlers {
	return handlers{
		service: headlineService,
	}
}

func parseRequest(req *http.Request) (HeadlineInput, error) {
	input := HeadlineInput{}
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&input)
	return input, err
}

func (h *handlers) GetHeadlinesByUUID(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json")
	UUIDs, err := parseRequest(req)
	if err != nil {
		panic(err)
	}

	if len(UUIDs.UUIDs) == 0 {
		resp.WriteHeader(400)
		resp.Write([]byte("{\"error\": \"No list of UUIDs have been provided.\"}"))
		return
	}

	output, err := h.service.getHeadlinesByUUID(UUIDs.UUIDs)

	logrus.Debugf("GetHeadlinesByUUID: %v", output)

	if len(output) > 0 {
		enc := json.NewEncoder(resp)
		err = enc.Encode(output)
		if err != nil {
			resp.WriteHeader(503)
			resp.Write([]byte("{\"error\": \"Error creating response\"}"))
		}
	}
}

func (h *handlers) GetListHeadlines(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(req)
	UUID := vars["uuid"]

	output, err := h.service.getHeadlinesByList(UUID)

	logrus.Debugf("GetListHeadlines: %v", output)

	if len(output) > 0 {
		enc := json.NewEncoder(resp)
		err = enc.Encode(output)
		if err != nil {
			resp.WriteHeader(503)
			resp.Write([]byte("{\"error\": \"Error creating response\"}"))
		}
	}

}

func (h *handlers) GetFlashBriefing(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(req)
	UUID := vars["uuid"]

	output, err := h.service.getFlashBriefingForList(UUID)

	logrus.Debugf("GetFlashBriefing: %v", output)

	if len(output) > 0 {
		enc := json.NewEncoder(resp)
		err = enc.Encode(output)
		if err != nil {
			resp.WriteHeader(503)
			resp.Write([]byte("{\"error\": \"Error creating response\"}"))
		}
	}

}

func (h *handlers) GetConceptHeadlines(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(req)
	UUID := vars["uuid"]

	output, err := h.service.getHeadlinesByConcept(UUID)

	logrus.Debugf("GetConceptHeadlines: %v", output)

	if len(output) > 0 {
		enc := json.NewEncoder(resp)
		err = enc.Encode(output)
		if err != nil {
			resp.WriteHeader(503)
			resp.Write([]byte("{\"error\": \"Error creating response\"}"))
		}
	}
}
