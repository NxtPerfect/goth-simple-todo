package main

import (
	"context"
	"net/http"
	"simple-todo/layout"
	"simple-todo/types"
	"strconv"

	"github.com/labstack/echo/v4"
)

func main() {
	var id int64 = 0

	// Database
	tasks := []types.Task{types.Task{Id: id, Title: "Titel", Description: "Descreptan", Finished: false}}

	e := echo.New()
	// Get home page
	e.GET("/", func(c echo.Context) error {
		return layout.Home(tasks).Render(context.Background(), c.Response().Writer)
	})

	// Add task to db
	e.POST("/task/add", func(c echo.Context) error {
		title := c.Request().PostFormValue("title")
		description := c.Request().PostFormValue("description")
		id = id + 1
		tasks = append(tasks, types.Task{Id: id, Title: title, Description: description, Finished: false})
		c.Redirect(http.StatusOK, "/")
		return layout.Home(tasks).Render(context.Background(), c.Response().Writer)
	})

	// Get form for editing task
	e.GET("/task/edit", func(c echo.Context) error {
		_id := c.QueryParam("id")
		title := c.QueryParam("title")
		description := c.QueryParam("description")
		return layout.TaskEditForm(_id, title, description).Render(context.Background(), c.Response().Writer)
	})

	// Send edited task
	e.POST("/task/edit", func(c echo.Context) error {
		_id := c.Request().PostFormValue("id")
		title := c.Request().PostFormValue("title")
		description := c.Request().PostFormValue("description")
		__id, err := strconv.ParseInt(_id, 10, 64)
		if err != nil {
			panic(err)
		}
		for i, task := range tasks {
			if task.Id == __id {
				tasks[i] = types.Task{Id: __id, Title: title, Description: description, Finished: false}
				return layout.TaskComponent(__id, title, description, "false").Render(context.Background(), c.Response().Writer)
			}
		}
		return c.HTML(http.StatusNotFound, "")
	})

	// Delete task
	e.DELETE("/task/remove", func(c echo.Context) error {
		_id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
		if err != nil {
			panic(err)
		}
		if _id > int64(len(tasks)) && _id < 0 {
			return c.HTML(http.StatusNotFound, "")
		}
		for i, task := range tasks {
			if task.Id == _id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				return c.HTML(http.StatusOK, "")
			}
		}
		return c.HTML(http.StatusNotFound, "")
	})

	e.Logger.Fatal(e.Start(":8000"))
}
