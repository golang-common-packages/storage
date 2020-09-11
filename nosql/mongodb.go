package nosql

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/golang-common-packages/database/model"
)

// MongoClient manage all mongo API
type MongoClient struct {
	Client *mongo.Client
	Config *model.MongoDB
}

// NewMongoDB function return a new mongo client
func NewMongoDB(config *model.MongoDB) INoSQL {
	currentSession := &MongoClient{nil, nil}

	// Init Client options base on URI
	clientOptions := options.Client().ApplyURI(getConnectionURI(config))

	// Establish MongoDB connection
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error when try to connect to Mongodb server: ", err)
		panic(err)
	}

	// Check the connection status
	if err := client.Ping(ctx, nil); err != nil {
		log.Println("Can not ping to Mongodb server: ", err)
		panic(err)
	}

	currentSession.Client = client
	currentSession.Config = config
	log.Println("Connected to MongoDB Server")

	return currentSession
}

// getConnectionURL return mongo connection URI
func getConnectionURI(config *model.MongoDB) (URI string) {
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
		log.Println("Error when try to start session: ", err)
		panic(err)
	}

	if err := session.StartTransaction(); err != nil {
		log.Println("Error when try to start transaction: ", err)
		panic(err)
	}

	return session
}

// GetALL from collection with pagination
func (m *MongoClient) GetALL(
	databaseName,
	collectionName,
	lastID,
	pageSize string,
	dataModel reflect.Type) (results interface{}, err error) {

	session := m.createSession()
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		var f interface{}
		if lastID != "" {
			f, err = bsonGenerator(Match{"_id", GreaterThan, lastID})
			if err != nil {
				return err
			}
		}
		f, err = bsonGenerator(Match{})
		if err != nil {
			return err
		}

		if filter, ok := f.(bson.M); ok {
			limit, err := strconv.ParseInt(pageSize, 10, 64)
			if err != nil {
				return err
			}

			findOptions := options.Find()
			findOptions.SetLimit(limit)
			findOptions.SetSort(bson.D{primitive.E{Key: "_id", Value: 1}})

			collection := m.Client.Database(databaseName).Collection(collectionName)
			cur, err := collection.Find(ctx, filter, findOptions)
			defer cur.Close(ctx)
			if err != nil {
				return err
			}

			// Decode cursor
			dataModel := reflect.Zero(reflect.SliceOf(dataModel)).Type()
			results = reflect.New(dataModel).Interface()
			err = cur.All(ctx, results)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Println("Error in GetALL method: ", err)
		return nil, err
	}

	return results, nil
}

// GetByField base on field and value
func (m *MongoClient) GetByField(
	databaseName,
	collectionName,
	field,
	value string,
	dataModel reflect.Type) (result interface{}, err error) {

	session := m.createSession()
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		f, err := bsonGenerator(Match{field, Equal, value})
		if err != nil {
			return err
		}

		if filter, ok := f.(bson.M); ok {
			collection := m.Client.Database(databaseName).Collection(collectionName)
			SR := collection.FindOne(ctx, filter)
			if SR.Err() != nil {
				return SR.Err()
			}

			result = reflect.New(dataModel).Interface()
			err = SR.Decode(result)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Println("Error in GetByField method: ", err)
		return nil, err
	}

	return result, nil
}

// Create new record base on model
func (m *MongoClient) Create(
	databaseName,
	collectionName string,
	dataModel interface{}) (result interface{}, err error) {

	session := m.createSession()
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.InsertOne(ctx, dataModel)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error in Create method: ", err)
		return nil, err
	}

	return result, nil
}

// Update record with new value base on _id and model
func (m *MongoClient) Update(
	databaseName,
	collectionName,
	ID string,
	dataModel interface{}) (result interface{}, err error) {

	session := m.createSession()
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		ud, err := bsonGenerator(Set{Replaces, dataModel})
		if err != nil {
			return err
		}

		f, err := bsonGenerator(Match{"_id", Equal, ID})
		if err != nil {
			return err
		}

		update, ok := ud.(bson.M)
		if !ok {
			return errors.New("something wrong with bson update at Update method")
		}

		filter, ok := f.(bson.M)
		if !ok {
			return errors.New("something wrong with bson filter at Update method")
		}

		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error in Update method: ", err)
		return nil, err
	}

	return result, nil
}

// Delete record base on _id
func (m *MongoClient) Delete(
	databaseName,
	collectionName,
	ID string) (result interface{}, err error) {

	session := m.createSession()
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		f, err := bsonGenerator(Match{"_id", Equal, ID})
		if err != nil {
			return err
		}

		if filter, ok := f.(bson.M); ok {
			collection := m.Client.Database(databaseName).Collection(collectionName)
			result, err = collection.DeleteOne(ctx, filter)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Println("Error in Delete method: ", err)
		return nil, err
	}

	return result, nil
}

// MatchAndLookup ...
func (m *MongoClient) MatchAndLookup(
	databaseName,
	collectionName string,
	model MatchLookup,
	dataModel reflect.Type) (results interface{}, err error) {

	session := m.createSession()
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		p, err := bsonGenerator(model)
		if err != nil {
			return err
		}

		if pipeline, ok := p.([]bson.M); ok {
			collection := m.Client.Database(databaseName).Collection(collectionName)
			cur, err := collection.Aggregate(ctx, pipeline)
			defer cur.Close(ctx)
			if err != nil {
				return err
			}

			// Decode cursor
			dataModel := reflect.Zero(reflect.SliceOf(dataModel)).Type()
			results = reflect.New(dataModel).Interface()
			err = cur.All(ctx, results)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Println("Error in MatchAndLookup method: ", err)
		return nil, err
	}

	return results, nil
}

// bsonGenerator return bson format based on model
func bsonGenerator(
	rawModel interface{}) (interface{}, error) {

	// Generate MatchLookup pipeline []bson.M
	if model, ok := rawModel.(MatchLookup); ok {
		value := reflect.Indirect(reflect.ValueOf(model))
		fields := value.MapKeys()
		var pipeline []bson.M
		for _, field := range fields {
			f := field.Interface()

			// Generate match pipeline type [] Match
			if matches, ok := f.([]Match); ok {
				for _, match := range matches {
					var filter bson.M
					if match.Field == "_id" {
						id, err := primitive.ObjectIDFromHex(match.Value)
						if err != nil {
							return nil, err
						}

						filter["$match"] = bson.M{
							match.Field: bson.M{string(match.Operator): id},
						}
					} else {
						filter["$match"] = bson.M{
							match.Field: bson.M{string(match.Operator): match.Value},
						}
					}
					pipeline = append(pipeline, filter)
				}
			}

			// Generate lookup pipeline type [] Lookup
			if lookups, ok := f.([]Lookup); ok {
				for _, lookup := range lookups {
					var filter bson.M
					filter["$lookup"] = bson.M{
						"from":         lookup.From,
						"localField":   lookup.LocalField,
						"foreignField": lookup.ForeignField,
						"as":           lookup.As,
					}
					pipeline = append(pipeline, filter)
				}
			}
		}

		return pipeline, nil
	}

	// Generate Match type bson.M
	if match, ok := rawModel.(Match); ok {
		emptyMatch := Match{}
		Match := bson.M{}

		if match == emptyMatch {
			return Match, nil
		}

		if match.Field == "_id" {
			id, err := primitive.ObjectIDFromHex(match.Value)
			if err != nil {
				return nil, err
			}

			Match = bson.M{
				match.Field: bson.M{string(match.Operator): id},
			}

			return Match, nil
		}

		Match = bson.M{
			match.Field: bson.M{string(match.Operator): match.Value},
		}

		return Match, nil

	}

	// Generate Set type bson.M
	if set, ok := rawModel.(Set); ok {
		setOperator := bson.M{
			string(set.Operator): set.Data,
		}

		return setOperator, nil
	}

	return nil, errors.New("error in bsonGenerator function")
}
