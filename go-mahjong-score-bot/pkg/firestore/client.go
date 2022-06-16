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

func (c *Client) Get(ctx context.Context, collection Collection, docID string, t interface{}) (interface{}, error) {
	ds, err := c.service.Collection(string(collection)).Doc(docID).Get(ctx)
	if err != nil {
		return nil, err
	}
	err = ds.DataTo(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *Client) Set(ctx context.Context, path Collection, docID string, data interface{}) error {
	_, err := c.service.Collection(string(path)).Doc(docID).Set(ctx, data)
	return err
}
