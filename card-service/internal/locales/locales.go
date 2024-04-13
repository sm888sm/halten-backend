package locales

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	message.SetString(language.English, "Service listening on port %d", "Service listening on port %d")
	message.SetString(language.German, "Service listening on port %d", "Dienst hört auf Port %d")

	message.SetString(language.English, "Failed to listen: %v", "Failed to listen: %v")
	message.SetString(language.German, "Failed to listen: %v", "Fehler beim Zuhören: %v")

	message.SetString(language.English, "Failed to serve: %v", "Failed to serve: %v")
	message.SetString(language.German, "Failed to serve: %v", "Fehler beim Bedienen: %v")
	// Add more translations as needed...
}
