package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const keycloakHealthCheckPath = "/auth/realms/cronpad"
const keycloakTimeout = 2 * time.Minute
const keycloakExpectedStatus = 200

func keycloakHealthCheck(keycloakUrl string) error {
	url := keycloakUrl + keycloakHealthCheckPath

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode == keycloakExpectedStatus {
		return nil
	}

	log.Printf("[INFO] response from keycloak server (%v) is: %v\n", url, resp.Status)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(keycloakTimeout)
	for {
		select {
		case <-timeoutExceeded:
			return fmt.Errorf("keycloak connection failed after %s timeout", keycloakTimeout)

		case <-ticker.C:
			if err != nil {
				return err
			}
			if resp.StatusCode == keycloakExpectedStatus {
				return nil
			}

			log.Printf("[INFO] response from keycloak server (%v) is: %v\n", url, resp.Status)
		}
	}
}
