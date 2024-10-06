package main

import "time"

type WebhookMinimalEvent struct {
	Event struct {
		EventType    string `json:"event_type,omitempty"`
		ResourceType string `json:"resource_type,omitempty"`
	} `json:"event,omitempty"`
}

type WebhookIncidentTriggered struct {
	Event struct {
		ID           string    `json:"id,omitempty"`
		EventType    string    `json:"event_type,omitempty"`
		ResourceType string    `json:"resource_type,omitempty"`
		OccurredAt   time.Time `json:"occurred_at,omitempty"`
		Agent        struct {
			HTMLURL string `json:"html_url,omitempty"`
			ID      string `json:"id,omitempty"`
			Self    string `json:"self,omitempty"`
			Summary string `json:"summary,omitempty"`
			Type    string `json:"type,omitempty"`
		} `json:"agent,omitempty"`
		Client any `json:"client,omitempty"`
		Data   struct {
			ID          string    `json:"id,omitempty"`
			Type        string    `json:"type,omitempty"`
			Self        string    `json:"self,omitempty"`
			HTMLURL     string    `json:"html_url,omitempty"`
			Number      int       `json:"number,omitempty"`
			Status      string    `json:"status,omitempty"`
			IncidentKey string    `json:"incident_key,omitempty"`
			CreatedAt   time.Time `json:"created_at,omitempty"`
			Title       string    `json:"title,omitempty"`
			Description string    `json:"description,omitempty"`
			Service     struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"service,omitempty"`
			Assignees []struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"assignees,omitempty"`
			EscalationPolicy struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"escalation_policy,omitempty"`
			Teams    []any `json:"teams,omitempty"`
			Priority struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"priority,omitempty"`
			Urgency          string `json:"urgency,omitempty"`
			ConferenceBridge any    `json:"conference_bridge,omitempty"`
			ResolveReason    any    `json:"resolve_reason,omitempty"`
			IncidentType     any    `json:"incident_type,omitempty"`
		} `json:"data,omitempty"`
	} `json:"event,omitempty"`
}

type WebhookIncidentResolved struct {
	Event struct {
		ID           string    `json:"id,omitempty"`
		EventType    string    `json:"event_type,omitempty"`
		ResourceType string    `json:"resource_type,omitempty"`
		OccurredAt   time.Time `json:"occurred_at,omitempty"`
		Agent        struct {
			HTMLURL string `json:"html_url,omitempty"`
			ID      string `json:"id,omitempty"`
			Self    string `json:"self,omitempty"`
			Summary string `json:"summary,omitempty"`
			Type    string `json:"type,omitempty"`
		} `json:"agent,omitempty"`
		Client any `json:"client,omitempty"`
		Data   struct {
			ID          string    `json:"id,omitempty"`
			Type        string    `json:"type,omitempty"`
			Self        string    `json:"self,omitempty"`
			HTMLURL     string    `json:"html_url,omitempty"`
			Number      int       `json:"number,omitempty"`
			Status      string    `json:"status,omitempty"`
			IncidentKey string    `json:"incident_key,omitempty"`
			CreatedAt   time.Time `json:"created_at,omitempty"`
			Title       string    `json:"title,omitempty"`
			Service     struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"service,omitempty"`
			Assignees        []any `json:"assignees,omitempty"`
			EscalationPolicy struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"escalation_policy,omitempty"`
			Teams    []any `json:"teams,omitempty"`
			Priority struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"priority,omitempty"`
			Urgency          string `json:"urgency,omitempty"`
			ConferenceBridge any    `json:"conference_bridge,omitempty"`
			ResolveReason    any    `json:"resolve_reason,omitempty"`
			IncidentType     any    `json:"incident_type,omitempty"`
		} `json:"data,omitempty"`
	} `json:"event,omitempty"`
}

type WebhookIncidentAcknowledged struct {
	Event struct {
		ID           string    `json:"id,omitempty"`
		EventType    string    `json:"event_type,omitempty"`
		ResourceType string    `json:"resource_type,omitempty"`
		OccurredAt   time.Time `json:"occurred_at,omitempty"`
		Agent        struct {
			HTMLURL string `json:"html_url,omitempty"`
			ID      string `json:"id,omitempty"`
			Self    string `json:"self,omitempty"`
			Summary string `json:"summary,omitempty"`
			Type    string `json:"type,omitempty"`
		} `json:"agent,omitempty"`
		Client any `json:"client,omitempty"`
		Data   struct {
			ID          string    `json:"id,omitempty"`
			Type        string    `json:"type,omitempty"`
			Self        string    `json:"self,omitempty"`
			HTMLURL     string    `json:"html_url,omitempty"`
			Number      int       `json:"number,omitempty"`
			Status      string    `json:"status,omitempty"`
			IncidentKey string    `json:"incident_key,omitempty"`
			CreatedAt   time.Time `json:"created_at,omitempty"`
			Title       string    `json:"title,omitempty"`
			Service     struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"service,omitempty"`
			Assignees []struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"assignees,omitempty"`
			EscalationPolicy struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"escalation_policy,omitempty"`
			Teams    []any `json:"teams,omitempty"`
			Priority struct {
				HTMLURL string `json:"html_url,omitempty"`
				ID      string `json:"id,omitempty"`
				Self    string `json:"self,omitempty"`
				Summary string `json:"summary,omitempty"`
				Type    string `json:"type,omitempty"`
			} `json:"priority,omitempty"`
			Urgency          string `json:"urgency,omitempty"`
			ConferenceBridge any    `json:"conference_bridge,omitempty"`
			ResolveReason    any    `json:"resolve_reason,omitempty"`
			IncidentType     any    `json:"incident_type,omitempty"`
		} `json:"data,omitempty"`
	} `json:"event,omitempty"`
}
