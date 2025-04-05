package tasks

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rebelnato/goAIRA/authorizations"
	"github.com/rebelnato/goAIRA/endpoints"
	"github.com/rebelnato/goAIRA/isolatedfunctions"
)

type Createincident struct {
	ConsumerId string `header:"ConsumerId" binding:"required"`
	ShortDesc  string `header:"ShortDesc" binding:"required"`
	Desc       string `header:"Desc" binding:"required"`
	Caller     string `header:"Caller" binding:"required"`
	Channel    string `header:"Channel" binding:"required"`
	Impact     string `header:"Impact" binding:"required"`
	Urgency    string `header:"Urgency" binding:"required"`
}

type DefaultMandate struct {
	ConsumerId  string `header:"ConsumerId" binding:"required"`
	IncidentNum string `header:"IncidentNum" binding:"required"`
}

type GetSysId struct {
	SysId string `json:"sys_id"`
}

type ExtractSysId struct {
	Result []GetSysId `json:"result"`
}

/*
  Available close codes:
  Duplicate
  Known error
  No resolution provided
  Resolved by caller
  Resolved by change
  Resolved by problem
  Resolved by request
  Solution provided
  Workaround provided
  User error
*/

type ResolveInc struct {
	CloseNotes string `header:"CloseNotes" binding:"required"`
	CloseCode  string `header:"CloseCode" binding:"required"`
}

type Updateincident struct {
	ConsumerId       string `header:"ConsumerId" binding:"required"`
	IncidentNum      string `header:"IncidentNum" binding:"required"`
	Comment          string `header:"Comment"`
	WorkNote         string `header:"WorkNote"`
	AssignmentGroup  string `header:"AssignmentGroup"`
	Description      string `header:"Description"`
	ShortDescription string `header:"ShortDescription"`
	Status           string `header:"Status"`
}

func CreateSNOWIncident(c *gin.Context) {

	var req Createincident
	var incidentDetails map[string]map[string]interface{}

	// Bind JSON and validate required fields
	validateRequest := c.ShouldBindHeader(&req)
	if validateRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": "Please provide all necessary headers ConsumerId, ShortDesc, Desc,Caller, Channel, Impact, Urgency."})
		return
	}

	if test := isolatedfunctions.ConsumerIDValidator(c, req.ConsumerId); test != true {
		return
	}

	snowURL, snowAuthToken, err := authorizations.GetSNOWAuthToken()
	if err != nil {
		log.Println("Failed to fetch SNOW auth token")
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": "Failed to fetch SNOW auth token",
		})
		return
	}

	payload := map[string]interface{}{
		"caller_id":         req.Caller,
		"short_description": req.ShortDesc,
		"desc":              req.Desc,
		"contact_type":      req.Channel,
		"impact":            req.Impact,
		"urgency":           req.Urgency,
	} // Generates payload for create incident API

	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal json")
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": "Failed to marshal request body for create incident API",
		})
		return
	}

	body, err := isolatedfunctions.POSTjsonPayload(c, snowAuthToken.AccessToken, "POST", snowURL+"api/now/table/incident", reqBody)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": err,
		})
		return
	} // Making new http request which returns []byte response and error

	// Insert request payload into db
	requestId := isolatedfunctions.UniqueIdGenerator()
	if endpoints.ConfigData.DbToggle {
		if err := isolatedfunctions.CreateRequestEntry(requestId, req.ConsumerId, "SNOW", "POST", reqBody); err != nil {
			log.Panic("Unable to insert data into requests table due to ", err)
		}
	}

	if err := json.Unmarshal(body, &incidentDetails); err != nil {
		log.Println("Response json unmarshal failed for CreateSNOWIncident")
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": "Failed to unmarshal response of create incident API",
		})
		return
	}
	// Parse JSON response

	// Insert response into db
	if endpoints.ConfigData.DbToggle {
		if err := isolatedfunctions.CreateResponseEntry(requestId, body); err != nil {
			log.Panic("Unable to insert data into responses table due to ", err)
		}
	}

	result := incidentDetails[`result`]
	incidentNum := result["number"]
	sysId := result["sys_id"]

	apiResponse := gin.H{
		"status": "success",
		"data": gin.H{
			"number":      incidentNum,
			"incidentURL": snowURL + "now/nav/ui/classic/params/target/incident.do%3Fsys_id%3D" + sysId.(string),
		},
	}

	c.JSON(http.StatusOK, apiResponse)
}

func GetSNOWIncident(c *gin.Context) {

	var req DefaultMandate
	var incidentDetails map[string]interface{}

	// Bind JSON and validate required fields
	validateRequest := c.ShouldBindHeader(&req)
	if validateRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": "Please provide all necessary headers consumerid , incidentNum."})
		return
	}

	if test := isolatedfunctions.ConsumerIDValidator(c, req.ConsumerId); test != true {
		return
	}

	snowURL, snowAuthToken, err := authorizations.GetSNOWAuthToken()
	if err != nil {
		log.Println("Failed to fetch SNOW auth token")
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": "Failed to fetch SNOW auth token",
		})
		return
	}

	body, err := isolatedfunctions.POSTjsonPayload(c, snowAuthToken.AccessToken, "GET", snowURL+"api/now/v1/table/incident?sysparm_query=number="+req.IncidentNum, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": err,
		})
		return
	} // Making new http request which returns []byte response and error

	// Insert request payload into db
	requestId := isolatedfunctions.UniqueIdGenerator()
	if endpoints.ConfigData.DbToggle {
		if err := isolatedfunctions.CreateRequestEntry(requestId, req.ConsumerId, "SNOW", "GET", nil); err != nil {
			log.Panic("Unable to insert data into requests table due to ", err)
		}
	}
	if err := json.Unmarshal(body, &incidentDetails); err != nil {
		log.Println("Response json unmarshal failed for GetSNOWIncident", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": "Failed to unmarshal response of get incident API",
		})
		return
	}
	// Parse JSON response

	// Insert response into db
	if endpoints.ConfigData.DbToggle {
		if err := isolatedfunctions.CreateResponseEntry(requestId, body); err != nil {
			log.Panic("Unable to insert data into responses table due to ", err)
		}
	}

	c.JSON(http.StatusOK, incidentDetails)
}

func UpdateSNOWIncident(c *gin.Context) {

	var resolution ResolveInc
	var req Updateincident
	var incidentDetails map[string]interface{}
	var extractSysId ExtractSysId

	// Bind JSON and validate required fields
	validateRequest := c.ShouldBindHeader(&req)
	if validateRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": "Please provide all necessary headers ConsumerId, IncidentNum."})
		return
	}

	if test := isolatedfunctions.ConsumerIDValidator(c, req.ConsumerId); test != true {
		return
	}

	if req.Status == "6" {
		validateResolution := c.ShouldBindHeader(&resolution)
		if validateResolution != nil {
			c.JSON(http.StatusBadRequest, gin.H{"response": "Please provide all necessary headers for moving incident to resolved status CloseNotes, CloseCode."})
			return
		}
	}

	snowURL, snowAuthToken, err := authorizations.GetSNOWAuthToken()
	if err != nil {
		log.Println("Failed to fetch SNOW auth token")
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": "Failed to fetch SNOW auth token",
		})
		return
	}

	sys_id, err := isolatedfunctions.POSTjsonPayload(c, snowAuthToken.AccessToken, "GET", snowURL+"/api/now/table/incident?sysparm_query=number="+req.IncidentNum+"&sysparm_fields=sys_id", nil)
	if err != nil {
		log.Println("Failed to get sys id for provided incident")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"reason": "Failed to get sys id for provided incident",
		})
		return
	} // Fetching sys id associated with the provided Incident number

	if err := json.Unmarshal(sys_id, &extractSysId); err != nil {
		log.Println("Response json unmarshal failed while extracting sysid due to ", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": "Failed to unmarshal response to get sysId",
		})
		return
	}

	if err == nil && len(extractSysId.Result) == 0 {
		log.Println("Please check the provided incident for genuinity.")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"reason": "Please check the provided incident for genuinity.",
		})
		return
	} // Check whether the sys_id was returned as response by checking array length of extractSysId.Result

	payload := make(map[string]string)

	if req.Comment != "" {
		payload["comments"] = req.Comment
	}
	if req.WorkNote != "" {
		payload["work_notes"] = req.WorkNote
	}
	if req.AssignmentGroup != "" {
		payload["assignment_group"] = req.AssignmentGroup
	}
	if req.Description != "" {
		payload["description"] = req.Description
	}
	if req.ShortDescription != "" {
		payload["short_description"] = req.ShortDescription
	}
	if req.Status != "" {
		payload["state"] = req.Status
	}
	if resolution.CloseNotes != "" {
		payload["close_notes"] = resolution.CloseNotes
	}
	if resolution.CloseCode != "" {
		payload["close_code"] = resolution.CloseCode
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("Json marshal for request payload failed")
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": err,
		})
	}

	body, err := isolatedfunctions.POSTjsonPayload(c, snowAuthToken.AccessToken, "PATCH", snowURL+"api/now/v1/table/incident/"+extractSysId.Result[0].SysId, reqBody)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": err,
		})
		return
	} // Making new http request which returns []byte response and error

	// Insert request payload into db
	requestId := isolatedfunctions.UniqueIdGenerator()
	if endpoints.ConfigData.DbToggle {
		if err := isolatedfunctions.CreateRequestEntry(requestId, req.ConsumerId, "SNOW", "PATCH", reqBody); err != nil {
			log.Panic("Unable to insert data into requests table due to ", err)
		}
	}

	if err := json.Unmarshal(body, &incidentDetails); err != nil {
		log.Println("Response json unmarshal failed for GetSNOWIncident", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"reason": "Failed to unmarshal response of get incident API",
		})
		return
	}
	// Parse JSON response

	// Insert response into db
	if endpoints.ConfigData.DbToggle {
		if err := isolatedfunctions.CreateResponseEntry(requestId, body); err != nil {
			log.Panic("Unable to insert data into responses table due to ", err)
		}
	}

	c.JSON(http.StatusOK, incidentDetails)
}
