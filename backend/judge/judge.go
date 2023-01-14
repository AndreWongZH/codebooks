package judge

import (
	"bytes"
	"codebooks/constants"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubmitRequest struct {
	SourceCode string  `json:"source_code"`
	Language   string  `json:"language"`
	RoomID     *string `json:"room_id"`
}

type SubmitResponse struct {
	Stdout        *string          `json:"stdout"`
	Time          string           `json:"time"`
	Memory        int              `json:"memory"`
	Stderr        string           `json:"stderr"`
	Token         string           `json:"token"`
	CompileOutput string           `json:"compile_output"`
	Message       string           `json:"message"`
	Status        SubmissionStatus `json:"status"`
}

type SubmissionStatus struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

func Submit(c *gin.Context) {
	// read request body
	b0, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "error",
		})
	}
	var req SubmitRequest
	err = json.Unmarshal(b0, &req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":      "error",
			"description": "unmarshal json error",
		})
	}
	if req.RoomID == nil {
		fmt.Println("no room id")
	}

	// set request to judge0 parameters
	postBody, _ := json.Marshal(map[string]string{
		"source_code": req.SourceCode,
		"language_id": constants.LanguageMap[req.Language],
	})
	resp, _ := http.Post(constants.Judge0Endpoint+"/submissions/?base64_encoded=true&wait=true", "application/json", bytes.NewBuffer(postBody))

	// read response from judge0
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))

	var res SubmitResponse
	json.Unmarshal(b, &res)
	fmt.Println(res)
	if res.Stdout != nil {
		fmt.Println(*res.Stdout)
	}

	// set response
	c.JSON(http.StatusOK, &res)
}
