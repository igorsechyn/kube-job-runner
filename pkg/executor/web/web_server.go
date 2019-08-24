package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"kube-job-runner/pkg/app"
	"kube-job-runner/pkg/app/job"

	"github.com/go-chi/chi"
)

type Server struct{}

func (webServer *Server) Run(ctx context.Context, app *app.App) {
	router := chi.NewRouter()

	router.Use(CreateHttResponseLogger(app.Reporter))

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Get("/execution/{jobId}", func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobId")
		details, err := app.JobService.GetJobDetails(jobID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(details)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(response)
		w.WriteHeader(http.StatusOK)
	})

	router.Post("/execute", func(writer http.ResponseWriter, request *http.Request) {
		var jobRequest job.Request
		rawBody, err := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			app.Reporter.Error("request.body.error", err, map[string]interface{}{})
			writer.WriteHeader(http.StatusInternalServerError)
		}

		err = json.Unmarshal(rawBody, &jobRequest)
		if err != nil {
			app.Reporter.Error("request.body.error", err, map[string]interface{}{})
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		jobId, err := app.JobService.SubmitJobCreationRequest(jobRequest)
		if err != nil {
			app.Reporter.Error("submit.job.error", err, map[string]interface{}{})
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Write([]byte(fmt.Sprintf(`{"callbackUrl":"/execution/%v"}`, jobId)))
		writer.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		app.Reporter.Error("server.startup.error", err, map[string]interface{}{})
	}
}
