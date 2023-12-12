package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func validateAgainstSSO(c *gin.Context) {
	token, err := c.Request.Cookie("auth_token")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Get("https://localhost/api/users/auth/" + token.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Next()
}

func tokenAuth(c *gin.Context) {
	token := c.Query("token")
	fmt.Println(token)
	fmt.Println(c.Request.URL)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Get("https://localhost/api/users/auth/" + token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.SetCookie("auth_token", token, 86400, "/", "dev.gfed", false, false)

	c.Redirect(http.StatusFound, "https://sample.dev.gfed:8080/home")
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")

	cookie, err := c.Request.Cookie("token")

	if cookie != nil && err == nil {
		c.SetCookie("token", "", -1, "/", "*", false, true)
	}

	if id == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No session."})
		return
	}
	session.Delete("id")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out!"})
}
