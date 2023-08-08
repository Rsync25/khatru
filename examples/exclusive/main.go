package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/plugins"
	"github.com/fiatjaf/khatru/plugins/storage/lmdbn"
	"github.com/nbd-wtf/go-nostr"
)

func main() {
	relay := khatru.NewRelay()

	db := lmdbn.LMDBBackend{Path: "/tmp/exclusive"}
	os.MkdirAll(db.Path, 0755)
	if err := db.Init(); err != nil {
		panic(err)
	}

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)

	relay.RejectEvent = append(relay.RejectEvent, plugins.PreventTooManyIndexableTags(10))
	relay.RejectFilter = append(relay.RejectFilter, plugins.NoPrefixFilters, plugins.NoComplexFilters)

	relay.OnEventSaved = append(relay.OnEventSaved, func(ctx context.Context, event *nostr.Event) {
	})

	fmt.Println("running on :3334")
	http.ListenAndServe(":3334", relay)
}

func deleteStuffThatCanBeFoundElsewhere() {
}
