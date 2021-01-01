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

// getConnectionURL return mongo connection URI
func getConnectionURI(config *MongoDB) (URI string) {
	host := strings.Join(config.Hosts, ",")
	opt := strings.Join(config.Options, "?")
	if config.User == "" && config.Password == "" {
		return fmt.Sprintf("%v?%v", host, opt)
	}
	URI = fmt.Sprintf("mongodb+srv://%v:%v@%v/%v", config.User, config.Password, host, opt)

	return URI
}

// createSession return a new mongo session & transaction
func (m *MongoClient) createSession() (session mongo.Session) {
	session, err := m.Client.StartSession()
	if err != nil {
		log.Fatalln("Unable to init new session: ", err)
	}

	if err := session.StartTransaction(); err != nil {
		log.Fatalln("Unable to start transaction: ", err)
	}

	return session
}

// Create the list of document on collection
func (m *MongoClient) Create(databaseName, collectionName string, documents []interface{}) (interface{}, error) {

	var result interface{}
	session := m.createSession()
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {

		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.InsertMany(ctx, documents)
		if err != nil {
			log.Println("Unable to create document: ", err)
			return err
		}

		return nil
	}); err != nil {
		log.Println("Unable to execute mongo session: ", err)
		return nil, err
	}

	return result, nil
}

// Read documents from collection based on filter
func (m *MongoClient) Read(databaseName, collectionName string, filter interface{}, limit int64, dataModel reflect.Type) (interface{}, error) {

	var results interface{}
	session := m.createSession()
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {

		findOptions := options.Find()
		findOptions.SetLimit(limit)
		findOptions.SetSort(bson.D{primitive.E{Key: "_id", Value: 1}})

		collection := m.Client.Database(databaseName).Collection(collectionName)
		cur, err := collection.Find(ctx, filter, findOptions)
		defer cur.Close(ctx)
		if err != nil {
			log.Println("Unable to read document: ", err)
			return err
		}

		// Decode cursor
		dataModel := reflect.Zero(reflect.SliceOf(dataModel)).Type()
		results = reflect.New(dataModel).Interface()
		err = cur.All(ctx, results)
		if err != nil {
			log.Println("Unable to decode cursor: ", err)
			return err
		}

		return nil
	}); err != nil {
		log.Println("Unable to execute mongo session: ", err)
		return nil, err
	}

	return results, nil
}

// Update document with new value based on filter condition
func (m *MongoClient) Update(databaseName, collectionName string, filter, update interface{}) (interface{}, error) {

	var result interface{}
	session := m.createSession()
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {

		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.UpdateMany(ctx, filter, update)
		if err != nil {
			log.Println("Unable to update: ", err)
			return err
		}

		return nil
	}); err != nil {
		log.Println("Unable to execute mongo session: ", err)
		return nil, err
	}

	return result, nil
}

// Delete document based on filter condition
func (m *MongoClient) Delete(databaseName, collectionName string, filter interface{}) (interface{}, error) {

	var result interface{}
	session := m.createSession()
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {

		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.DeleteMany(ctx, filter)
		if err != nil {
			log.Println("Unable to delete: ", err)
			return err
		}

		return nil
	}); err != nil {
		log.Println("Unable to execute mongo session: ", err)
		return nil, err
	}

	return result, nil
}
