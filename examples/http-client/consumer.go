package main

import (
	"net/http"
	"net/url"
)

// registerConsumerRequest/-Response spiegeln den JSON-Vertrag von
// `POST /consumers` (`LH-FA-CON-001`, `internal/adapters/driving/http/registerconsumer.go`):
// beide Request-Felder Pflicht, `AlreadyRegistered` trägt die
// Idempotenz-Antwort ohne eigenen Statuscode.
type registerConsumerRequest struct {
	ConsumerID string `json:"consumer_id"`
	Name       string `json:"name"`
}

type registerConsumerResponse struct {
	ConsumerID        string `json:"consumer_id"`
	Name              string `json:"name"`
	AlreadyRegistered bool   `json:"already_registered"`
}

// registerConsumer ruft den `admin`-Endpunkt `POST /consumers` auf.
func registerConsumer(client *http.Client, cfg config) (registerConsumerResponse, error) {
	var resp registerConsumerResponse
	err := doRequestJSON(client, http.MethodPost, "http://"+cfg.addr+"/consumers", cfg.adminToken,
		registerConsumerRequest{ConsumerID: cfg.consumerID, Name: cfg.name}, &resp, http.StatusCreated)
	return resp, err
}

// acknowledgeConsumerRequest/-Response spiegeln `POST /consumers/acknowledge`
// (`LH-FA-CON-004`).
type acknowledgeConsumerRequest struct {
	ConsumerID string `json:"consumer_id"`
	SourceID   string `json:"source_id"`
	Offset     uint64 `json:"offset"`
}

type acknowledgeConsumerResponse struct {
	ConsumerID string `json:"consumer_id"`
	SourceID   string `json:"source_id"`
	Offset     uint64 `json:"offset"`
}

// acknowledgeConsumer ruft den `admin`-Endpunkt `POST /consumers/acknowledge`
// auf.
func acknowledgeConsumer(client *http.Client, cfg config) (acknowledgeConsumerResponse, error) {
	var resp acknowledgeConsumerResponse
	err := doRequestJSON(client, http.MethodPost, "http://"+cfg.addr+"/consumers/acknowledge", cfg.adminToken,
		acknowledgeConsumerRequest{ConsumerID: cfg.consumerID, SourceID: cfg.source, Offset: cfg.offset}, &resp, http.StatusOK)
	return resp, err
}

// consumerPositionResponse spiegelt `GET /consumers/position`
// (`LH-FA-CON-005`): `Acknowledged` unterscheidet die definierte
// Anfangsposition (kein Nachweis) von einer echten Bestätigung mit
// Offset 0.
type consumerPositionResponse struct {
	ConsumerID   string `json:"consumer_id"`
	SourceID     string `json:"source_id"`
	Offset       uint64 `json:"offset"`
	Acknowledged bool   `json:"acknowledged"`
}

// ConsumerPositionURL baut die Lese-Adresse von `GET /consumers/position`:
// ihr einziges Pflichtfeld ist `consumer_id`.
func ConsumerPositionURL(addr, consumerID string) string {
	u := url.URL{Scheme: "http", Host: addr, Path: "/consumers/position"}
	q := u.Query()
	q.Set("consumer_id", consumerID)
	u.RawQuery = q.Encode()
	return u.String()
}

// consumerPosition ruft den `reader`-Endpunkt `GET /consumers/position` auf.
func consumerPosition(client *http.Client, cfg config) (consumerPositionResponse, error) {
	var resp consumerPositionResponse
	err := doRequestJSON(client, http.MethodGet, ConsumerPositionURL(cfg.addr, cfg.consumerID), cfg.token, nil, &resp, http.StatusOK)
	return resp, err
}

// removeConsumerRequest/-Response spiegeln `POST /consumers/remove`
// (`LH-FA-CON-006`): `Removed` trägt den Idempotenz-Ausgang — ein nie
// registrierter Consumer meldet `false`, kein `404`.
type removeConsumerRequest struct {
	ConsumerID string `json:"consumer_id"`
}

type removeConsumerResponse struct {
	ConsumerID string `json:"consumer_id"`
	Removed    bool   `json:"removed"`
}

// removeConsumer ruft den `admin`-Endpunkt `POST /consumers/remove` auf.
func removeConsumer(client *http.Client, cfg config) (removeConsumerResponse, error) {
	var resp removeConsumerResponse
	err := doRequestJSON(client, http.MethodPost, "http://"+cfg.addr+"/consumers/remove", cfg.adminToken,
		removeConsumerRequest{ConsumerID: cfg.consumerID}, &resp, http.StatusOK)
	return resp, err
}
