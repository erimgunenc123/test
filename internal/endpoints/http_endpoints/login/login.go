package login

import (
	"encoding/json"
	"genericAPI/internal/customErrors"
	"genericAPI/internal/models/refresh_token"
	"genericAPI/internal/models/user"
	"genericAPI/internal/utils/authentication_utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type loginRequestBody struct {
	Mail                 *string `json:"mail"`
	Username             *string `json:"username"`
	Password             *string `json:"password"`
	GenerateRefreshToken bool    `json:"generate_refresh_token"`
}

func LoginEndpoint(ctx *gin.Context) {
	body, _ := io.ReadAll(ctx.Request.Body)
	var reqBody loginRequestBody
	if err := json.Unmarshal(body, &reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON body!"})
		return
	}

	if err := validateBody(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userModel *user.User
	if reqBody.Mail != nil {
		userModel = &user.User{Mail: *reqBody.Mail}
	} else {
		userModel = &user.User{Username: *reqBody.Username}
	}

	if userModel = userModel.Find(); userModel == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "User not found"})
		return
	}

	if authentication_utils.ComparePassword(userModel.Password, *reqBody.Password) {
		respBody := gin.H{}
		token, err := authentication_utils.CreateAccessToken(userModel.ID, userModel.PublicID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		respBody["access_token"] = token

		if reqBody.GenerateRefreshToken {
			refreshToken := authentication_utils.CreateRefreshToken(userModel.ID)
			refreshTokenModel := refresh_token.RefreshToken{Token: refreshToken, UserID: userModel.ID}
			err = refreshTokenModel.Save()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			respBody["refresh_token"] = refreshToken
		}
		ctx.JSON(http.StatusOK, respBody)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Wrong password."})
	}
}

func validateBody(body *loginRequestBody) error {
	if (body.Mail == nil && body.Username == nil) || body.Password == nil {
		return customErrors.ErrMissingField
	}
	return nil
}
