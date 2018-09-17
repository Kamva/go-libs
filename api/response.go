package api

import "github.com/kataras/iris"

type Response struct {
	context iris.Context
}

func (r *Response) Ok(data interface{}, code string) {
	r.context.StatusCode(iris.StatusOK)
	r.context.JSON(iris.Map{
		"code": code,
		"data": data,
	})
}

func (r *Response) Created(message string, code string, data interface{}) {
	r.context.StatusCode(iris.StatusCreated)
	r.context.JSON(iris.Map{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func (r *Response) Accepted(message string, code string, data interface{}) {
	r.context.StatusCode(iris.StatusAccepted)
	r.context.JSON(iris.Map{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func (r *Response) BadRequest(message string, code string) {
	r.context.StatusCode(iris.StatusBadRequest)
	r.context.JSON(iris.Map{
		"code":    code,
		"message": message,
	})
}
func (r *Response) NotAcceptable(message string, code string, data ...interface{}) {
	r.context.StatusCode(iris.StatusNotAcceptable)
	r.context.JSON(iris.Map{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func (r *Response) InternalServerError(message string, code string) {
	r.context.StatusCode(iris.StatusInternalServerError)
	r.context.JSON(iris.Map{
		"code":    code,
		"message": message,
	})
}

func NewResponse(context iris.Context) Response {
	return Response{context: context}
}
