// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"math/rand"
	"time"
	"strconv"
	"strings"
	"os"
	"golang.org/x/net/context"
	"github.com/line/line-bot-sdk-go/linebot"
	"googlemaps.github.io/maps"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := handleText(message, event.ReplyToken, event.Source); err != nil {
					log.Print(err)
				}
			case *linebot.ImageMessage:
				if err := handleImage(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.VideoMessage:
				if err := handleVideo(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.AudioMessage:
				if err := handleAudio(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.LocationMessage:
				if err := handleLocation(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				if err := handleSticker(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			default:
				log.Printf("Unknown message: %v", message)
			}
		case linebot.EventTypeFollow:
			log.Print("Followed this bot: %v", event)
		case linebot.EventTypeUnfollow:
			log.Printf("Unfollowed this bot: %v", event)
		case linebot.EventTypeJoin:
			log.Print("Join: %v", event)
		case linebot.EventTypeLeave:
			log.Printf("Left: %v", event)
		case linebot.EventTypePostback:
			log.Print("Postback: %v", event)
		case linebot.EventTypeBeacon:
			log.Print("Beacon: %v", event)
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
	
	/*
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				action := strings.Contains(message.Text, "吃")
				target := strings.Contains(message.Text, "什麼")
				if !target {
					target = strings.Contains(message.Text, "啥")
				}
								
				if action && target {
					log.Print("SIVA: BINGO")
										
					i := random(1, 10)
					env := strconv.FormatInt(int64(i), 10)
					env = "SWFood"+env
					ans := os.Getenv(env)
					log.Print("SIVA: "+ans)
					
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ans)).Do(); err != nil {
						log.Print(err)
					}
				}
				//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
				//	log.Print(err)
				//}
			}
		}
	}
	*/
}

func handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	
	action := strings.Contains(message.Text, "吃")
	target := strings.Contains(message.Text, "什麼")
	if !target {
		target = strings.Contains(message.Text, "啥")
	}
								
	if action && target {
		log.Print("SIVA: BINGO")
										
		i := random(1, 10)
		env := strconv.FormatInt(int64(i), 10)
		env = "SWFood"+env
		ans := os.Getenv(env)
		log.Print("SIVA: "+ans)
					
		if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(ans)).Do(); err != nil {
			log.Print(err)
		}
	}
	
	return nil
}

func handleImage(message *linebot.ImageMessage, replyToken string) error {
	return nil
}

func handleVideo(message *linebot.VideoMessage, replyToken string) error {
	return nil
}

func handleAudio(message *linebot.AudioMessage, replyToken string) error {
	return nil
}

func handleLocation(message *linebot.LocationMessage, replyToken string) error {
	log.Print("handleLocation")
	
	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GMAP_KEY")))
    if err != nil {
        log.Fatalf("fatal error: %s", err)
    }
            
    ll := &maps.LatLng{
    	Lat: message.Latitude,
    	Lng: message.Longitude,
    }
    
    r := &maps.NearbySearchRequest{
        Location: ll,
        Type: "food",
        Language: "zh",
        //Region: "TW",
        OpenNow: true,
        RankBy: "distance",
    }
    
    resp, err := c.NearbySearch(context.Background(), r)
    if err != nil {
        log.Fatalf("fatal error: %s", err)
    } else {
    	len := len(resp.Results)
    	log.Print("length = %d", len)
    	
    	i := random(0, len - 1)
    	log.Print("target = %d", i)
    	log.Print("Name = %s", resp.Results[i].Name)
    	log.Print("FormattedAddress = %s", resp.Results[i].FormattedAddress)
    	log.Print("Vicinity = %s", resp.Results[i].Vicinity)
    	log.Print("Rating = %f", resp.Results[i].Rating)
    	ratingStr := FloatToString(float64(resp.Results[i].Rating))
    	
    	messsage := resp.Results[i].Name+"("+ratingStr+")\n"+resp.Results[i].Vicinity
    	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(messsage)).Do(); err != nil {
			log.Print(err)
		}
    }
	
	return nil
}

func handleSticker(message *linebot.StickerMessage, replyToken string) error {
	return nil
}

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', -1, 64)
}