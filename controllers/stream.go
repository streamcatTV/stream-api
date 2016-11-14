package controllers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strconv"
	"stream-api/data"
	"stream-api/models"
	"stream-api/services"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
)

// GetStream Gets stream.
func GetStream(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	stream, err := data.GetStreamByID(id)
	if err != nil {
		fmt.Println(err)
		resp := models.DoesNotExist{
			Code:    404,
			Message: "Stream does not exist",
		}
		return c.JSON(http.StatusNotFound, resp)
	}
	fmt.Println(stream)

	return c.JSON(http.StatusOK, stream)
}

// GetStreams Gets stream.
func GetStreams(c echo.Context) error {
	stream := data.GetStreams()
	fmt.Println(stream)

	return c.JSON(http.StatusOK, stream)
}

// CreateStream Creates a stream
func CreateStream(c echo.Context) error {

	// Bind model to json request body.
	s := new(models.Stream)
	if err := c.Bind(s); err != nil {
		return err
	}

	// Check if stream exists.
	if data.StreamExistsByName(s.StreamName) {
		resp := models.AlreadyExists{
			Code:    409,
			Message: "Stream name is taken.",
		}
		return c.JSON(http.StatusConflict, resp)
	}

	// Generate Key.
	streamKey := generateKey(10)
	s.StreamKey = streamKey

	fmt.Print("validating stream")
	// Validate Stream
	result, err := valid.ValidateStruct(s)
	if err != nil {
		fmt.Println(err)
		ve := models.ValidationError{
			ValidationErrors: valid.ErrorsByField(err),
		}
		return c.JSON(http.StatusBadRequest, ve)
	}
	fmt.Println(result)

	// Create record
	stream := data.CreateStream(s)
	return c.JSON(http.StatusCreated, stream)
}

// UpdateStream Updates a stream
func UpdateStream(c echo.Context) error {
	// Bind model to json request body.
	s := new(models.Stream)
	if err := c.Bind(s); err != nil {
		return err
	}

	id, _ := strconv.Atoi(c.Param("id"))

	stream, err := data.GetStreamByID(id)
	if err != nil {
		fmt.Println(err)
		resp := models.DoesNotExist{
			Code:    404,
			Message: "Stream does not exist",
		}
		return c.JSON(http.StatusNotFound, resp)
	}

	stream.Title = s.Title
	stream.Type = s.Type
	stream.Description = s.Description
	stream.Private = s.Private

	updatedStream := data.UpdateStreamByID(id, *stream)
	return c.JSON(http.StatusOK, updatedStream)
}

// DeleteStream deletes a stream
func DeleteStream(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	err := data.DeleteStreamByID(id)
	if err != nil {
		fmt.Println(err)
		resp := models.DoesNotExist{
			Code:    404,
			Message: "Stream does not exist",
		}
		return c.JSON(http.StatusNotFound, resp)
	}

	// resp := models.StreamDoesNotExist{
	// 	Code:    200,
	// 	Message: "Stream deleted.",
	// }
	return c.NoContent(http.StatusNoContent)
}

// AuthenticateStream Authenticates stream by checking StreamKey.
func AuthenticateStream(c echo.Context) error {
	streamName := c.QueryParam("name")
	streamKey := c.QueryParam("key")
	if data.ValidateStreamKey(streamName, streamKey) {
		return c.String(http.StatusOK, "OK")
	}
	return c.String(http.StatusForbidden, "Forbidden")
}

// GetAllStreamStats Gets all stream stats from rtmp server.
func GetAllStreamStats(c echo.Context) error {
	client, e := services.New()
	resp, e := client.GetRTMPStats()
	if e != nil {
		fmt.Print(e)
	}
	return c.JSON(http.StatusOK, resp)
}

func generateKey(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%x", b)
	return string(s)
}
