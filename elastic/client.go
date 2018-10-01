package elastic

import (
	"github.com/kamva/go-libs/exceptions"
	"github.com/kamva/go-libs/translation"
	"github.com/kamva/go-libs/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"gopkg.in/olivere/elastic.v3"
	"log"
	"os"
	"strings"
)

type Client struct {
	client    *elastic.Client
	issueCode string
	indexName string
	typeName  string
}

func (c *Client) Index(data Type) {
	c.setIndexAndType(data)
	c.checkIndex(data)

	_, err := c.client.Index().Index(c.indexName).Type(c.typeName).BodyJson(data.GetBody()).Do()

	c.handleError(err)
}

func (c *Client) setIndexAndType(data Type) {
	c.indexName = c.getIndexName(data.GetIndex())
	c.typeName = c.getTypeName(data)
}

func (c *Client) checkIndex(data Type) {
	exists, err := c.client.IndexExists(c.indexName).Do()

	if !exists {
		createIndex, err := c.client.CreateIndex(c.indexName).
			BodyString(data.GetIndex().GetMapping().Serialize()).Do()

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

func (c *Client) getIndexName(index Index) string {
	typeName := utils.ToSnake(utils.GetType(index))
	nameParts := strings.Split(typeName, "_")

	if nameParts[len(nameParts)] != "index" {
		c.handleError(errors.New("Invalid index name"))
	}

	return strings.Join(nameParts[:len(nameParts)-1], "_")
}

func (c *Client) getTypeName(data Type) string {
	typeName := utils.ToSnake(utils.GetType(data))
	nameParts := strings.Split(typeName, "_")

	if nameParts[len(nameParts)] != "type" {
		c.handleError(errors.New("Invalid type name"))
	}

	return strings.Join(nameParts[:len(nameParts)-1], "_")
}

func throwException(err error, exceptionCode string) {
	panic(exceptions.Exception{
		ResponseMessage: translation.Translate("internal_error"),
		Message:         err.Error(),
		Code:            exceptionCode,
		StatusCode:      iris.StatusInternalServerError,
	})
}

func NewClient(exceptionCode string, url string, options ...elastic.ClientOptionFunc) *Client {
	defaultOptions := []elastic.ClientOptionFunc{
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetErrorLog(log.New(os.Stderr, "[ELASTIC ERR] ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "[ELASTIC] ", log.LstdFlags)),
	}

	options = append(defaultOptions, options...)

	client, err := elastic.NewClient(options...)

	if err != nil {
		throwException(err, exceptionCode)
	}

	return &Client{
		client:    client,
		issueCode: exceptionCode,
	}
}
