// Bulk of code from Auth0's setup instructions

package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"slugquest.com/backend/crud"
)

const FRONTEND_HOST string = "localhost:5185"

// Checks if user is authenticated before redirecting to next page
func IsAuthenticated(c *gin.Context) {
	// Auth token: for direct calls to this endpoint
	auth_token := c.GetHeader("Authorization")

	// Should have user profile saved to session
	user_profile := sessions.Default(c).Get("profile")

	if auth_token == "" && user_profile == nil {
		// c.Redirect(http.StatusSeeOther, "/") // TODO: maybe make an "Oops, wrong page"
		c.String(http.StatusUnauthorized, "Forbidden")
		c.Abort()
	} else {
		c.Next()
	}
}

// Handler for our login.
func LoginHandler(auth *Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session.
		session := sessions.Default(c)
		session.Clear()

		session.Set("state", state)
		if err := session.Save(); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func LogoutHandler(c *gin.Context) {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// Return to the not logged in page
	// returnTo, err := url.Parse(scheme + "://" + c.Request.Host)
	returnTo, err := url.Parse(scheme + "://" + FRONTEND_HOST)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}

func CallbackHandler(auth *Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if c.Query("state") != session.Get("state") {
			c.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(c.Request.Context(), c.Query("code"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(c.Request.Context(), token)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		var userInfoStruct *crud.User = getUserInfo(c)
		if userInfoStruct == nil {
			c.String(http.StatusInternalServerError, "Couldn't retrieve user profile.")
			return
		}

		session.Set("user_profile", userInfoStruct)
		if err := session.Save(); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		c.Redirect(http.StatusTemporaryRedirect, "http://"+FRONTEND_HOST+"/loggedin")
	}
}

func getUserInfo(c *gin.Context) *crud.User {
	session := sessions.Default(c)

	profile, ok := session.Get("profile").(map[string]interface{})
	if !ok || profile == nil {
		c.String(http.StatusInternalServerError, "Couldn't retrieve user profile.")
		return nil
	}

	// No user id? No SlugQuest.
	sesUID, ok := profile["sub"].(string)
	if !ok {
		c.String(http.StatusInternalServerError, "Couldn't resolve user id.")
		return nil
	}

	sesUsername, ok := profile["name"].(string)
	if !ok {
		c.String(http.StatusInternalServerError, "Couldn't resolve username.")
		return nil
	}

	sesPFP, ok := profile["picture"].(string)
	if !ok {
		c.String(http.StatusInternalServerError, "Couldn't resolve profile picture URL.")
		return nil
	}

	// Check if user exists in our DB
	user, found, err := crud.GetUser(sesUID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't fetch user.")
		return nil
	}

	if found {
		// Check if any edits need to be made
		if user.Username != sesUsername || user.Picture != sesPFP {
			user.Username = sesUsername
			user.Picture = sesPFP

			edited, err := crud.EditUser(user, sesUID)
			if err != nil || !edited {
				c.String(http.StatusInternalServerError, "Couldn't update user.")
				return nil
			}
		}
	} else {
		// Need to populate and add a new user
		user = crud.User{
			UserID:   sesUID,
			Username: sesUsername,
			Picture:  sesPFP,
			Points:   0,
			BossId:   1,
		}

		added, err := crud.AddUser(user)
		if err != nil || !added {
			c.String(http.StatusInternalServerError, "Couldn't register user into our records.")
			return nil
		}
	}

	return &user
}

// Sends user profile from the current session as JSON
// func UserProfileHandler(c *gin.Context) {
// 	session := sessions.Default(c)
// 	profile := session.Get("profile")

// 	if profile == nil {
// 		c.String(http.StatusInternalServerError, "Failed to retrieve user information.")
// 		return
// 	}

// 	c.HTML(http.StatusOK, "/template/user.html", profile)
// }
