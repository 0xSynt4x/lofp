// Package hub provides cross-machine coordination for multiplayer state
// using MongoDB Change Streams as a message bus and presence tracker.
package hub

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// PlayerPresence represents an online player across any machine.
type PlayerPresence struct {
	MachineID  string    `bson:"machineId"`
	FirstName  string    `bson:"firstName"`
	FullName   string    `bson:"fullName"`
	RoomNumber int       `bson:"roomNumber"`
	Race       int       `bson:"race"`
	RaceName   string    `bson:"raceName"`
	Position   int       `bson:"position"`
	IsGM       bool      `bson:"isGM"`
	GMHat      bool      `bson:"gmHat"`
	GMHidden   bool      `bson:"gmHidden"`
	GMInvis    bool      `bson:"gmInvis"`
	Hidden     bool      `bson:"hidden"`
	UpdatedAt  time.Time `bson:"updatedAt"`
}

// Event represents a cross-machine broadcast event.
type Event struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	MachineID     string        `bson:"machineId"`
	Type          string        `bson:"type"` // room_broadcast, global_broadcast, send_to_player, gm_broadcast
	RoomNumber     int            `bson:"roomNumber,omitempty"`
	TargetPlayer   string         `bson:"targetPlayer,omitempty"`
	ExcludePlayers []string       `bson:"excludePlayers,omitempty"`
	Messages       []string       `bson:"messages,omitempty"`
	Data           bson.M         `bson:"data,omitempty"` // opaque payload for room_state_change etc.
	CreatedAt      time.Time      `bson:"createdAt"`
}

// DeliveryFunc is called by the hub to deliver messages to local WebSocket connections.
type DeliveryFunc func(event *Event)

// Hub coordinates multiplayer state across machines via MongoDB.
type Hub struct {
	db         *mongo.Database
	machineID  string
	presence   *mongo.Collection
	events     *mongo.Collection
	deliverFn  DeliveryFunc

	// Local cache of all online players (refreshed via change stream)
	mu           sync.RWMutex
	allPlayers   []PlayerPresence
	cacheExpiry  time.Time
}

// New creates a new Hub. If db is nil, operates in local-only mode.
func New(db *mongo.Database, machineID string) *Hub {
	h := &Hub{
		db:        db,
		machineID: machineID,
	}
	if db != nil {
		h.presence = db.Collection("presence")
		h.events = db.Collection("events")
		h.setupIndexes()
	}
	return h
}

func (h *Hub) setupIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TTL index on presence — stale entries expire after 30 seconds
	h.presence.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "updatedAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(30),
	})
	h.presence.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "machineId", Value: 1}, {Key: "firstName", Value: 1}},
	})

	// TTL index on events — auto-clean after 60 seconds
	h.events.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "createdAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(60),
	})
}

// SetDeliveryFunc sets the callback for delivering remote events to local connections.
func (h *Hub) SetDeliveryFunc(fn DeliveryFunc) {
	h.deliverFn = fn
}

// Start begins watching for remote events via Change Stream and runs heartbeat.
func (h *Hub) Start() {
	if h.db == nil {
		return
	}
	go h.watchEvents()
	go h.heartbeatLoop()
}

// Stop cleans up this machine's presence entries.
func (h *Hub) Stop() {
	if h.presence == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	h.presence.DeleteMany(ctx, bson.M{"machineId": h.machineID})
}

// RegisterPlayer adds/updates a player in the cross-machine presence table.
func (h *Hub) RegisterPlayer(firstName, fullName string, roomNumber, race int, raceName string, position int, isGM, gmHat, gmHidden, gmInvis, hidden bool) {
	if h.presence == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"machineId": h.machineID, "firstName": firstName}
	update := bson.M{"$set": PlayerPresence{
		MachineID:  h.machineID,
		FirstName:  firstName,
		FullName:   fullName,
		RoomNumber: roomNumber,
		Race:       race,
		RaceName:   raceName,
		Position:   position,
		IsGM:       isGM,
		GMHat:      gmHat,
		GMHidden:   gmHidden,
		GMInvis:    gmInvis,
		Hidden:     hidden,
		UpdatedAt:  time.Now(),
	}}
	opts := options.UpdateOne().SetUpsert(true)
	h.presence.UpdateOne(ctx, filter, update, opts)
	h.invalidateCache()
}

// UpdatePlayerRoom updates just the room number for a player.
func (h *Hub) UpdatePlayerRoom(firstName string, roomNumber int) {
	if h.presence == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	h.presence.UpdateOne(ctx,
		bson.M{"machineId": h.machineID, "firstName": firstName},
		bson.M{"$set": bson.M{"roomNumber": roomNumber, "updatedAt": time.Now()}},
	)
	h.invalidateCache()
}

// UnregisterPlayer removes a player from cross-machine presence.
func (h *Hub) UnregisterPlayer(firstName string) {
	if h.presence == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	h.presence.DeleteOne(ctx, bson.M{"machineId": h.machineID, "firstName": firstName})
	h.invalidateCache()
}

// AllOnlinePlayers returns all players across all machines (cached).
func (h *Hub) AllOnlinePlayers() []PlayerPresence {
	if h.presence == nil {
		return nil
	}
	h.mu.RLock()
	if time.Now().Before(h.cacheExpiry) {
		players := h.allPlayers
		h.mu.RUnlock()
		return players
	}
	h.mu.RUnlock()

	// Refresh cache
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := h.presence.Find(ctx, bson.M{})
	if err != nil {
		return nil
	}
	var players []PlayerPresence
	cursor.All(ctx, &players)

	h.mu.Lock()
	h.allPlayers = players
	h.cacheExpiry = time.Now().Add(2 * time.Second)
	h.mu.Unlock()

	return players
}

func (h *Hub) invalidateCache() {
	h.mu.Lock()
	h.cacheExpiry = time.Time{}
	h.mu.Unlock()
}

// Publish sends an event to all machines (including self for remote players).
func (h *Hub) Publish(event *Event) {
	if h.events == nil {
		return
	}
	event.MachineID = h.machineID
	event.CreatedAt = time.Now()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := h.events.InsertOne(ctx, event); err != nil {
			log.Printf("hub: failed to publish event: %v", err)
		}
	}()
}

// watchEvents listens for events from OTHER machines and delivers them locally.
func (h *Hub) watchEvents() {
	for {
		if err := h.doWatch(); err != nil {
			log.Printf("hub: change stream error: %v, reconnecting in 2s", err)
			time.Sleep(2 * time.Second)
		}
	}
}

func (h *Hub) doWatch() error {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"operationType": "insert",
			"fullDocument.machineId": bson.M{"$ne": h.machineID},
		}}},
	}
	opts := options.ChangeStream().SetFullDocument(options.FullDocument("updateLookup"))
	ctx := context.Background()
	stream, err := h.events.Watch(ctx, pipeline, opts)
	if err != nil {
		return err
	}
	defer stream.Close(ctx)

	for stream.Next(ctx) {
		var change struct {
			FullDocument Event `bson:"fullDocument"`
		}
		if err := stream.Decode(&change); err != nil {
			continue
		}
		if h.deliverFn != nil {
			h.deliverFn(&change.FullDocument)
		}
	}
	return stream.Err()
}

// heartbeatLoop refreshes presence timestamps so TTL doesn't expire active players.
func (h *Hub) heartbeatLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		h.presence.UpdateMany(ctx,
			bson.M{"machineId": h.machineID},
			bson.M{"$set": bson.M{"updatedAt": time.Now()}},
		)
		cancel()
	}
}
