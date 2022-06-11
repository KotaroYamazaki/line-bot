package linebot

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"unicode/utf8"

	"example.com/kanojo-bot/pkg/googlemap"
	"github.com/line/line-bot-sdk-go/linebot"
)

const maxTextWC = 60
const maxTitleWC = 40
const dots = "..."

func Webhook(w http.ResponseWriter, r *http.Request) {
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		http.Error(w, "Error init line bot", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	gmap, err := googlemap.New(os.Getenv("GOOGLE_API_KEY"))
	if err != nil {
		http.Error(w, "Error init google place API", http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		http.Error(w, "Error parse request", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	ctx := r.Context()
	for _, e := range events {
		switch e.Type {
		case linebot.EventTypeMessage:
			switch message := e.Message.(type) {
			case *linebot.TextMessage:
				query := message.Text
				res, err := gmap.TextSearch(ctx, &googlemap.TextSearchRequest{
					Query: fmt.Sprintf("静岡県浜松市 %s", query),
					Type:  googlemap.PlaceTypeCafe,
				})
				if err != nil {
					log.Fatal(err)
					continue
				}
				if len(res) == 0 {
					_, err = bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage("該当する結果がありませんでした")).Do()
					if err != nil {
						log.Println(err)
					}
					continue
				}
				max := 10
				ccs := make([]*linebot.CarouselColumn, 0, max)
				for i, result := range res {
					if i >= max {
						break
					}
					title := result.Name
					if utf8.RuneCountInString(title) > maxTitleWC {
						title = title[:(maxTitleWC-utf8.RuneCountInString(dots))] + dots
					}
					text := fmt.Sprintf("⭐ %.1f\n%s", result.Rating, result.Addr)
					if utf8.RuneCountInString(text) > maxTextWC {
						text = string([]rune(text)[:(maxTextWC-utf8.RuneCountInString(dots))]) + dots
					}
					cc := linebot.NewCarouselColumn(
						result.PhotoURL,
						title,
						text,
						linebot.NewURIAction("GoogleMapで開く", result.URL),
					).WithImageOptions("#FFFFFF")
					ccs = append(ccs, cc)
				}
				reply := linebot.NewTemplateMessage("カフェ一覧", linebot.NewCarouselTemplate(ccs...).
					WithImageOptions("rectangle", "cover"))
				_, err = bot.ReplyMessage(e.ReplyToken, reply).Do()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	fmt.Fprint(w, "ok")
}

type PlacesSearchResult struct {
	// FormattedAddress is the human-readable address of this place
	FormattedAddress string `json:"formatted_address,omitempty"`
	// Geometry contains geometry information about the result, generally including the
	// location (geocode) of the place and (optionally) the viewport identifying its
	// general area of coverage.

	// Name contains the human-readable name for the returned result. For establishment
	// results, this is usually the business name.
	Name string `json:"name,omitempty"`
	// Icon contains the URL of a recommended icon which may be displayed to the user
	// when indicating this result.
	Icon string `json:"icon,omitempty"`
	// PlaceID is a textual identifier that uniquely identifies a place.
	PlaceID string `json:"place_id,omitempty"`
	// Rating contains the place's rating, from 1.0 to 5.0, based on aggregated user
	// reviews.
	Rating float32 `json:"rating,omitempty"`
	// UserRatingsTotal contains total number of the place's ratings
	UserRatingsTotal int `json:"user_ratings_total,omitempty"`
	// Types contains an array of feature types describing the given result.
	Types []string `json:"types,omitempty"`
	// OpeningHours may contain whether the place is open now or not.
	// Photos is an array of photo objects, each containing a reference to an image.
	Photos []Photo `json:"photos,omitempty"`
	// PriceLevel is the price level of the place, on a scale of 0 to 4.
	PriceLevel int `json:"price_level,omitempty"`
	// Vicinity contains a feature name of a nearby location.
	Vicinity string `json:"vicinity,omitempty"`
	// PermanentlyClosed is a boolean flag indicating whether the place has permanently
	// shut down.
	PermanentlyClosed bool `json:"permanently_closed,omitempty"`
	// BusinessStatus is a string indicating the operational status of the
	// place, if it is a business.
	BusinessStatus string `json:"business_status,omitempty"`
	// ID is an identifier.
	ID string `json:"id,omitempty"`
}

type Photo struct {
	// PhotoReference is used to identify the photo when you perform a Photo request.
	PhotoReference string `json:"photo_reference"`
	// Height is the maximum height of the image.
	Height int `json:"height"`
	// Width is the maximum width of the image.
	Width int `json:"width"`
	// htmlAttributions contains any required attributions.
	HTMLAttributions []string `json:"html_attributions"`
}
