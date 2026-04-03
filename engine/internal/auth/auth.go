package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Account represents a user account linked to a Google identity.
type Account struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	GoogleID  string        `bson:"googleId" json:"googleId"`
	Email     string        `bson:"email" json:"email"`
	Name      string        `bson:"name" json:"name"`
	Picture   string        `bson:"picture" json:"picture"`
	IsAdmin   bool          `bson:"isAdmin" json:"isAdmin"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

// HasAnyAdmin checks if any admin account exists.
func (s *Service) HasAnyAdmin(ctx context.Context) (bool, error) {
	if s.db == nil {
		return false, fmt.Errorf("no database connection")
	}
	coll := s.db.Collection("accounts")
	count, err := coll.CountDocuments(ctx, bson.M{"isAdmin": true})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListAccounts returns all accounts, sorted by updatedAt descending.
func (s *Service) ListAccounts(ctx context.Context) ([]Account, error) {
	if s.db == nil {
		return nil, fmt.Errorf("no database connection")
	}
	coll := s.db.Collection("accounts")
	opts := options.Find().SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var accounts []Account
	if err := cursor.All(ctx, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

// SetAdmin sets or clears the admin flag on an account.
func (s *Service) SetAdmin(ctx context.Context, accountID string, isAdmin bool) (*Account, error) {
	if s.db == nil {
		return nil, fmt.Errorf("no database connection")
	}
	oid, err := bson.ObjectIDFromHex(accountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID")
	}
	coll := s.db.Collection("accounts")
	_, err = coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"isAdmin": isAdmin}})
	if err != nil {
		return nil, err
	}
	return s.GetAccount(ctx, accountID)
}

// Service handles authentication and account management.
type Service struct {
	db             *mongo.Database
	googleClientID string
	jwtSecret      []byte
	googleKeys     map[string]*rsa.PublicKey
	keysMu         sync.RWMutex
	keysExpiry     time.Time
}

// NewService creates a new auth service.
func NewService(db *mongo.Database, googleClientID, jwtSecret string) *Service {
	s := &Service{
		db:             db,
		googleClientID: googleClientID,
		jwtSecret:      []byte(jwtSecret),
		googleKeys:     make(map[string]*rsa.PublicKey),
	}
	if db != nil {
		// Create index on googleId for fast lookups
		coll := db.Collection("accounts")
		coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys:    bson.D{{Key: "googleId", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	}
	return s
}

// GoogleClaims represents the claims in a Google ID token.
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Subject       string `json:"sub"`
	jwt.RegisteredClaims
}

// googleCertsURL is the endpoint for Google's public keys.
const googleCertsURL = "https://www.googleapis.com/oauth2/v3/certs"

type jwksResponse struct {
	Keys []jwksKey `json:"keys"`
}

type jwksKey struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
}

// fetchGoogleKeys fetches and caches Google's public RSA keys.
func (s *Service) fetchGoogleKeys() error {
	s.keysMu.RLock()
	if time.Now().Before(s.keysExpiry) && len(s.googleKeys) > 0 {
		s.keysMu.RUnlock()
		return nil
	}
	s.keysMu.RUnlock()

	resp, err := http.Get(googleCertsURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Google certs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Google certs: %w", err)
	}

	var jwks jwksResponse
	if err := json.Unmarshal(body, &jwks); err != nil {
		return fmt.Errorf("failed to parse Google certs: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey)
	for _, k := range jwks.Keys {
		if k.Kty != "RSA" {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil {
			continue
		}
		n := new(big.Int).SetBytes(nBytes)
		e := 0
		for _, b := range eBytes {
			e = e*256 + int(b)
		}
		keys[k.Kid] = &rsa.PublicKey{N: n, E: e}
	}

	s.keysMu.Lock()
	s.googleKeys = keys
	s.keysExpiry = time.Now().Add(1 * time.Hour)
	s.keysMu.Unlock()

	return nil
}

// VerifyGoogleToken verifies a Google ID token and returns the claims.
func (s *Service) VerifyGoogleToken(idToken string) (*GoogleClaims, error) {
	if err := s.fetchGoogleKeys(); err != nil {
		return nil, err
	}

	// Parse the token header to get kid
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid token header")
	}
	var header struct {
		Kid string `json:"kid"`
	}
	json.Unmarshal(headerJSON, &header)

	s.keysMu.RLock()
	key, ok := s.googleKeys[header.Kid]
	s.keysMu.RUnlock()
	if !ok {
		// Key might have rotated — force refresh and retry
		s.keysMu.Lock()
		s.keysExpiry = time.Time{}
		s.keysMu.Unlock()
		if err := s.fetchGoogleKeys(); err != nil {
			return nil, fmt.Errorf("failed to refresh Google keys: %w", err)
		}
		s.keysMu.RLock()
		key, ok = s.googleKeys[header.Kid]
		s.keysMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("unknown signing key (kid=%s)", header.Kid)
		}
	}

	claims := &GoogleClaims{}
	token, err := jwt.ParseWithClaims(idToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Verify audience
	aud, err := claims.GetAudience()
	if err != nil || len(aud) == 0 {
		return nil, fmt.Errorf("missing audience")
	}
	validAud := false
	for _, a := range aud {
		if a == s.googleClientID {
			validAud = true
			break
		}
	}
	if !validAud {
		return nil, fmt.Errorf("invalid audience: token has %v, expected %q", aud, s.googleClientID)
	}

	// Verify issuer
	issuer, _ := claims.GetIssuer()
	if issuer != "accounts.google.com" && issuer != "https://accounts.google.com" {
		return nil, fmt.Errorf("invalid issuer: %s", issuer)
	}

	return claims, nil
}

// FindOrCreateAccount finds an existing account by Google ID or creates a new one.
func (s *Service) FindOrCreateAccount(ctx context.Context, claims *GoogleClaims) (*Account, error) {
	if s.db == nil {
		return nil, fmt.Errorf("no database connection")
	}
	coll := s.db.Collection("accounts")

	// Try to find existing
	var account Account
	err := coll.FindOne(ctx, bson.M{"googleId": claims.Subject}).Decode(&account)
	if err == nil {
		// Update profile info
		account.Email = claims.Email
		account.Name = claims.Name
		account.Picture = claims.Picture
		account.UpdatedAt = time.Now()
		coll.ReplaceOne(ctx, bson.M{"_id": account.ID}, account)
		return &account, nil
	}

	// Create new account
	account = Account{
		GoogleID:  claims.Subject,
		Email:     claims.Email,
		Name:      claims.Name,
		Picture:   claims.Picture,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result, err := coll.InsertOne(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		account.ID = oid
	}
	return &account, nil
}

// IssueJWT creates a signed JWT for the given account.
func (s *Service) IssueJWT(account *Account) (string, error) {
	claims := jwt.MapClaims{
		"sub":     account.ID.Hex(),
		"email":   account.Email,
		"name":    account.Name,
		"picture": account.Picture,
		"isAdmin": account.IsAdmin,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateJWT validates a JWT and returns the account ID.
// JWTClaims holds the decoded fields from a validated JWT.
type JWTClaims struct {
	AccountID string
	Email     string
	Name      string
}

func (s *Service) ValidateJWT(tokenStr string) (string, error) {
	claims, err := s.ValidateJWTFull(tokenStr)
	if err != nil {
		return "", err
	}
	return claims.AccountID, nil
}

func (s *Service) ValidateJWTFull(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid claims")
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return nil, fmt.Errorf("missing subject")
	}
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)
	return &JWTClaims{AccountID: sub, Email: email, Name: name}, nil
}

// GetAccount loads an account by ID.
func (s *Service) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	if s.db == nil {
		return nil, fmt.Errorf("no database connection")
	}
	oid, err := bson.ObjectIDFromHex(accountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID")
	}
	coll := s.db.Collection("accounts")
	var account Account
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&account); err != nil {
		return nil, fmt.Errorf("account not found")
	}
	return &account, nil
}
