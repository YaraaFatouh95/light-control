package dkron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"light-control/models"
	"net/http"
)

type dkronJob struct {
	Name           string            `json:"name"`
	Schedule       string            `json:"schedule"`
	Executor       string            `json:"executor"`
	ExecutorConfig map[string]string `json:"executor_config"`
}

func CreateDkronJob(command models.Command) error {

	jsonBytes, err := json.Marshal(command)
	if err != nil {
		return err
	}

	scheduleHour := command.ScheduledTime.Hour()
	scheduleMin := command.ScheduledTime.Minute()

	job := dkronJob{
		Name:     command.ID.String(),
		Schedule: fmt.Sprintf("%v %v * * * *", scheduleMin, scheduleHour),
		Executor: "http",
		ExecutorConfig: map[string]string{
			"url":     "http://localhost:8081/command/exec",
			"method":  "POST",
			"headers": `{"Content-Type":"application/json"}`,
			"body":    string(jsonBytes),
		},
	}

	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/v1/jobs", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("can not create dkron Job")
	}

	return nil
}

func DeleteDkronJob(command models.Command) error {

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/v1/jobs/%v", command.ID), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return fmt.Errorf("can not delete dkron Job")
	}

	return nil
}
