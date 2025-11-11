package http

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/go-resty/resty/v2"

	"github.com/theinfosecguy/archer/internal/constants"
	"github.com/theinfosecguy/archer/internal/logger"
	"github.com/theinfosecguy/archer/internal/models"
	"github.com/theinfosecguy/archer/internal/variables"
)

// Client wraps Resty client for API validation requests
type Client struct {
	restyClient *resty.Client
}

// NewClient creates a new HTTP client
func NewClient() *Client {
	return &Client{
		restyClient: resty.New(),
	}
}

// ExecuteRequest executes an HTTP request based on the template and variables
func (c *Client) ExecuteRequest(
	template *models.SecretTemplate,
	vars map[string]string,
) (*models.ValidationResult, error) {
	// Process URL
	requestURL, maskedURL := variables.ProcessURL(template.APIURL, vars)

	// Process headers
	requestHeaders, maskedHeaders := variables.ProcessHeaders(template.Request.Headers, vars)

	// Process query parameters
	requestQueryParams, maskedQueryParams := variables.ProcessQueryParams(template.Request.QueryParams, vars)

	// Process data
	requestData, maskedData := variables.ProcessData(template.Request.Data, vars)

	// Process JSON data
	requestJSONData, maskedJSONData := variables.ProcessJSONData(template.Request.JSONData, vars)

	// Log request preparation with masked values
	logger.Debug("Preparing %s request to %s", template.Method, maskedURL)
	logger.Debug("Request headers (masked): %v", maskedHeaders)
	if maskedQueryParams != nil && len(maskedQueryParams) > 0 {
		logger.Debug("Query parameters (masked): %v", maskedQueryParams)
	}
	if maskedData != nil {
		logger.Debug("Request data (masked): %s", *maskedData)
	}
	if maskedJSONData != nil {
		logger.Debug("Request JSON data (masked): %v", maskedJSONData)
	}

	// Configure client with timeout and retries
	restyClient := c.restyClient.
		SetTimeout(time.Duration(template.Request.Timeout) * time.Second)

	// Configure retries if specified
	if template.ErrorHandling.MaxRetries > 0 {
		restyClient.
			SetRetryCount(template.ErrorHandling.MaxRetries).
			SetRetryWaitTime(time.Duration(template.ErrorHandling.RetryDelay) * time.Second).
			// Retry on 5xx server errors and 429 rate limiting
			AddRetryCondition(func(r *resty.Response, err error) bool {
				// Retry on network errors
				if err != nil {
					return true
				}
				// Retry on server errors (5xx) and rate limiting (429)
				return r.StatusCode() >= 500 || r.StatusCode() == 429
			})
	}

	// Create request
	req := restyClient.R()

	// Set headers
	if len(requestHeaders) > 0 {
		req.SetHeaders(requestHeaders)
	}

	// Set query parameters
	if len(requestQueryParams) > 0 {
		req.SetQueryParams(requestQueryParams)
	}

	// Set body
	if requestData != nil {
		req.SetBody(*requestData)
	} else if requestJSONData != nil {
		req.SetBody(requestJSONData)
	}

	// Execute request
	logger.Debug("Sending request (timeout: %ds)...", template.Request.Timeout)
	resp, err := req.Execute(template.Method, requestURL)
	if err != nil {
		// Check if it's a timeout error
		errMsg := err.Error()
		if strings.Contains(errMsg, "context deadline exceeded") ||
			strings.Contains(errMsg, "Client.Timeout exceeded") ||
			strings.Contains(errMsg, "timeout") {
			logger.Info("Request timeout after %ds", template.Request.Timeout)
			return &models.ValidationResult{
				Valid: false,
				Error: constants.RequestTimeout,
			}, nil
		}
		logger.Info("Request failed: %s", err.Error())
		return &models.ValidationResult{
			Valid: false,
			Error: fmt.Sprintf(constants.RequestFailed, err.Error()),
		}, nil
	}

	logger.Info("Request completed with status code: %d", resp.StatusCode())

	// Log response content in debug mode
	if logger.IsDebug() {
		bodyStr := string(resp.Body())
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "... (truncated)"
		}
		logger.Debug("Response content: %s", bodyStr)
	}

	// Check response against success criteria
	return c.checkResponse(resp, template)
}

// checkResponse validates the response against template success criteria
func (c *Client) checkResponse(
	resp *resty.Response,
	template *models.SecretTemplate,
) (*models.ValidationResult, error) {
	logger.Debug("Validating response against template success criteria")
	statusCode := resp.StatusCode()

	// Check status code
	statusCodeValid := false
	for _, code := range template.SuccessCriteria.StatusCode {
		if statusCode == code {
			statusCodeValid = true
			break
		}
	}

	if !statusCodeValid {
		logger.Info("Status code validation failed: got %d, expected one of %v", statusCode, template.SuccessCriteria.StatusCode)
		errorMsg := fmt.Sprintf("HTTP %d", statusCode)
		if customMsg, ok := template.ErrorHandling.ErrorMessages[statusCode]; ok {
			errorMsg = customMsg
		}
		return &models.ValidationResult{
			Valid: false,
			Error: errorMsg,
		}, nil
	}

	logger.Debug("Status code validation passed: %d is in expected range", statusCode)

	// Check required fields if specified
	if len(template.SuccessCriteria.RequiredFields) > 0 {
		logger.Debug("Checking %d required fields in JSON response", len(template.SuccessCriteria.RequiredFields))
		var responseData interface{}
		if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
			logger.Info("Response validation failed: API returned invalid JSON")
			return &models.ValidationResult{
				Valid: false,
				Error: constants.InvalidJSONResponse,
			}, nil
		}

		// Check each required field
		for _, fieldPath := range template.SuccessCriteria.RequiredFields {
			value, err := jsonpath.Get(fieldPath, responseData)
			if err != nil || value == nil {
				logger.Info("Required field validation failed: '%s' not found in response", fieldPath)
				return &models.ValidationResult{
					Valid: false,
					Error: fmt.Sprintf(constants.RequiredFieldNotFound, fieldPath),
				}, nil
			}
			logger.Debug("Required field validation passed: '%s' found in response", fieldPath)
		}
	}

	logger.Info("Validation successful")
	return &models.ValidationResult{
		Valid:   true,
		Message: constants.SecretValid,
	}, nil
}
