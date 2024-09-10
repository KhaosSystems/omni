// This file contains helper functions for the REST API.
// Primarly, it's concered with creating easy-to-use constructs following the Khaos Collective REST API Specification 2024.
// See: https://www.notion.so/khaosgroup/Khaos-Collective-REST-API-Specification-2024-WIP-9b276e93b64c46ccb09d25e9757b3161
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type Collection struct {
	Count int           `json:"count"`
	Total int           `json:"total"`
	Items []interface{} `json:"items"`
}

// Query paramers that can always be used.
type MetaQueryParams struct {
	// Meta fields to include in the response, "none" to exclude all meta field, "all" to include all meta fields.
	Meta []string
}

func ParseMetaQueryParams(r *http.Request) (*MetaQueryParams, error) {
	var meta []string = []string{"all"}

	// Parse the meta query parameter.
	metaStr := r.URL.Query().Get("meta")
	if metaStr != "" {
		meta = strings.Split(metaStr, ",")
	}

	return &MetaQueryParams{
		Meta: meta,
	}, nil
}

type CollectionQueryParams struct {
	Limit  int
	Offset int
	Expand []string
}

func ParseCollectionQueryParams(r *http.Request) (*CollectionQueryParams, error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	expandStr := r.URL.Query().Get("expand")

	if limitStr == "" {
		limitStr = "50"
	}

	if offsetStr == "" {
		offsetStr = "0"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid limit value: %s", limitStr)
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return nil, fmt.Errorf("invalid offset value: %s", offsetStr)
	}

	var expand []string
	if expandStr != "" {
		expand = strings.Split(expandStr, ",")
	}

	return &CollectionQueryParams{
		Limit:  limit,
		Offset: offset,
		Expand: expand,
	}, nil
}

type ResourceQueryParams struct {
	Expand []string
}

func ParseResourceQueryParams(r *http.Request) (*ResourceQueryParams, error) {
	expandStr := r.URL.Query().Get("expand")

	var expand []string
	if expandStr != "" {
		expand = strings.Split(expandStr, ",")
	}

	return &ResourceQueryParams{
		Expand: expand,
	}, nil
}

// Generic function to write a collection response.
// A collection response could look like this:
//
//	{
//	  "@links": {
//	    "self": "/v1/tasks",
//	    "next": "/v1/tasks?offset=50&limit=50"
//	  },
//	  "@query": {
//	    "limit": 50,
//	    "offset": 0
//	  },
//	  "count": 50,
//	  "total": 100,
//	  "results": [
//	    { resource response },
//	 ]
//	}
func WriteCollectionResponse[T any](w http.ResponseWriter, code int, data []T, count int, total int, collectionQueryParams *CollectionQueryParams, metaParams *MetaQueryParams) {
	var response map[string]interface{} = make(map[string]interface{})

	// Add meta fields to the response.
	if !slices.Contains(metaParams.Meta, "none") {
		all := slices.Contains(metaParams.Meta, "all")

		if all || slices.Contains(metaParams.Meta, "links") {
			response["@links"] = map[string]interface{}{
				"self": "/v1/tasks",
			}

			if count < total {
				response["@links"].(map[string]interface{})["next"] = fmt.Sprintf("/v1/tasks?offset=%d&limit=%d", count, collectionQueryParams.Limit)
			}
		}

		if all || slices.Contains(metaParams.Meta, "query") {
			response["@query"] = map[string]interface{}{
				"limit":  collectionQueryParams.Limit,
				"offset": collectionQueryParams.Offset,
			}
		}
	}

	// Create the collection, and add it to the response.
	response["count"] = count
	response["total"] = total
	response["results"] = data

	// Write the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// Generic function to write a resource response.
// A resource response could look like this:
//
//		{
//		  "@links": {
//		    "self": "/v1/tasks"
//		  },
//		  "@query": {
//			"expand": ["author"]
//		 },
//		  "@expandable": { "owner": "/v1/users/123e4567-e89b-12d3-a456-426614174000" }, ,
//		  "uuid": "123e4567-e89b-12d3-a456-426614174000",
//		  "title": "My Task",
//		  "description": "This is a task."
//	   "author": {
//	     "@links": {...},
//	     "@query": {...},

//	     "uuid": "123e4567-e89b-12d3-a456-426614174000",
//	     "name": "John Doe"
//	   }
//		}
func WriteResourceResponse[T any](w http.ResponseWriter, code int, data T, metaParams *MetaQueryParams) {
	var response map[string]interface{} = make(map[string]interface{})

	// Add meta fields to the response.
	if !slices.Contains(metaParams.Meta, "none") {
		all := slices.Contains(metaParams.Meta, "all")

		if all || slices.Contains(metaParams.Meta, "links") {
			response["@links"] = map[string]interface{}{
				"self": "/v1/tasks",
			}
		}

		if all || slices.Contains(metaParams.Meta, "query") {
			response["@query"] = map[string]interface{}{}
		}
	}

	// Serialize the provided data and add it to the response.
	// Reflect on the data structure to allow generic serialization.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to serialize response data", http.StatusInternalServerError)
		return
	}

	// Deserialize back into a map to merge into the response.
	var dataMap map[string]interface{}
	err = json.Unmarshal(dataBytes, &dataMap)
	if err != nil {
		http.Error(w, "Failed to deserialize response data", http.StatusInternalServerError)
		return
	}

	// Merge the data into the response.
	for key, value := range dataMap {
		response[key] = value
	}

	// Write the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func RespondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	// if data is list

	// Write the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
