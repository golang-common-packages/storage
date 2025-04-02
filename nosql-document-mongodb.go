package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/golang-common-packages/hash"
)

// MongoClient manage all mongodb actions
type MongoClient struct {
	Client *mongo.Client
	Cancel context.CancelFunc
	Config *MongoDB
}

var (
	// mongoClientSessionMapping singleton pattern
	mongoClientSessionMapping = make(map[string]*MongoClient)
)

// newMongoDB init new instance
func newMongoDB(config *MongoDB) INoSQLDocument {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("Unable to marshal MongoDB configuration: ", err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentMongoSession := mongoClientSessionMapping[configAsString]
	if currentMongoSession == nil {
		currentMongoSession = &MongoClient{nil, nil, nil}

		// Establish MongoDB connection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(getConnectionURI(config)))
		if err != nil {
			cancel()
			log.Fatalln("Unable to connect to MongoDB: ", err)
		}

		// Check the connection status
		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			cancel()
			log.Fatalln("Unable to ping to MongoDB: ", err)
		}

		currentMongoSession.Client = client
		currentMongoSession.Cancel = cancel
		currentMongoSession.Config = config
		mongoClientSessionMapping[configAsString] = currentMongoSession
		log.Println("Connected to MongoDB")
	}

	return currentMongoSession
}

// getConnectionURI returns mongo connection URI
// It properly formats the connection string based on authentication requirements
func getConnectionURI(config *MongoDB) string {
	if config == nil {
		log.Println("Warning: MongoDB config is nil")
		return ""
	}
	
	host := strings.Join(config.Hosts, ",")
	opt := strings.Join(config.Options, "&")
	
	// Handle connection without authentication
	if config.User == "" && config.Password == "" {
		if opt != "" {
			return fmt.Sprintf("mongodb://%v?%v", host, opt)
		}
		return fmt.Sprintf("mongodb://%v", host)
	}
	
	// Connection with authentication
	if opt != "" {
		return fmt.Sprintf("mongodb+srv://%v:%v@%v/%v?%v", 
			config.User, 
			config.Password, 
			host, 
			config.DB, 
			opt)
	}
	
	return fmt.Sprintf("mongodb+srv://%v:%v@%v/%v", 
		config.User, 
		config.Password, 
		host, 
		config.DB)
}

// createSession returns a new mongo session & transaction
// It handles session creation and transaction initialization
func (m *MongoClient) createSession() (mongo.Session, error) {
	if m.Client == nil {
		return nil, fmt.Errorf("MongoDB client is not initialized")
	}
	
	session, err := m.Client.StartSession()
	if err != nil {
		log.Printf("Unable to init new MongoDB session: %v", err)
		return nil, err
	}

	if err := session.StartTransaction(); err != nil {
		log.Printf("Unable to start MongoDB transaction: %v", err)
		session.EndSession(ctx)
		return nil, err
	}

	return session, nil
}

// Create inserts a list of documents into the specified collection
func (m *MongoClient) Create(databaseName, collectionName string, documents []interface{}) (interface{}, error) {
	if len(documents) == 0 {
		return nil, fmt.Errorf("no documents to insert")
	}

	var result interface{}
	session, err := m.createSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.InsertMany(ctx, documents)
		if err != nil {
			log.Printf("Unable to create documents in %s.%s: %v", databaseName, collectionName, err)
			return err
		}

		return nil
	}); err != nil {
		log.Printf("Unable to execute mongo session for Create: %v", err)
		return nil, err
	}

	return result, nil
}

// Read retrieves documents from the specified collection based on filter
func (m *MongoClient) Read(databaseName, collectionName string, filter interface{}, limit int64, dataModel reflect.Type) (interface{}, error) {
	if dataModel == nil {
		return nil, fmt.Errorf("data model cannot be nil")
	}

	var results interface{}
	session, err := m.createSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		findOptions := options.Find()
		findOptions.SetLimit(limit)
		findOptions.SetSort(bson.D{primitive.E{Key: "_id", Value: 1}})

		collection := m.Client.Database(databaseName).Collection(collectionName)
		cur, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			log.Printf("Unable to read documents from %s.%s: %v", databaseName, collectionName, err)
			return err
		}
		defer cur.Close(ctx)

		// Decode cursor
		sliceType := reflect.Zero(reflect.SliceOf(dataModel)).Type()
		results = reflect.New(sliceType).Interface()
		err = cur.All(ctx, results)
		if err != nil {
			log.Printf("Unable to decode cursor: %v", err)
			return err
		}

		return nil
	}); err != nil {
		log.Printf("Unable to execute mongo session for Read: %v", err)
		return nil, err
	}

	return results, nil
}

// Update modifies documents in the specified collection based on filter
func (m *MongoClient) Update(databaseName, collectionName string, filter, update interface{}) (interface{}, error) {
	if filter == nil || update == nil {
		return nil, fmt.Errorf("filter and update cannot be nil")
	}

	var result interface{}
	session, err := m.createSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.UpdateMany(ctx, filter, update)
		if err != nil {
			log.Printf("Unable to update documents in %s.%s: %v", databaseName, collectionName, err)
			return err
		}

		return nil
	}); err != nil {
		log.Printf("Unable to execute mongo session for Update: %v", err)
		return nil, err
	}

	return result, nil
}

// Delete removes documents from the specified collection based on filter
func (m *MongoClient) Delete(databaseName, collectionName string, filter interface{}) (interface{}, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter cannot be nil")
	}

	var result interface{}
	session, err := m.createSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.DeleteMany(ctx, filter)
		if err != nil {
			log.Printf("Unable to delete documents from %s.%s: %v", databaseName, collectionName, err)
			return err
		}

		return nil
	}); err != nil {
		log.Printf("Unable to execute mongo session for Delete: %v", err)
		return nil, err
	}

	return result, nil
}
