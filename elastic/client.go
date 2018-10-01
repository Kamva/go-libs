package elastic

import (
	"github.com/kamva/go-libs/exceptions"
	"github.com/kamva/go-libs/translation"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"gopkg.in/olivere/elastic.v3"
)

type Client struct {
	client    *elastic.Client
	issueCode string
}

func (c *Client) Index(data Type) {
	c.checkIndex(data)

	_, err := c.client.Index().Index(data.GetIndexName()).
		Type(data.GetTypeName()).
		BodyJson(data.GetBody()).Do()

	c.handleError(err)
}

func (c *Client) checkIndex(data Type) {
	exists, err := c.client.IndexExists(data.GetIndexName()).Do()

	if !exists {
		createIndex, err := c.client.CreateIndex(data.GetIndexName()).
			BodyString(data.GetIndexMapping()).Do()

		c.handleError(err)

		if !createIndex.Acknowledged {
			throwException(errors.New("Index is not created."), c.issueCode)
		}
	}

	c.handleError(err)
}

func (c *Client) handleError(err error) {
	if err != nil {
		throwException(err, c.issueCode)
	}
}

func NewClient(exceptionCode string, options ...elastic.ClientOptionFunc) *Client {
	client, err := elastic.NewClient(options...)

	if err != nil {
		throwException(err, exceptionCode)
	}

	return &Client{
		client:    client,
		issueCode: exceptionCode,
	}
}

func throwException(err error, exceptionCode string) {
	panic(exceptions.Exception{
		ResponseMessage: translation.Translate("internal_error"),
		Message:         err.Error(),
		Code:            exceptionCode,
		StatusCode:      iris.StatusInternalServerError,
	})
}
