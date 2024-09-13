package krest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

/*
* Generic handler for http api. Implements basic CRUD operations.
 */
type Handler[T any] struct {
	service Service[T]
}

func NewHandler[T any](service Service[T]) *Handler[T] {
	return &Handler[T]{service: service}
}

func (h *Handler[T]) Get(w http.ResponseWriter, r *http.Request) {
	// Get the uuid from the url param  [GET /v1/tasks/{uuid}]
	// TODO: Remove dependency on chi, this is the only place we use it (for now).
	uuidStr := chi.URLParam(r, "uuid")
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the query parameters
	query, err := ParseResourceQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the meta query parameters
	metaQuery, err := ParseMetaQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the resource
	resource, err := h.service.Get(r.Context(), uuid, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	WriteResourceResponse(w, http.StatusOK, resource, query, metaQuery)

}

func (h *Handler[T]) List(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters
	query, err := ParseCollectionQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the meta query parameters
	metaQuery, err := ParseMetaQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the resources
	resources, err := h.service.List(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	WriteCollectionResponse(w, http.StatusOK, resources, len(resources), len(resources), query, metaQuery)
}

func (h *Handler[T]) Create(w http.ResponseWriter, r *http.Request) {
	// Parse the request body.
	var resource T
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the resource.
	createdResource, err := h.service.Create(r.Context(), resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response.
	WriteResourceResponse(w, http.StatusCreated, createdResource, ResourceQuery{}, MetaQuery{})
}

func (h *Handler[T]) Update(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler[T]) Delete(w http.ResponseWriter, r *http.Request) {
}

// Query paramers that can always be used.
type MetaQuery struct {
	// Meta fields to include in the response, "none" to exclude all meta field, "all" to include all meta fields.
	// links, query, expandable
	Meta []string
}

/*
* Extracts the resource query parameters from the http request.
 */
func ParseResourceQuery(r *http.Request) (ResourceQuery, error) {
	expandStr := r.URL.Query().Get("expand")

	var expand []string
	if expandStr != "" {
		expand = strings.Split(expandStr, ",")
	}

	return ResourceQuery{
		Expand: expand,
	}, nil
}

/*
*
 */
func ParseMetaQuery(r *http.Request) (MetaQuery, error) {
	var meta []string = []string{"all"}

	// Parse the meta query parameter.
	metaStr := r.URL.Query().Get("meta")
	if metaStr != "" {
		meta = strings.Split(metaStr, ",")
	}

	return MetaQuery{
		Meta: meta,
	}, nil
}

/*
*
 */
func ParseCollectionQuery(r *http.Request) (CollectionQuery, error) {
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
		return CollectionQuery{}, fmt.Errorf("invalid limit value: %s", limitStr)
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return CollectionQuery{}, fmt.Errorf("invalid offset value: %s", offsetStr)
	}

	var expand []string
	if expandStr != "" {
		expand = strings.Split(expandStr, ",")
	}

	return CollectionQuery{
		Limit:  limit,
		Offset: offset,
		Expand: expand,
	}, nil
}

func WriteCollectionResponse[T any](w http.ResponseWriter, code int, data []T, count int, total int, collectionQueryParams CollectionQuery, metaParams MetaQuery) {
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
func WriteResourceResponse[T any](w http.ResponseWriter, code int, data T, query ResourceQuery, metaParams MetaQuery) {
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

		// Reflect on T to find expandable fields (assumed tagged with `krest:"expandable"`)
		// TODO: Abstract this into a function.
		if all || slices.Contains(metaParams.Meta, "expandable") {
			/*allExpandableFields, err := ReflectExpandableFields[T]()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			expandableFields := make(map[string]interface{})
			for key := range allExpandableFields {
				if !slices.Contains(query.Expand, key) {
					expandableFields[key] = ""
				}
			}

			response["@expandable"] = expandableFields*/
			response["@expandable"] = map[string]interface{}{"TODO": "This was commented out in refactor.."}
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
