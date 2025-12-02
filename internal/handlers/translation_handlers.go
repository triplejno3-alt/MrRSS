package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleTranslateArticle translates an article's title.
func (h *Handler) HandleTranslateArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ArticleID  int64  `json:"article_id"`
		Title      string `json:"title"`
		TargetLang string `json:"target_language"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.TargetLang == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Translate the title
	translatedTitle, err := h.Translator.Translate(req.Title, req.TargetLang)
	if err != nil {
		log.Printf("Error translating article %d: %v", req.ArticleID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the article with the translated title
	if err := h.DB.UpdateArticleTranslation(req.ArticleID, translatedTitle); err != nil {
		log.Printf("Error updating article translation: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"translated_title": translatedTitle,
	})
}

// HandleClearTranslations clears all translated titles from the database.
func (h *Handler) HandleClearTranslations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.DB.ClearAllTranslations(); err != nil {
		log.Printf("Error clearing translations: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// HandleTranslateText translates any text to the target language.
// This is used for translating content, summaries, etc.
func (h *Handler) HandleTranslateText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Text       string `json:"text"`
		TargetLang string `json:"target_language"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Text == "" || req.TargetLang == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Translate the text
	translatedText, err := h.Translator.Translate(req.Text, req.TargetLang)
	if err != nil {
		log.Printf("Error translating text: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"translated_text": translatedText,
	})
}
