package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minilikmila/standard-auth-go/internal/domain/model"
	"github.com/minilikmila/standard-auth-go/internal/domain/model/enum"
	"github.com/minilikmila/standard-auth-go/internal/domain/validators"
	"github.com/minilikmila/standard-auth-go/pkg/utils"

	"github.com/minilikmila/standard-auth-go/internal/service"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// SignUp handles user registration
func SignUp(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body model.SignUpForm

		// Bind and validate request body
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid request body",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		// Validate required fields
		if body.Email == "" || body.Password == "" || body.Name == "" {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Email, password, and name are required",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		// Check if user already exists
		exists, err := authService.UserExists(ctx, body.Email)
		if err != nil {
			log.Errorf("Error checking user existence: %v", err)
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "Error checking user existence",
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			})
			return
		}
		if exists {
			ctx.JSON(http.StatusConflict, model.Response{
				Message:    "User already exists",
				StatusCode: http.StatusConflict,
			})
			return
		}

		// Create user
		user := &model.User{
			Name:         &body.Name,
			Email:        &body.Email,
			Phone:        &body.Phone,
			Password:     &body.Password,
			Role:         "user",             // Default role
			SignUpMethod: enum.RegularSignUp, // Default sign-up method
		}

		if err := authService.CreateUser(ctx, user); err != nil {
			log.Errorf("Error creating user: %v", err)
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "Error creating user",
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			})
			return
		}

		// Send verification email if email is provided
		if body.Email != "" {
			logrus.Infoln("Sending email verification")
			if err := authService.SendVerificationEmail(ctx, user); err != nil {
				log.Errorf("Error sending verification email: %v", err)
				ctx.JSON(http.StatusInternalServerError, model.Response{
					Message:    "Error sending verification email",
					StatusCode: http.StatusInternalServerError,
					Error:      err,
				})
				return
			}
		}

		ctx.JSON(http.StatusCreated, model.Response{
			Message:    "User created successfully. Please check your email to verify your account.",
			StatusCode: http.StatusCreated,
			Data: gin.H{
				"user": gin.H{
					"id":    user.ID,
					"name":  user.Name,
					"email": user.Email,
					"role":  user.Role,
				},
			},
		})
	}
}

// Login handles user authentication
func Login(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body model.LoginForm

		// Bind and validate request body
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid request body",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		// Validate required fields
		if body.Email == "" || body.Password == "" {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Email and password are required",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		// Authenticate user
		user, err := authService.Login(ctx, body.Email, body.Password)
		if err != nil {
			switch err {
			case service.ErrUserNotFound, service.ErrEmailNotVerified, service.ErrPhoneNotVerified:
				ctx.JSON(http.StatusNotFound, model.Response{
					Message:    "User not found",
					StatusCode: http.StatusNotFound,
				})
			case service.ErrInvalidCredentials:
				ctx.JSON(http.StatusUnauthorized, model.Response{
					Message:    "Invalid credentials",
					StatusCode: http.StatusUnauthorized,
				})
			default:
				log.Errorf("Error during login: %v", err)
				ctx.JSON(http.StatusInternalServerError, model.Response{
					Message:    "Error during login",
					StatusCode: http.StatusInternalServerError,
					Error:      err,
				})
			}
			return
		}

		// Generate tokens
		accessToken, refreshToken, err := authService.GenerateTokens(ctx, user)
		if err != nil {
			log.Errorf("Error generating tokens: %v", err)
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "Error generating tokens",
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			})
			return
		}

		// Set cookies
		ctx.SetCookie("access_token", accessToken, int(time.Hour*24), "/", "", false, true)
		ctx.SetCookie("refresh_token", refreshToken, int(time.Hour*24*7), "/", "", false, true)

		ctx.JSON(http.StatusOK, model.Response{
			Message:    "Login successful",
			StatusCode: http.StatusOK,
			Data: gin.H{
				"user": gin.H{
					"id":    user.ID,
					"name":  user.Name,
					"email": user.Email,
					"role":  user.Role,
				},
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		})
	}
}

// Logout handles user logout
func Logout(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Try to get token from Authorization header
		token := ctx.GetHeader("Authorization")
		accessToken := ""
		if token != "" {
			accessToken = utils.ExtractTokenFromHeader(token)
		}
		// If not found, try to get from cookie
		if accessToken == "" {
			cookie, err := ctx.Cookie("access_token")
			if err == nil {
				accessToken = cookie
			}
		}
		if accessToken == "" {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Authorization token is required",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		// Proceed with logout logic (e.g., blacklist token)
		err := authService.Logout(ctx, accessToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "Logout failed",
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			})
			return
		}

		// Optionally clear cookies
		ctx.SetCookie("access_token", "", -1, "/", "", false, true)
		ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

		ctx.JSON(http.StatusOK, model.Response{
			Message:    "Logout successful",
			StatusCode: http.StatusOK,
		})
	}
}

// RefreshToken handles token refresh
func RefreshToken(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshToken := ctx.GetHeader("Authorization")
		if refreshToken == "" {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Refresh token is required",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		accessToken, newRefreshToken, err := authService.RefreshToken(ctx, refreshToken)
		if err != nil {
			log.Errorf("Error refreshing token: %v", err)
			ctx.JSON(http.StatusUnauthorized, model.Response{
				Message:    "Invalid refresh token",
				StatusCode: http.StatusUnauthorized,
				Error:      err,
			})
			return
		}

		// Set cookies
		ctx.SetCookie("access_token", accessToken, int(time.Hour*24), "/", "", false, true)
		ctx.SetCookie("refresh_token", newRefreshToken, int(time.Hour*24*7), "/", "", false, true)

		ctx.JSON(http.StatusOK, model.Response{
			Message:    "Token refreshed successfully",
			StatusCode: http.StatusOK,
			Data: gin.H{
				"access_token":  accessToken,
				"refresh_token": newRefreshToken,
			},
		})
	}
}

// VerifyEmail handles email verification
func VerifyEmail(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body struct {
			Token string `json:"token" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid request body",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		if err := authService.VerifyEmail(ctx, body.Token); err != nil {
			log.Errorf("Error verifying email: %v", err)
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid verification token",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message:    "Email verified successfully",
			StatusCode: http.StatusOK,
		})
	}
}

// VerifyPhone handles phone verification
func VerifyPhone(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body struct {
			Code  string `json:"code" binding:"required"`
			Phone string `json:"phone,omitempty"`
		}

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid request body",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		if err := authService.VerifyPhone(ctx, body.Code, body.Phone); err != nil {
			log.Errorf("Error verifying phone: %v", err)
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid verification code",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message:    "Phone verified successfully",
			StatusCode: http.StatusOK,
		})
	}
}

// ForgotPassword handles password reset request
func ForgotPassword(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logrus.Info("Processing forgot password request")
		var body struct {
			Email string `json:"email" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&body); err != nil {
			logrus.Errorf("Error binding JSON: %v", err)

			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid request body",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		valid := validators.IsValidEmail(body.Email)
		if !valid {
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "Error processing forgot password request",
				StatusCode: http.StatusInternalServerError,
				Error:      errors.New("Invalid email format"),
			})
			return
		}

		if err := authService.ForgotPassword(ctx, body.Email); err != nil {
			log.Errorf("Error processing forgot password request: %v", err)
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "Error processing forgot password request",
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message:    "Password reset instructions sent to your email",
			StatusCode: http.StatusOK,
		})
	}
}

// ResetPassword handles password reset
func ResetPassword(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body struct {
			Token       string `json:"token" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid request body",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		if err := authService.ResetPassword(ctx, body.Token, body.NewPassword); err != nil {
			log.Errorf("Error resetting password: %v", err)
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid or expired reset token",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message:    "Password reset successfully",
			StatusCode: http.StatusOK,
		})
	}
}

// GetProfile handles profile retrieval
func GetProfile(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("id")
		if userID == "" {
			ctx.JSON(http.StatusUnauthorized, model.Response{
				Message:    "Unauthorized",
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		user, err := authService.GetProfile(ctx, userID)
		if err != nil {
			log.Errorf("Error getting profile: %v", err)
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "Error getting profile",
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message:    "Profile retrieved successfully",
			StatusCode: http.StatusOK,
			Data: gin.H{
				"user": gin.H{
					"id":                user.ID,
					"name":              user.Name,
					"email":             user.Email,
					"phone":             user.Phone,
					"role":              user.Role,
					"profile_picture":   user.ProfilePicture,
					"sign_up_method":    user.SignUpMethod,
					"is_email_verified": user.IsEmailVerified,
					"is_phone_verified": user.IsPhoneVerified,
					"last_login_at":     user.LastLoginAt,
					"email_verified_at": user.EmailVerifiedAt,
					"phone_verified_at": user.PhoneVerifiedAt,
				},
			},
		})
	}
}

// UpdateProfile handles profile updates
func UpdateProfile(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString("user_id")
		if userID == "" {
			ctx.JSON(http.StatusUnauthorized, model.Response{
				Message:    "Unauthorized",
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		var updates map[string]interface{}
		if err := ctx.ShouldBindJSON(&updates); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{
				Message:    "Invalid request body",
				StatusCode: http.StatusBadRequest,
				Error:      err,
			})
			return
		}

		if err := authService.UpdateProfile(ctx, userID, updates); err != nil {
			log.Errorf("Error updating profile: %v", err)
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message:    "Error updating profile",
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message:    "Profile updated successfully",
			StatusCode: http.StatusOK,
		})
	}
}

// // ChangePassword handles password change
// func ChangePassword(authService service.AuthService) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		userID := ctx.GetString("user_id")
// 		if userID == "" {
// 			ctx.JSON(http.StatusUnauthorized, model.Response{
// 				Message:    "Unauthorized",
// 				StatusCode: http.StatusUnauthorized,
// 			})
// 			return
// 		}

// 		var body struct {
// 			CurrentPassword string `json:"current_password" binding:"required"`
// 			NewPassword     string `json:"new_password" binding:"required,min=8"`
// 		}

// 		if err := ctx.ShouldBindJSON(&body); err != nil {
// 			ctx.JSON(http.StatusBadRequest, model.Response{
// 				Message:    "Invalid request body",
// 				StatusCode: http.StatusBadRequest,
// 				Error:      err,
// 			})
// 			return
// 		}

// 		if err := authService.UpdatePassword(ctx, userID, body.CurrentPassword, body.NewPassword); err != nil {
// 			log.Errorf("Error changing password: %v", err)
// 			ctx.JSON(http.StatusBadRequest, model.Response{
// 				Message:    "Error changing password",
// 				StatusCode: http.StatusBadRequest,
// 				Error:      err,
// 			})
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, model.Response{
// 			Message:    "Password changed successfully",
// 			StatusCode: http.StatusOK,
// 		})
// 	}
// }
