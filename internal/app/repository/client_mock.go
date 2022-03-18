package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Mock_Client struct {
	M_Connect     func() error
	M_Disconnect  func() error
	M_GetDatabase func() *mongo.Database
	M_GetClient   func() *mongo.Client
	M_Ping        func() error
}

func (r *Mock_Client) Connect() error {
	return r.M_Connect()
}

func (r *Mock_Client) Disconnect() error {
	return r.M_Disconnect()
}

func (r *Mock_Client) GetDatabase() *mongo.Database {
	return r.M_GetDatabase()
}

func (r *Mock_Client) GetClient() *mongo.Client {
	return r.M_GetClient()
}

func (r *Mock_Client) Ping() error {
	return r.M_Ping()
}
