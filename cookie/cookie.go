package cookie

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minilikmila/standard-auth-go/config"
)

func SetCookie(c *gin.Context, data string, config *config.Config) {

	// c.SetCookie("genie_session", data, 3600, "/", config.Host, false, true)
	// maxAge is more more essential and powerful to define the duration of cookie stay valid in client's browser.
	// make it maxAte = -1 to expire the cookie
	cookie := &http.Cookie{
		Name:     "go_auth_session",
		Value:    data,
		Path:     "/",
		Domain:   config.Host, //The Domain attribute in a cookie refers to the domain of the server that sets the cookie,
		Secure:   false,
		HttpOnly: true,
		MaxAge:   int(config.JWT.Exp * 60),
		// Expires:  time.Now().Add(time.Minute * config.JWT.Exp),
	}
	http.SetCookie(c.Writer, cookie)

}

func GetCookie(c *gin.Context) (string, error) {
	cookie, err := c.Request.Cookie("go_auth_session") // direct http request / mux
	// cookie, err := c.Cookie("genie_session") // *gin version
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
