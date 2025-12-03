package opml

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	"MrRSS/internal/handlers/core"
	"MrRSS/internal/opml"
)

// HandleOPMLImport handles OPML file import.
func HandleOPMLImport(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleOPMLImport: ContentLength: %d", r.ContentLength)
	contentType := r.Header.Get("Content-Type")
	log.Printf("HandleOPMLImport: Content-Type: %s", contentType)

	var file io.Reader

	if strings.Contains(contentType, "multipart/form-data") {
		f, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error getting form file: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()
		log.Printf("HandleOPMLImport: Received file %s, size: %d", header.Filename, header.Size)

		if header.Size == 0 {
			http.Error(w, "Uploaded file is empty", http.StatusBadRequest)
			return
		}
		file = f
	} else {
		// Handle raw body upload
		file = r.Body
		defer r.Body.Close()
	}

	feeds, err := opml.Parse(file)
	if err != nil {
		log.Printf("Error parsing OPML: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		// Collect feed IDs for the newly imported feeds
		var feedIDs []int64
		for _, f := range feeds {
			feedID, err := h.Fetcher.ImportSubscription(f.Title, f.URL, f.Category)
			if err != nil {
				log.Printf("Error importing feed %s: %v", f.Title, err)
				continue
			}
			feedIDs = append(feedIDs, feedID)
		}

		// Fetch articles for the newly imported feeds with progress tracking
		if len(feedIDs) > 0 {
			h.Fetcher.FetchFeedsByIDs(context.Background(), feedIDs)
		}
	}()

	w.WriteHeader(http.StatusOK)
}

// HandleOPMLExport handles OPML file export.
func HandleOPMLExport(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	feeds, err := h.DB.GetFeeds()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := opml.Generate(feeds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=subscriptions.opml")
	w.Header().Set("Content-Type", "text/xml")
	w.Write(data)
}
