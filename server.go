package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MTRNord/cachet_go"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/go-pagerduty/webhookv3"
	"github.com/gorilla/mux"
	"github.com/maniartech/gotime"
)

var (
	cachetURL    = os.Getenv("CACHET_URL")
	cachetKey    = os.Getenv("CACHET_KEY")
	pagerdutyKey = os.Getenv("PAGERDUTY_KEY")
	secret       = os.Getenv("WEBHOOK_SECRET")
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/webhook", handler).Methods("POST")
	router.Use(mux.CORSMethodMiddleware(router))
	http.Handle("/", router)

	go func() {
		client := pagerduty.NewClient(pagerdutyKey)

		cachetClient, _ := cachet_go.NewClient(cachetURL, nil)
		cachetClient.Authentication.SetTokenAuth(cachetKey)
		_, status, err := cachetClient.General.Ping()
		if err != nil || status.StatusCode != 200 {
			log.Println(err)
			return
		}

		// Fethc maintenance windows every 30 minutes
		for {
			maintenanceWindows, err := fetchMaintenanceWindows(client)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("Fetched %d maintenance windows\n", len(maintenanceWindows))

			// Set them in cachet as schedules if they don't exist

			// Get all schedules
			schedules, _, err := cachetClient.Schedules.GetAll(&cachet_go.SchedulesQueryParams{
				QueryOptions: cachet_go.QueryOptions{
					PerPage: 100,
				},
			})
			if err != nil {
				log.Println(err)
				continue
			}

			// Create a map of schedules by start and end time
			schedulesMap := map[string]*cachet_go.Schedule{}
			for _, s := range schedules.Schedules {
				schedulesMap[s.ScheduledAt+"_"+s.CompletedAt] = &s
			}

			// Load Europe/Berlin timezone
			loc, err := time.LoadLocation("Europe/Berlin")
			if err != nil {
				log.Println(err)
				continue
			}

			// Create a map of maintenance windows by start and end time if they don't exist
			for _, mw := range maintenanceWindows {
				// Convert the pagerduty time in the format "2015-11-09T20:00:00-05:00" to "Y-m-d H:i:sO".
				// Keep in mind that cachet is in Europe/Berlin timezone
				// Keep in mind that pagerduty is giving us a string

				// Parse the start and end time
				startTime, err := time.Parse(time.RFC3339, mw.StartTime)
				if err != nil {
					log.Println(err)
					continue
				}
				endTime, err := time.Parse(time.RFC3339, mw.EndTime)
				if err != nil {
					log.Println(err)
					continue
				}

				// If the endTime is before now, skip
				if endTime.Before(time.Now()) {
					continue
				}

				// Convert the start and end time to Europe/Berlin timezone
				startTime = startTime.In(loc)
				endTime = endTime.In(loc)

				// Convert the start and end time to the format "Y-m-d H:i:sO"
				mw.StartTime = gotime.Format(startTime, "yyyy-mm-dd hhh:ii")
				mw.EndTime = gotime.Format(endTime, "yyyy-mm-dd hhh:ii")
				// Also convert it with space instead of + for the lookup table
				mwStart := gotime.Format(startTime, "yyyy-mm-dd hhh:ii:ss")
				mwEnd := gotime.Format(endTime, "yyyy-mm-dd hhh:ii:ss")

				// If there is no summary then we want to set it to "Maintenance"
				if mw.Description == "" {
					mw.Description = "Maintenance"
				}

				log.Printf("Maintenance window (%s): %s - %s\n", mw.Summary, mw.StartTime, mw.EndTime)

				var components []cachet_go.Component
				if mw.Services != nil {
					for _, s := range mw.Services {
						component, err := findComponentByName(cachetClient, s.Summary)
						if err != nil {
							log.Println(err)
							continue
						}
						if component != nil {
							components = append(components, *component)
						}
					}
				}

				if _, ok := schedulesMap[mwStart+"_"+mwEnd]; !ok {
					newSchedule := &cachet_go.Schedule{
						Name:        mw.Description,
						Message:     mw.Description,
						Status:      "0",
						ScheduledAt: mw.StartTime,
						CompletedAt: mw.EndTime,
						Components:  components,
					}
					// Print the schedule to the console
					log.Printf("%+v\n", newSchedule)

					_, _, err := cachetClient.Schedules.Create(newSchedule)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}

			time.Sleep(30 * time.Minute)
		}
	}()

	log.Println("Listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

func findComponentByName(client *cachet_go.Client, name string) (*cachet_go.Component, error) {
	components, _, err := client.Components.GetAll(&cachet_go.ComponentsQueryParams{
		QueryOptions: cachet_go.QueryOptions{
			PerPage: 100,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, c := range components.Components {
		if c.Name == name {
			return &c, nil
		}
	}

	return nil, nil
}

func fetchMaintenanceWindows(client *pagerduty.Client) ([]pagerduty.MaintenanceWindow, error) {
	ctx := context.TODO()
	maintenanceWindows := []pagerduty.MaintenanceWindow{}
	opts := pagerduty.ListMaintenanceWindowsOptions{}
	for {
		resp, err := client.ListMaintenanceWindowsWithContext(ctx, opts)
		if err != nil {
			return nil, err
		}
		maintenanceWindows = append(maintenanceWindows, resp.MaintenanceWindows...)
		if !resp.More {
			break
		}
		opts.Offset = resp.Offset
	}
	return maintenanceWindows, nil
}

func fetchIncident(client *cachet_go.Client, incidentID string) (*cachet_go.Incident, error) {
	incidents, _, err := client.Incidents.GetAll(&cachet_go.IncidentsQueryParams{
		QueryOptions: cachet_go.QueryOptions{
			PerPage: 100,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, i := range incidents.Incidents {
		v, ok := i.Meta.(map[string]interface{})
		if ok {
			v, ok := v["pagerduty"].(map[string]interface{})
			if ok {
				if v["incident_id"] == incidentID {
					return &i, nil
				}
			}
		}
	}

	return nil, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	client, _ := cachet_go.NewClient(cachetURL, nil)
	client.Authentication.SetTokenAuth(cachetKey)
	_, status, err := client.General.Ping()
	if err != nil || status.StatusCode != 200 {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = webhookv3.VerifySignature(r, secret)
	if err != nil {
		switch err {
		case webhookv3.ErrNoValidSignatures:
			log.Println("no valid signatures")
			w.WriteHeader(http.StatusUnauthorized)

		case webhookv3.ErrMalformedBody, webhookv3.ErrMalformedHeader:
			log.Println("malformed body or header")
			w.WriteHeader(http.StatusBadRequest)

		default:
			log.Println("internal server error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Println(err)

		fmt.Fprintf(w, "%v", err)
		return
	}

	// Read body as json
	bytedata, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Println("error reading body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error reading body")
		return
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bytedata))
	decodedJSON := WebhookMinimalEvent{}
	err = json.NewDecoder(r.Body).Decode(&decodedJSON)
	if err != nil {
		log.Println("error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error decoding json")
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bytedata))

	if decodedJSON.Event.EventType == "incident.triggered" {
		incident := WebhookIncidentTriggered{}
		err = json.NewDecoder(r.Body).Decode(&incident)
		if err != nil {
			log.Println("error decoding json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding json")
			return
		}
		log.Printf("%+v\n", incident)

		// Get description or fallback to title
		description := "This incident has been triggered automatically. Please stand by for updates."

		// If the urgency is high, set the status to degraded otherwise set the component status to normal
		componentStatus := cachet_go.ComponentStatusOperational
		if incident.Event.Data.Urgency == "high" {
			componentStatus = cachet_go.ComponentStatusPerformanceIssues
		}

		/*
			Create a meta map which looks like this:
			{
			    "meta": {
				    "pagerduty": {
				        "incident_id": "ABC"
				    }
				}
			}
		*/
		vars := map[string]map[string]string{
			"pagerduty": {
				"incident_id": incident.Event.Data.ID,
			},
		}

		// Create a new incident in cachet
		cachetIncident := &cachet_go.Incident{
			Name:            incident.Event.Data.Title,
			Message:         description,
			Status:          cachet_go.IncidentStatusInvestigating,
			ComponentID:     1,
			ComponentStatus: componentStatus,
			Meta:            vars,
		}
		_, cachetResp, err := client.Incidents.Create(cachetIncident)
		if err != nil {
			log.Println(err)
			// Print the cachetResp body as a string to the console
			bytedata, err := io.ReadAll(cachetResp.Body)
			r.Body.Close()
			if err != nil {
				log.Println("error reading cachetResp body")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println(string(bytedata))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else if decodedJSON.Event.EventType == "incident.acknowledged" {
		incident := WebhookIncidentAcknowledged{}
		err = json.NewDecoder(r.Body).Decode(&incident)
		if err != nil {
			log.Println("error decoding json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding json")
			return
		}
		log.Printf("%+v\n", incident)

		cachetIncident, err := fetchIncident(client, incident.Event.Data.ID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Fallthrough without error if the incident is not found
		if cachetIncident == nil {
			log.Println("incident not found")
			return
		}

		newCachetIncidentUpdate := &cachet_go.IncidentUpdate{
			Status:          cachet_go.IncidentStatusWatching,
			HumanStatus:     "Watching",
			Message:         "This incident has been acknowledged. We are currently investigating the issue. Thank you for your patience.",
			ComponentID:     cachetIncident.ComponentID,
			ComponentStatus: cachet_go.ComponentStatusPerformanceIssues,
		}
		_, cachetResp, err := client.IncidentUpdates.Create(cachetIncident.ID, newCachetIncidentUpdate)
		if err != nil {
			log.Println(err)
			// Print the cachetResp body as a string to the console
			bytedata, err := io.ReadAll(cachetResp.Body)
			r.Body.Close()
			if err != nil {
				log.Println("error reading cachetResp body")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println(string(bytedata))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else if decodedJSON.Event.EventType == "incident.resolved" {
		incident := WebhookIncidentResolved{}
		err = json.NewDecoder(r.Body).Decode(&incident)
		if err != nil {
			log.Println("error decoding json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding json")
			return
		}
		log.Printf("%+v\n", incident)

		cachetIncident, err := fetchIncident(client, incident.Event.Data.ID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Fallthrough without error if the incident is not found
		if cachetIncident == nil {
			log.Println("incident not found")
			return
		}

		newCachetIncidentUpdate := &cachet_go.IncidentUpdate{
			Status:          cachet_go.IncidentStatusFixed,
			HumanStatus:     "Fixed",
			Message:         "This incident has been resolved. Thank you for your patience.",
			ComponentID:     cachetIncident.ComponentID,
			ComponentStatus: cachet_go.ComponentStatusOperational,
		}

		_, cachetResp, err := client.IncidentUpdates.Create(cachetIncident.ID, newCachetIncidentUpdate)
		if err != nil {
			log.Println(err)
			// Print the cachetResp body as a string to the console
			bytedata, err := io.ReadAll(cachetResp.Body)
			r.Body.Close()
			if err != nil {
				log.Println("error reading cachetResp body")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println(string(bytedata))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("unknown event type")

		// Print the body as a string to the console
		log.Println(string(bytedata))

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unknown event type")
		return
	}

	fmt.Fprintf(w, "received signed webhook")
}
