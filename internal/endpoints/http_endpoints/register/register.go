package register

import (
	"encoding/json"
	"genericAPI/internal/customErrors"
	"genericAPI/internal/models/user"
	"genericAPI/internal/utils/authentication_utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/mail"
)

type registerRequestBody struct {
	Mail     *string `json:"mail"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func RegisterEndpoint(ctx *gin.Context) {
	body, _ := io.ReadAll(ctx.Request.Body)
	var reqBody registerRequestBody
	if err := json.Unmarshal(body, &reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON body!"})
		return
	}

	if err := validateBody(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userModel := user.User{
		Username: *reqBody.Username,
		Mail:     *reqBody.Mail,
	}

	if userModel.Exists() {
		ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists."})
		return
	}

	userModel.Password = authentication_utils.HashPassword(*reqBody.Password)

	if err := userModel.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// todo mail verification
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully created user."})
}

func validateBody(body *registerRequestBody) error {
	if body.Mail == nil || body.Username == nil || body.Password == nil {
		return customErrors.ErrMissingField
	}

	_, err := mail.ParseAddress(*body.Mail)
	if err != nil {
		return err
	}

	return nil
}
