package main

import (
	"crow/frontend/mark/comp"
	"crow/frontend/mark/pages"
	"crow/help"
	"crow/help/strs"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)



func main() {


	//TODO:=CREATE DEV INTERFACE FOR RESET PASSWORDs (dev and farid), CREATE FORMS FOR EVERYTHING
	//projects, jobAppilcations, cannot
	//

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
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
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
		return c.SendStatus(200)
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
		res, err := db.Exec(`INSERT INTO messages (firstName, lastName, email, object, message) VALUES (?,?,?,?,?)`,
		    message.FirstName, message.LastName, message.Email, message.Object, message.Message)
		if err != nil { 
			return c.SendString(help.ShowError(err)) 
		}

		id, _ := res.LastInsertId()
		_, err = db.Exec(`INSERT INTO messNot (messageId) VALUES (?)`, id)
		if err!=nil {
			return c.SendString(help.ShowError(err)) 
		}

		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/deleteMessage", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		messageID:= c.FormValue("messageID")
		_,err:= db.Exec(`DELETE FROM messages WHERE id = ?`, messageID)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/viewMessage", help.Authmid, func (c fiber.Ctx) error  {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		notId:= c.FormValue("notId")
		_,err:= db.Exec(`DELETE FROM messNot WHERE id = ?`, notId)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	//users
	app.Post("/createUser", help.Authmid, func(c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		err:= help.Signup(c, db)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
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
			return c.SendString(help.ShowError(err))
		} else {
			nUser.Password = string(bytes)
		}
		_,err := db.Exec(`UPDATE users SET userName = ?, password = ? WHERE userId = ?`, nUser.UserName, nUser.Password, nUser.UserId)
		if err!=nil{
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
		
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
		return c.SendStatus(200)
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
		return c.SendStatus(200)
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
			if err2:= help.DeleteImage(category.ImagePath); err2!=nil {
				return c.SendString(help.ShowError(err2))
			}
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/deleteCategory", help.Authmid, func (c fiber.Ctx) error {
		user:= c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		category := c.FormValue("categoryId")
		_,err := db.Exec("DELETE FROM categories WHERE categoryId = ?", category)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/updateCategory", help.Authmid, func (c fiber.Ctx) error {
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
			if err2 := help.DeleteImage(category.ImagePath); err2!=nil {
				return c.SendString(help.ShowError(err2))
			}
			return c.SendString(help.ShowError(err))
		}
		if err := help.DeleteImage(oldPath); err!=nil{
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
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
			if err2:= help.DeleteImage(partenair.ImagePath); err2!=nil{
				return c.SendString(help.ShowError(err2))
			}
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/deletePartenair", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
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
		return c.SendStatus(200)
	})

	//ratings
	app.Post("/rate", func (c fiber.Ctx) error {
		rating:= new(strs.Rating)
		if err:=c.Bind().Form(rating); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(rating); err!=nil {
			return c.SendString(help.ShowError(err))
		} 	
		_,err:= db.Exec(`INSERT INTO ratings (name, comment, rating) VALUES (?,?,?)`,
						rating.Name, rating.Comment, rating.Rating * 10)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/deleteRate", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		ratingId:= c.FormValue("ratingId")
		_,err := db.Exec(`DELETE FROM ratings WHERE ratingId = ?`, ratingId)
		if err!= nil {
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	//members
	app.Post("/addMember", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		member:=new(strs.Member)
		if err:=c.Bind().Form(member); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(member); err!=nil {
			return c.SendString(help.ShowError(err))
		} 
		
		path, err := help.SaveImage(c, member.MemberImage)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}

		member.MemberImagePath = path
		_,err = db.Exec(`INSERT INTO members (memberName, memberTitle, memberDescription, memberImagePath) VALUES (?,?,?,?)`,
						member.MemberName, member.MemberTitle, member.MemberDescription, member.MemberImagePath)
		if err!=nil {
			if err2:= help.DeleteImage(member.MemberImagePath); err2!=nil {
				return c.SendString(help.ShowError(err2))
			}
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/updateMember", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		member:=new(strs.Member)
		if err:=c.Bind().Form(member); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(member); err!=nil {
			return c.SendString(help.ShowError(err))
		} 

		oldPath:= member.MemberImagePath

		if path, err:= help.SaveImage(c, member.MemberImage); err!=nil {
			return c.SendString(help.ShowError(err))
		} else {
			member.MemberImagePath = path
		}

		_,err:= db.Exec(`UPDATE members SET memberName=?, memberTitle=?, memberDescription=?, memberImagePath=? WHERE memberId = ?`,
						member.MemberName, member.MemberTitle, member.MemberDescription, member.MemberImagePath, member.MemberId)
		if err!=nil {
			if err2 := help.DeleteImage(member.MemberImagePath); err2!=nil {
				return c.SendString(help.ShowError(err2))
			}
			return c.SendString(help.ShowError(err))
		}
		if err := help.DeleteImage(oldPath); err!=nil{
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/deleteMember", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		id:=c.FormValue("memberId")
		path:=c.FormValue("memberImagePath")

		if _,err:= db.Exec(`DELETE FROM members WHERE memberID = ?`, id); err!= nil {
			return c.SendString(help.ShowError(err))
		}

		if err:= help.DeleteImage(path); err!=nil {
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	//job applications
	app.Post("/apply", func (c fiber.Ctx) error  {
		//logic
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/dlApp", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		//logic for later
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
	})

	app.Post("/viewAp", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		//logic for later
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)
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
				if v!= "" {
					if err :=help.DeleteImage(v); err!= nil {
						return c.SendString(help.ShowError(err))
					}
				}
			}
			if err:= help.DeleteImage(project.MImagePath); err!=nil {
				return c.SendString(help.ShowError(err))
			}
		    return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendStatus(200)

	})

	app.Post("/updateNotif", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		var messages int = -1
		if err := db.Get(&messages, `SELECT COUNT(*) FROM messNot`); err!=nil {
			return c.SendStatus(500)
		}
		var candidats int = -1
		if err := db.Get(&candidats, `SELECT COUNT(*) FROM canNot`); err!=nil {
			return c.SendStatus(500)
		}

		trigger := fmt.Sprintf(`{"notif": {"messages": %d, "candidats": %d}}`, messages, candidats)
		c.Set("HX-Trigger", trigger)
		return c.SendStatus(200)
	})


	app.Post("/alpinetest", func(c fiber.Ctx) error {
		return help.Render(c, comp.Alpinetest())
	})

	app.Post("/restartDatabase", func(c fiber.Ctx) error {
		help.DeleteTable(db)
		return c.SendString("restarted")
	})

	app.Listen(":3000")
}