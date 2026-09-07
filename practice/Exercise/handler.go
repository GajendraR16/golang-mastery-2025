package main

import (
	"encoding/json"
	"net/http"
)

func (app *App) CreateURLResponseHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Decode the request body into URLCheckRequest.
		Pass the URLs to the concurrency implementation.
		Collect the processed results.
		Store the results in PostgreSQL in bulk.
		Send the response as JSON.
	*/
	var request URLCheckRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	results := app.BatchProcessor(request.URLs, r)
	err = app.Store.SaveResults(r.Context(), results)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonHandler(w, http.StatusCreated, results)

}

func (app *App) URLResultsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := app.Store.db.QueryContext(r.Context(), `
		SELECT url, status, status_code, response_time_ms, error
		FROM url_checks
		ORDER BY id
	`)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var results []*URLCheckResult

	for rows.Next() {
		var result URLCheckResult

		err := rows.Scan(
			&result.URL,
			&result.Status,
			&result.StatusCode,
			&result.ResponseTimeMs,
			&result.Error,
		)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}

		results = append(results, &result)
	}

	if err := rows.Err(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonHandler(w, http.StatusOK, URLCheckResponse{
		Results: results,
	})
}
