package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

const Now = firestore.ServerTimestamp

type Client struct {
	service *firestore.Client
}

type Collection string

func New(ctx context.Context) (*Client, error) {
	app, err := firebase.NewApp(ctx, nil) // set credential file path as GOOGLE_APPLICATION_CREDENTIALS
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		service: client,
	}, nil
}

func (c *Client) GetCollectionDocs(ctx context.Context, collection Collection) *firestore.DocumentIterator {
	return c.service.Collection(string(collection)).Documents(ctx)
}

func (c *Client) WhereDocumentsItr(ctx context.Context, collection, key, op string, value interface{}) *firestore.DocumentIterator {
	return c.service.Collection(string(collection)).Where(key, op, value).Documents(ctx)
}
