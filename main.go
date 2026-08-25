package main

import (
	"crow/frontend/mark/comp"
	"crow/frontend/mark/pages"
	"crow/help"
	"crow/help/strs"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)



func main() {


	//TODO:=CREATE POLLING SYSTEM FOR  NOTIFICATIONS CREATE DEV INTERFACE FOR RESET PASSWORDs (dev and farid), CREATE FORMS FOR EVERYTHING

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
		case "Projets":
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

	//contact
	app.Post("/contact", help.Authmid, func(c fiber.Ctx) error {
		contact := new(strs.Contact)
		if err:=c.Bind().Form(contact); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(contact); err!=nil {
			return c.SendString(help.ShowError(err))
		} 
		var exists bool
		err := db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM contact)`)
		if err != nil {
		    return c.SendString(help.ShowError(err))
		}
		if exists {
			_, err = db.Exec(
			    `UPDATE contact SET address = ?, baladya = ?, wilaya = ?, email = ?, number = ?, location = ?`,
			    contact.Address, contact.Baladya, contact.Wilaya, contact.Email, contact.Number, contact.Location,
			)
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
		}else{
			_, err := db.Exec(
			    `INSERT INTO contact (address, baladya, wilaya, email, number, location) VALUES (?, ?, ?, ?, ?, ?)`,
			    contact.Address, contact.Baladya, contact.Wilaya, contact.Email, contact.Number, contact.Location,
			)
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	//messages
	app.Post("/uploadMessage", func (c fiber.Ctx) error {
		message := new(strs.Message)
		if err:=c.Bind().Form(message); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(message); err!=nil {
			return c.SendString(help.ShowError(err))
		} 
		_,err := db.Exec(`INSERT INTO messages (firstName, lastName, email, object, message) VALUES (?,?,?,?,?)`,
						message.FirstName, message.LastName, message.Email, message.Email, message.Object, message.Message)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return nil
	})



	//users
	app.Post("/createUser", func(c fiber.Ctx) error {
		err:= help.Signup(c, db)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	app.Post("/updateUser", help.Authmid, func(c fiber.Ctx) error {
		user:= c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		nUser:=new(strs.User)
		if err:=c.Bind().Form(nUser); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(nUser); err!=nil {
			return c.SendString(help.ShowError(err))
		} 
		if bytes, err := bcrypt.GenerateFromPassword([]byte(nUser.Password), 12); err!= nil{
			return err
		} else {
			nUser.Password = string(bytes)
		}
		_,err := db.Exec(`UPDATE users SET userName = ?, password = ?, WHERE userId = ?`, nUser.UserName, nUser.Password, nUser.UserId)
		if err!=nil{
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil
		
	})

	app.Post("/deleteUser", help.Authmid, func (c fiber.Ctx) error {
		user:= c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		userID := c.FormValue("userId")
		if _, err := db.Exec(`DELETE FROM users WHERE userID = ?`, userID); err!= nil{
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	app.Post("/reset", help.Authmid, func (c fiber.Ctx) error {
		user:= c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "sss"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		password:= c.FormValue("password")
		if bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12); err!= nil{
			return c.SendString(help.ShowError(err))
		} else {
			password = string(bytes)
		}
		if _, err:=db.Exec(`UPDATE users SET password = ? WHERE access = ?`, password, "fcmd"); err!= nil{
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	//categories
	app.Post("/createCategory", help.Authmid, func(c fiber.Ctx) error {
		user:= c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		category:=new(strs.Category)
		if err:=c.Bind().Form(category); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(category); err!=nil {
			return c.SendString(help.ShowError(err))
		} 

		path, err:= help.SaveImage(c, category.Image)
		if err!= nil {
			return c.SendString(help.ShowError(err))
		}
		category.ImagePath = path
		_,err = db.NamedExec(`INSERT INTO categories (categoryName, imagePath) VALUES (:categoryName, :imagePath)`, category)
		if err!= nil {
			if err:= help.DeleteImage(category.ImagePath); err!=nil {
				return c.SendString(help.ShowError(err))
			}
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	app.Post("deleteCategory", help.Authmid, func (c fiber.Ctx) error {
		user:= c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		category := c.FormValue("category")
		_,err := db.Exec("DELETE FROM categoryies WHERE categoryName = $1", category)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	app.Post("updateCategory", help.Authmid, func (c fiber.Ctx) error {
		user:= c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		category:=new(strs.Category)
		if err:=c.Bind().Form(category); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(category); err!=nil {
			return c.SendString(help.ShowError(err))
		} 
		oldPath := category.ImagePath
		path, err:= help.SaveImage(c, category.Image)
		if err!=nil{
			return c.SendString(help.ShowError(err)) 
		}
		category.ImagePath = path
		_, err = db.Exec(`UPDATE categories SET categoryName = ?, imagePath = ? WHERE categoryId = ?`,
						category.CategoryName, category.ImagePath,  category.CategoryId)
		if err!=nil {
			err = help.DeleteImage(category.ImagePath)
			return c.SendString(help.ShowError(err))
		}
		if err := help.DeleteImage(oldPath); err!=nil{
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	//partnairs
	app.Post("/uploadPartenair", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		partenair:=new(strs.Partenair)
		if err:=c.Bind().Form(partenair); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(partenair); err!=nil {
			return c.SendString(help.ShowError(err))
		} 
		if path, err:= help.SaveImage(c, partenair.PartenairImg); err!=nil {
			return c.SendString(help.ShowError(err))
		} else {
			partenair.ImagePath = path
		}
		_,err := db.Exec(`INSERT INTO partenairs(partenairName, imagePath) VALUES(?, ?)`, partenair.PartenairName, partenair.ImagePath)
		if err!=nil {
			err= help.DeleteImage(partenair.ImagePath)
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	app.Post("/deletePartenair", help.Authmid, func (c fiber.Ctx) error {
		partenairId := c.FormValue("partenairId")
		partenairImagePath := c.FormValue("imagePath")
		_,err := db.Exec(`DELETE FROM partenairs WHERE PartenairID = ?`, partenairId)
		if err!=nil{
			return c.SendString(help.ShowError(err))
		}
		if err := help.DeleteImage(partenairImagePath); err!=nil{
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil
	})

	
	

	//upload project
	app.Post("/uploadProject", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "c"){
			return c.Status(403).SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		project := new(strs.Project)
		if err:=c.Bind().Form(project); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(project); err!=nil {
			return c.SendString(help.ShowError(err))
		} 

		//userid
		project.UserId = user.UserId

		//main image path
		mPath, err := help.SaveImage(c, project.MImage)
		if err != nil {
			return c.SendString(help.ShowError(err))
		}
		project.MImagePath = mPath

		//image paths
		for _, image := range project.Images {
			path, err := help.SaveImage(c, image)
			if err!= nil {
				if err:= help.DeleteImage(project.MImagePath); err!=nil {
					return c.SendString(help.ShowError(err))
				}
				project.ImagePaths = strings.TrimSuffix(project.ImagePaths, ";")
				toDelete := strings.Split(project.ImagePaths, ";")
				for _, v := range toDelete {
					if err :=help.DeleteImage(v); err!= nil {
						return c.SendString(help.ShowError(err))
					}
				}
				return c.SendString(help.ShowError(err))
			}
			project.ImagePaths = project.ImagePaths + path + ";"
		}
		project.ImagePaths = strings.TrimSuffix(project.ImagePaths, ";")

		//enter into te database
		_, err = db.NamedExec(
		    `INSERT INTO projects (userId, categoryId, projectName, description, imagePaths, mImagePath)
		     VALUES (:userId, :categoryId, :projectName, :description, :imagePaths, :mImagePath)`,
		    project,
		)
		if err != nil {
			toDelete := strings.Split(project.ImagePaths, ";")
			for _, v := range toDelete {
				if err :=help.DeleteImage(v); err!= nil {
					return c.SendString(help.ShowError(err))
				}
			}
			if err:= help.DeleteImage(project.MImagePath); err!=nil {
				return c.SendString(help.ShowError(err))
			}
		    return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return nil

	})


	app.Post("/alpinetest", func(c fiber.Ctx) error {
		return help.Render(c, comp.Alpinetest())
	})

	app.Post("/restartDatabase", func(c fiber.Ctx) error {
		help.DeleteTable(db)
		return c.SendString("restarted")
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