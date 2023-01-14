package judge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Submit(c *gin.Context) {
	postBody, _ := json.Marshal(map[string]string{
		"source_code": "cGFja2FnZSBtYWluCgppbXBvcnQgImZtdCIKCmZ1bmMgbWFpbigpIHsKCWZtdC5QcmludGxuKCJIZWxsbywg5LiW55WMIikKfQ==",
		"language_id": "60",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, _ := http.Post("http://localhost:2358/submissions/?base64_encoded=true&wait=true", "application/json", responseBody)
	fmt.Println(resp)
}
