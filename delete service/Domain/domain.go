package domain

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var lock = &sync.Mutex{}

type single struct {
}

var Err error
var Client *mongo.Client
var Ctx context.Context
var cancel context.CancelFunc

// Implementing builder design pattern
type QueryBuilder interface {
	Db(database string) QueryBuilder
	Cols(cols string) QueryBuilder
	Find(id int32) QueryBuilder
	Create(item bson.D) QueryBuilder
	Update(filter bson.D, item bson.D) QueryBuilder
	Delete(filter bson.D) QueryBuilder
}

type MongoDbQueryBuilder struct {
	Database *mongo.Database
	Coll     *mongo.Collection
	Result   *mongo.Cursor
}

func (m *MongoDbQueryBuilder) Db(database string) QueryBuilder {
	Err = nil
	m.Database = Client.Database(database)
	return &MongoDbQueryBuilder{}
}
func (m *MongoDbQueryBuilder) Cols(cols string) QueryBuilder {
	Err = nil
	m.Coll = m.Database.Collection(cols)
	return &MongoDbQueryBuilder{}
}
func (m *MongoDbQueryBuilder) Find(id int32) QueryBuilder {
	Err = nil
	if id != 0 {
		filter := bson.D{{"Id", id}}
		m.Result, Err = m.Coll.Find(Ctx, filter)
		if Err != nil {
			fmt.Println("can't find the value", Err)
		}
	} else {
		m.Result, Err = m.Coll.Find(Ctx, nil)
		if Err != nil {
			fmt.Println("can't find the value", Err)
		}
	}
	return &MongoDbQueryBuilder{}
}

func (m *MongoDbQueryBuilder) Create(item bson.D) QueryBuilder {
	Err = nil
	a, erro := m.Coll.InsertOne(Ctx, item)
	Err = erro
	fmt.Println(a)
	if Err != nil {
		fmt.Println(Err)
	}
	return &MongoDbQueryBuilder{}
}
func (m *MongoDbQueryBuilder) Update(filter bson.D, item bson.D) QueryBuilder {
	m.Coll.DeleteOne(Ctx, filter)
	m.Coll.InsertOne(Ctx, item)
	return &MongoDbQueryBuilder{}
}
func (m *MongoDbQueryBuilder) Delete(filter bson.D) QueryBuilder {
	m.Coll.DeleteOne(Ctx, filter)
	return &MongoDbQueryBuilder{}
}

// Implementing singleton design pattern
func Connect_to_db() (*mongo.Client, context.Context) {
	if Client == nil {
		lock.Lock()
		defer lock.Unlock()
		if Client == nil {
			fmt.Println("Creating single instance now.")
			clientOptions := options.Client().
				ApplyURI("mongodb://localhost:27017")

			Ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
			//	defer cancel()
			Client, _ = mongo.Connect(Ctx, clientOptions)

			Ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
			//defer cancel()
			if err := Client.Ping(Ctx, nil); err != nil {
				fmt.Println(err)
				fmt.Errorf("failed to ping mongodb %w", err)
				return nil, nil
			}

		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return Client, Ctx
}
