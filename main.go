package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"triplecheck/gemini_manager"
	"triplecheck/instructor_manager"

	// Import the package containing the function 	_ "triplecheck/migrations"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type Person struct {
	ID   int
	Name string
}

func parseThinkingAndAnswer(input string) (thinking, answer string) {
	// Find the start and end of the thinking section
	log.Println("-")
	log.Println(input)

	thinkingStart := strings.Index(input, "<thinking>")
	thinkingEnd := strings.Index(input, "</thinking>")

	// Find the start and end of the answer section
	answerStart := strings.Index(input, "<answer>")
	answerEnd := strings.Index(input, "</answer>")

	// Extract thinking and answer if found
	if thinkingStart != -1 && thinkingEnd != -1 && thinkingStart < thinkingEnd {
		thinking = strings.TrimSpace(input[thinkingStart+len("<thinking>") : thinkingEnd])
	}

	if answerStart != -1 && answerEnd != -1 && answerStart < answerEnd {
		answer = strings.TrimSpace(input[answerStart+len("<answer>") : answerEnd])
	}

	return thinking, answer
}

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		e.Router.GET("/triplecheck", func(c echo.Context) error {
			name := c.PathParam("name")

			return c.JSON(http.StatusOK, map[string]string{"message": "Hello " + name})
		})
		return nil
	})

	app.OnRecordAfterCreateRequest("messages").Add(func(e *core.RecordCreateEvent) error {
		contentValue, ok := e.Record.Get("content").(string)
		if !ok {
			log.Println("Content field not found in the record")
			return nil
		}

		response, err := gemini_manager.CallAI(contentValue)
		if err != nil {
			log.Println("Error calling AI:", err)
			return err
		}

		thinking, answer := parseThinkingAndAnswer(response)

		// Update the record with thinking and answer
		// Update the record with thinking and answer
		e.Record.Set("thought_process", thinking)
		e.Record.Set("answer", answer)

		// Save the updated record
		err = app.Dao().SaveRecord(e.Record)
		if err != nil {
			log.Println("Error saving updated record:", err)
			return err
		}

		// Save the updated record
		err = app.Dao().SaveRecord(e.Record)
		if err != nil {
			log.Println("Error saving updated record:", err)
			return err
		}

		log.Println("Record updated successfully with thinking and answer")
		return nil
	})
	app.OnModelAfterCreate("messagesrrrr").Add(func(e *core.ModelEvent) error {
		log.Println(e.Model.TableName())
		log.Println(e.Model.GetId())
		log.Println()
		response, err := instructor_manager.CallAI("chat message")
		//kernelBootArgs, err := firecracker_manager.SetupFirecrackerEnv("./")
		if err != nil {
			fmt.Println("Error:", err)
		}
		log.Printf(response)
		return nil
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
