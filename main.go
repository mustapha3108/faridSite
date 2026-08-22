package main

import (
	"crow/frontend/mark/comp"
	"crow/frontend/mark/pages"
	"crow/help"
	"github.com/gofiber/fiber/v3"
	//"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/joho/godotenv"
)


func main() {


	//todo:="make an app where multiple users can create accounts and upload posts with pictures, dashboard + global + posts filtering, if have time add animation"

	godotenv.Load()

	//body limit man, body limit
	app:= fiber.New(fiber.Config{
	    // Default is 	10MB. For a SaaS, you might want 10MB or 20MB.
	    BodyLimit: 10 * 1024 * 1024, 
	})

	db:=help.BasicDb()

	help.BasicStart(app)


	//home
	app.Get("/", func(c fiber.Ctx) error {
		return help.Hrender(c, pages.Home(help.Logcheck(c)))
	})

	//atelier
	app.Get("/atelier", func(c fiber.Ctx) error {
		return help.Hrender(c, pages.Atelier())
	})

	//projets
	app.Get("/projets", func(c fiber.Ctx) error {
		return help.Hrender(c, pages.Projets())
	})

	//book
	app.Get("/book", func(c fiber.Ctx) error {
		return help.Hrender(c, pages.Book())
	})

	//candidature
	app.Get("/candidature", func(c fiber.Ctx) error{
		return help.Hrender(c, pages.Candidature())
	})

	//contact
	app.Get("/contact", func(c fiber.Ctx) error {
		return help.Hrender(c, pages.Contact())
	})

	//rating
	app.Get("/star", func(c fiber.Ctx) error {
		return help.Hrender(c, pages.Star())
	})

	//chief
	app.Get("/chief", func(c fiber.Ctx) error {
		return help.Hrender(c, pages.Chief())
	})

	app.Post("/login", func(c fiber.Ctx) error {
		if err:=help.Login(c, db); err!=nil {
			return c.SendString((help.ShowError(err)))
		}
		return help.Redirect(c, "/dash")
	})

	//dashboard
	app.Get("/dash", help.Authmid, func(c fiber.Ctx) error {
		user:= c.Locals("user").(*help.User)
		return help.Hrender(c, pages.Dash(user))
	})

	app.Post("/dashNav", help.Authmid, func (c fiber.Ctx) error {
		target := c.FormValue("target")
		switch target {
		case "Contact":
			return help.Render(c, comp.ContactMod())
		case "Candidats":
			return help.Render(c, comp.CandidatsMod())
		case "Atelier":
			return help.Render(c, comp.AtelierMod())
		case "Comptes":
			return help.Render(c, comp.ComptesMod())
		case "Project":
			return help.Render(c, comp.ProjetsMod())
		case "Stars":
			return help.Render(c, comp.StarsMod())
		case "Messages":
			return help.Render(c, comp.MessagesMod())
		}
		return c.SendString("page introuvable")
	})


	app.Post("/alpinetest", func(c fiber.Ctx) error {
		return help.Render(c, comp.Alpinetest())
	})

	app.Post("/createUser", func(c fiber.Ctx) error {

		err:= help.Signup(c, db)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	
/*
	app.Get("/posts", func(c fiber.Ctx) error {
		var posts []help.Post
		if err:= db.Select(&posts, "select * from posts"); err!= nil {
			return c.SendString(help.ShowError(err))
		}
		return help.Hrender(c, pages.Posts(posts))
	})

	app.Get("/account", help.Authmid, func(c fiber.Ctx) error {
		user := c.Locals("user").(*help.User)
		return help.Hrender(c, pages.Account(user))
	})

	app.Get("/sign", help.Authmidr, func(c fiber.Ctx) error {
		return help.Hrender(c, pages.Sign())
	})
	

	//Posts
	app.Post("/signup", func(c fiber.Ctx) error {

		err:= help.Signup(c, db)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		return help.Redirect(c, "/account")
	})

	app.Post("/login", help.Authmidr, func(c fiber.Ctx) error {
		if  err:=help.Login(c, db); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		return help.Redirect(c, "/account")

	})

	app.Post("/logout", help.Authmid, func(c fiber.Ctx) error {
		if err:=help.Logout(c); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		return help.Redirect(c, "/")
	})

	app.Post("/uploadPoste", func(c fiber.Ctx) error {
		poste:= new(help.Post)
		if err:= c.Bind().Form(poste); err!= nil {
			return c.SendString(help.ShowError(err))
		}
		if imagepath, err :=  help.SaveImage(c, poste.Image); err!= nil{
			return c.SendString(help.ShowError(err))
		} else {
			poste.ImagePath = imagepath
		}

		user:= session.FromContext(c).Get("user").(*help.User)

		if _, err:= db.Exec("insert into posts(userid, postname, description, imagepath) VALUES($1, $2, $3, $4)",
		user.UserID, poste.PostName, poste.Description, poste.ImagePath);
		err!= nil {
			return c.SendString(help.ShowError(err))
		 }

		return help.Redirect(c, "/posts")
	})
*/
	app.Listen(":3000")
}