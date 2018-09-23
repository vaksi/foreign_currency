/*  APIResponse.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 05, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 05/09/18 04:07
 */

package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/vaksi/foreign_currency/constants"
)

// APIResponse defines attributes for api Response
type APIResponse struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Write writes the data to http response writer
func Write(w http.ResponseWriter, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

// Defaults API Response with standard HTTP Status Code. The default value can
// be changed either using `ModifyMessage` or `ModifyHTTPCode`. You can call it
// directly by :
//
// response.Write(res, response.APIErrorUnknown)
// return
var (
	// A generic error message, given when an unexpected condition was encountered.
	APIErrorUnknown = APIResponse{
		Code:    constants.ErrorUnknownCode,
		Message: "Internal Server Error, Please Try Again",
	}

	// Standard response for successful HTTP requests.
	APIOK = APIResponse{
		Code:    constants.GeneralSuccessCode,
		Message: "Success",
	}

	// The request has been fulfilled, resulting in the creation of a new resource
	APICreated = APIResponse{
		Code:    constants.GeneralSuccessCode,
		Message: "Success",
	}

	// The request has been accepted for processing, but the processing has not been completed.
	APIAccepted = APIResponse{
		Code:    constants.GeneralSuccessCode,
		Message: "Success",
	}

	// API Data not found
	APIDataNotFound = APIResponse{
		Code:    constants.DataNotFound,
		Message: "Success",
	}

	// API Data Duplicate
	APIOKDuplicatedExpression = APIResponse{
		Code:    constants.DuplicatedExpressionCode,
		Message: "Dupplicated Expression",
	}

	// The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax
	// , size too large, invalid request message framing, or deceptive request routing).
	APIErrorValidation = APIResponse{
		Code:    constants.InvalidKeyFieldValue,
		Message: "Invalid Input Data",
	}

	APIErrorInvalidData = APIResponse{
		Code:    constants.ErrorInvalidData,
		Message: "Invalid Data Format",
	}

	// The request was valid, but the server is refusing action. The user might not have the necessary permissions for
	// a resource, or may need an account of some sort.
	APIErrorForbidden = APIResponse{
		Code:    constants.ErrorForbiddenCode,
		Message: "Action forbidden",
	}

	// The request was valid, but the server is refusing action. The user might not have the necessary permissions for
	// a resource, or may need an account of some sort.
	APIErrorUnauthorized = APIResponse{
		Code:    constants.ErrorUnauthorizedCode,
		Message: "Unauthorized",
	}
)
