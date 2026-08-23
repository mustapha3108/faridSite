package main

import (
	"crow/frontend/mark/comp"
	"crow/frontend/mark/pages"
	"crow/help"
	"crow/help/strs"
	"strings"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)


func main() {


	//todo:="make an app where multiple users can create accounts and upload posts with pictures, dashboard + global + posts filtering, if have time add animation"

	godotenv.Load()

	//body limit man, body limit
	app:= fiber.New(fiber.Config{
	    // 100 mb body limit
	    BodyLimit: 100 * 1024 * 1024, 
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
		user:= c.Locals("user").(*strs.User)
		return help.Hrender(c, pages.Dash(user))
	})

	app.Post("/dashNav", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		target := c.FormValue("target")
		switch target {
		case "Contact":
			return help.Render(c, comp.ContactMod(user))
		case "Candidats":
			return help.Render(c, comp.CandidatsMod(user))
		case "Atelier":
			return help.Render(c, comp.AtelierMod(user))
		case "Comptes":
			return help.Render(c, comp.ComptesMod(user))
		case "Project":
			return help.Render(c, comp.ProjetsMod(user))
		case "Stars":
			return help.Render(c, comp.StarsMod(user))
		case "Messages":
			return help.Render(c, comp.MessagesMod(user))
		case "Logout":
			if err:=help.Logout(c); err!=nil{
				return c.SendString(help.ShowError(err))
			} else {
				return help.Redirect(c, "/")
			}
		}	
		return c.SendString("page introuvable")
	})

	app.Post("/uploadProject", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "c"){
			return c.Status(403).SendString("vous n'avez pas le droit de creer des projets")
		}
		project := new(strs.Project)
		if err:=c.Bind().Form(project); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(project); err!=nil {
			return c.SendString(help.ShowError(err))
		} 
		//i am fucking stupid, just create a random name for each one, then save all the pics in dynamix folder, wtf man don't over complicate your life
		//TODO: WRITE A FORM FOR PROJECTS, ADD CRUD FOR REST TOO


		return c.SendString("Projet créé avec succès")

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

	import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func createProjectHandler(c *fiber.Ctx) error {
	var project Project

	// Fiber's BodyParser handles multipart forms and maps fields via `form:"..."` tags,
	// including []*multipart.FileHeader for file inputs.
	if err := c.BodyParser(&project); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid form data"})
	}

	validate := validator.New()
	if err := validate.Struct(project); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// sanitize project name into a filesystem-safe folder name
	dirName := sanitizeFolderName(project.ProjectName)
	projectDir := filepath.Join("uploads", dirName)

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create project directory"})
	}

	var savedPaths []string
	for _, fh := range project.Images {
		filename := uniqueFilename(fh.Filename)
		dest := filepath.Join(projectDir, filename)

		if err := c.SaveFile(fh, dest); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to save " + fh.Filename})
		}
		savedPaths = append(savedPaths, dest)
	}

	project.ImagePaths = strings.Join(savedPaths, ",")

	// ... insert project into DB here (ProjectId, UserId, ProjectName, Description, ImagePaths)

	return c.JSON(fiber.Map{"message": "project created", "paths": savedPaths})
}

// strips anything that isn't safe for a directory name — no spaces, slashes, dots-only, etc.
func sanitizeFolderName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	name = reg.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	if name == "" {
		name = "project"
	}
	return name
}

// avoids overwrites if two uploaded files share a name, or two projects share a sanitized folder name
func uniqueFilename(original string) string {
	ext := filepath.Ext(original)
	base := strings.TrimSuffix(filepath.Base(original), ext)
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%s-%s%s", base, hex.EncodeToString(b), ext)
}








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