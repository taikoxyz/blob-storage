package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Server struct {
	cfg *Config
}

func NewServer(cfg *Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	fmt.Println("Server strated!")
	r := mux.NewRouter()
	r.HandleFunc("/getBlob", s.getBlobHandler).Methods("GET")

	http.Handle("/", r)
	return http.ListenAndServe(":27001", nil)
}

func (s *Server) getBlobHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["blobHash"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "Url Param 'blobHash' is missing", http.StatusBadRequest)
		return
	}

	blobHashes := keys[0]

	data, err := GetBlobData(s.cfg, blobHashes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response struct {
		Data []struct {
			Blob          string `json:"blob"`
			KzgCommitment string `json:"kzg_commitment"`
		} `json:"data"`
	}

	// Convert data to the correct type
	for _, d := range data {
		response.Data = append(response.Data, struct {
			Blob          string `json:"blob"`
			KzgCommitment string `json:"kzg_commitment"`
		}{Blob: d.Blob, KzgCommitment: d.KzgCommitment})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBlobData retrieves blob data from MongoDB based on blobHashes.
func GetBlobData(cfg *Config, blobHashes string) ([]struct {
	Blob          string `bson:"blob_data"`
	KzgCommitment string `bson:"kzg_commitment"`
}, error) {
	mongoClient, err := NewMongoDBClient(cfg.MongoDB)
	if err != nil {
		return nil, err
	}
	defer mongoClient.Close()

	collection := mongoClient.Client.Database(cfg.MongoDB.Database).Collection("blobs")

	cursor, err := collection.Find(context.Background(), bson.M{"blob_hash": blobHashes})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []struct {
		Blob          string `bson:"blob_data"`
		KzgCommitment string `bson:"kzg_commitment"`
	}
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
