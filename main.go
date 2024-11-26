package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TeamMember struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Image    string `json:"image"`
	IdNo     string `json:"idno"`
	Email    string `json:"email"`
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
}

var db *mongo.Database

func main() {
	// Read MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	// MongoDB connection details
	const dbName = "mydatabase"
	const collectionName = "team"

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Access the database
	db = client.Database(dbName)
	log.Println("Connected to MongoDB!")

	// Serve static files
	fs := http.FileServer(http.Dir("./site"))
	http.Handle("/", fs)

	// Endpoint to fetch team data
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		// Fetch team data from the database
		collection := db.Collection(collectionName)
		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		var members []TeamMember

		if err := cursor.All(context.Background(), &members); err != nil {
			http.Error(w, "Failed to parse data", http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(members)
	})

	// API route to get a member's profile
	http.HandleFunc("/api/member/", func(w http.ResponseWriter, r *http.Request) {
		memberID := r.URL.Path[len("/api/member/"):]
		if memberID == "" {
			http.Error(w, "Member ID is required", http.StatusBadRequest)
			return
		}

		collection := db.Collection(collectionName)
		var member TeamMember
		err := collection.FindOne(context.Background(), bson.M{"id": memberID}).Decode(&member)
		if err != nil {
			http.Error(w, "Member not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(member)
	})

	// Run server
	log.Println("Server is running on port 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
