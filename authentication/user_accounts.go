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
)

const FRONTEND_HOST string = "localhost:5185"

// TODO: make this more elegant with Gin sessions or something
var Curr_user_id string = "hi"

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

		// Extract Auth0's provided user vid
		if profile["sub"] == nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		user_id := profile["sub"].(string)[len("auth0|"):]
		Curr_user_id = user_id

		// Redirect to logged in page.
		c.Redirect(http.StatusTemporaryRedirect, "http://"+FRONTEND_HOST+"/loggedin")
	}
}

// Displays user profile from the current session
// func UserProfileHandler(c *gin.Context) {
// 	session := sessions.Default(c)
// 	profile := session.Get("profile")

// 	if profile == nil {
// 		c.String(http.StatusInternalServerError, "Failed to retrieve user information.")
// 		return
// 	}

// 	c.HTML(http.StatusOK, "/template/user.html", profile)
// }
