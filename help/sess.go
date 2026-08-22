package help

import (
	"encoding/gob"
	"errors"
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
	"golang.org/x/crypto/bcrypt"
	"github.com/jmoiron/sqlx"
	"github.com/gofiber/storage/sqlite3/v2"
	"crow/help/strs"
	
)

var validate = validator.New()
func BasicStart(app *fiber.App){

	app.Use("/*", static.New("frontend/glitter/"))

	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} ${path}\n",
	}))

	gob.Register(&strs.User{}) 

	storage := sqlite3.New(sqlite3.Config{
        Database: "./data.db",   // can be the same file as your app DB, or separate
        Table:    "sessions",
        Reset:    false,
    })

	sessionMiddleware, sessionStore := session.NewWithStore(session.Config{
		Storage:          storage,
	    CookieSecure:      false,
	    CookieHTTPOnly:    true,
	    CookieSameSite:    "Lax",
	    IdleTimeout:       6*30*24 * time.Hour,
	    AbsoluteTimeout:   5*365*24 * time.Hour,
	})

	app.Use(sessionMiddleware) 

	app.Use(csrf.New(csrf.Config{
	    CookieHTTPOnly:    false,
	    CookieSameSite:    "Lax",
	    CookieSessionOnly: true,
	    Extractor:         extractors.FromHeader("X-Csrf-Token"),
	    Session:           sessionStore,
	}))
}

func Authmid(c fiber.Ctx) error {

	sess := session.FromContext(c)
	user:= sess.Get("user")
	if user == nil {
		return Redirect(c, "/")
	}
	c.Locals("user", user)
	return c.Next()
}

func Authmidr(c fiber.Ctx) error {

	sess := session.FromContext(c)
	if user:= sess.Get("user"); user != nil {
		return Redirect(c, "/")
	}
	return c.Next()
}

func Logcheck(c fiber.Ctx) bool {
	sess:= session.FromContext(c)
	if user:= sess.Get("user"); user==nil {
		return false
	}
	return true
}

func Signup(c fiber.Ctx, db *sqlx.DB) error {
    user := new(strs.User)

    if err := c.Bind().Form(user); err != nil {
        return err
    }
	if err := validate.Struct(user); err != nil {
        return err
    }

	if bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12); err!= nil{
		return err
	} else {
		user.Password = string(bytes)
	}

	if err := db.QueryRowx("INSERT INTO users(userName, password, access) VALUES ($1, $2, $3) RETURNING *", user.UserName, user.Password, user.Access).StructScan(user); err != nil {
		return err
	}

	sess:= session.FromContext(c)
	sess.Set("user", user)

    return nil
}

func Login(c fiber.Ctx, db *sqlx.DB) error {
	user:= new(strs.User)

	if err:= c.Bind().Form(user); err!=nil {
		return err
	}
	if err := validate.Struct(user); err != nil {
        return err
    }

	dbuser := new(strs.User)
	err:= db.Get(dbuser, "select * from users where userName= $1", user.UserName)

	if err!= nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(user.Password)); err == nil{
		sess:= session.FromContext(c)
		sess.Set("user", dbuser)
		return nil
	}
	return errors.New("wrong credentials")
	
}

func Logout(c fiber.Ctx) error {
	if err := csrf.HandlerFromContext(c).DeleteToken(c); err != nil {
        return err
    }
	if err :=session.FromContext(c).Destroy(); err!=nil{
		return err
	}
	return nil
}