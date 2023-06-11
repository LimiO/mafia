package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"mafia/internal/db"
	"mafia/internal/helpers"
)

func passFilter(user *db.User) {
	user.PassHash = ""
}

func (c *Controller) getUser(ctx *gin.Context) {
	var req *db.User
	if err := ctx.BindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, "failed to parse req")
		return
	}
	if req.ID == "" {
		ctx.IndentedJSON(http.StatusBadRequest, "empty id")
		return
	}

	user, err := c.dbcontroller.GetUser(req.ID, passFilter)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, "failed to get user")
		return
	}
	if user == nil {
		ctx.IndentedJSON(http.StatusNotFound, "user not found")
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

func (c *Controller) mutateUser(ctx *gin.Context, mutationFunction func(*db.User) error) {
	var req *db.User

	if err := ctx.BindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, "failed to parse request")
		return
	}

	if req.ID == "" {
		ctx.IndentedJSON(http.StatusBadRequest, "ID cannot be empty")
		return
	}
	err := mutationFunction(req)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("failed to mutate user: %v", err))
		return
	}
	ctx.IndentedJSON(http.StatusOK, req)
}

func (c *Controller) newUser(ctx *gin.Context) {
	c.mutateUser(ctx, c.dbcontroller.CreateUser)
}

func (c *Controller) updateUser(ctx *gin.Context) {
	c.mutateUser(ctx, c.dbcontroller.UpdateUser)
}

func (c *Controller) deleteUser(ctx *gin.Context) {
	c.mutateUser(ctx, c.dbcontroller.DeleteUser)
}

func (c *Controller) getUsers(ctx *gin.Context) {
	body, err := c.dbcontroller.SelectUsers(passFilter)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("failed to select users: %v", err))
	}
	ctx.IndentedJSON(http.StatusOK, body)
}

func (c *Controller) pdfStats(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	request, ok := params["request"]
	if !ok {
		ctx.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title": "Bad request",
		})
		return
	}
	data, ok := c.pdfController.data[request[0]]
	if !ok {
		ctx.HTML(http.StatusNotFound, "index.html", gin.H{
			"title": "Not found",
		})
		return
	}
	switch data.status {
	case 0:
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Not yet",
		})
	case 1:
		filename := fmt.Sprintf("/infoserver/content/pdf/%s.pdf", request[0])
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Transfer-Encoding", "binary")
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		ctx.Header("Content-Type", "application/pdf")
		ctx.File(filename)
	case 2:
		ctx.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title": "Error",
		})
	}

}

func (c *Controller) pdfGen(ctx *gin.Context) {
	response := helpers.RandStringRunes(10)
	stats, err := c.dbcontroller.SelectStats()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title": fmt.Sprintf("failed to select stats: %v", err),
		})
		return
	}
	users, err := c.dbcontroller.SelectUsers(passFilter)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title": fmt.Sprintf("failed to select users: %v", err),
		})
		return
	}
	mergedData := mergeUserData(stats, users)

	request := &Request{
		Name: response,
		Data: mergedData,
	}
	c.pdfController.queue <- request

	ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("http://localhost:%d/getPdf?request=%s", c.cfg.Port, request.Name))
}

func (c *Controller) pdfGenUser(ctx *gin.Context) {
	var req *db.User

	if err := ctx.BindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, "failed to parse request")
		return
	}
	response := helpers.RandStringRunes(10)

	stat, err := c.dbcontroller.GetStats(req.ID)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title": fmt.Sprintf("failed to get stats: %v", err),
		})
		return
	}
	user, err := c.dbcontroller.GetUser(req.ID, passFilter)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title": fmt.Sprintf("failed to get user: %v", err),
		})
		return
	}
	mergedData := mergeUserData([]*db.Stats{stat}, []*db.User{user})
	request := &Request{
		Name: response,
		Data: mergedData,
	}

	c.pdfController.queue <- request

	ctx.IndentedJSON(http.StatusOK, fmt.Sprintf("http://localhost:%d/getPdf?request=%s", c.cfg.Port, request.Name))
}
