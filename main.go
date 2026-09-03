package main

import (
	"crow/frontend/mark/comp"
	"crow/frontend/mark/comp/pageComp"
	"crow/frontend/mark/pages"
	"crow/help"
	"crow/help/strs"
	"fmt"
	"strings"
	"strconv"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)



func main() {


	//TODO:= projects / logs 

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

	app.Get("/partenairs", func(c fiber.Ctx) error {
		var partenairs []strs.Partenair
		err := db.Select(&partenairs, `SELECT * FROM partenairs`)
		if err != nil {
			return c.SendString(help.ShowError(err))
		}
		return help.Hrender(c, pages.PartenairFace(partenairs))
	})

	//atelier
	app.Get("/atelier", func(c fiber.Ctx) error {
		var members []strs.Member
		err := db.SelectContext(c, &members, `SELECT * FROM members`)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		return help.Hrender(c, pages.Atelier(members))
	})

	//projets
	app.Get("/projets", func(c fiber.Ctx) error {
		var list []strs.ProjectWithCategory
		query := `
		    SELECT 
		        p.*,
		        c.categoryName
		    FROM projects p
		    LEFT JOIN categories c ON p.categoryId = c.categoryId
		    ORDER BY p.projectId DESC LIMIT 4
		`
		err := db.SelectContext(c.Context(), &list, query)
		if err != nil {
		    return err
		}
		var categories []strs.Category
		err = db.SelectContext(c, &categories, `SELECT * FROM categories`)
		if err!= nil {
			return c.SendString(help.ShowError(err))
		}
		con := len(list) > 3
		if con {
			list = list[:3]
		}
		return help.Hrender(c, pages.Projets(list, categories, 3, 0, con, 0))
	})


	app.Get("/projets/category/:category", func(c fiber.Ctx) error {
		category, _ := strconv.Atoi(c.Params("category"))
		var list []strs.ProjectWithCategory
		query := `
		    SELECT 
		        p.*,
		        c.categoryName
		    FROM projects p
		    LEFT JOIN categories c ON p.categoryId = c.categoryId
			WHERE p.categoryId = ?
		    ORDER BY p.projectId DESC LIMIT 4
		`
		err := db.SelectContext(c.Context(), &list, query, category)
		if err != nil {
		    return err
		}
		var categories []strs.Category
		err = db.SelectContext(c, &categories, `SELECT * FROM categories`)
		if err!= nil {
			return c.SendString(help.ShowError(err))
		}
		con := len(list) > 3
		if con {
			list = list[:3]
		}
		return help.Hrender(c, pages.Projets(list, categories, 3, 0, con, category))
		
	})

	app.Get("/projets/:id", func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		var project strs.ProjectWithCategory
		query := `
		    SELECT 
		        p.*,
		        c.categoryName
		    FROM projects p
		    LEFT JOIN categories c ON p.categoryId = c.categoryId
			WHERE p.projectId = ?
		    ORDER BY p.projectId DESC LIMIT 4
		`
		err := db.GetContext(c.Context(), &project, query, id)
		if err != nil {
		    return err
		}
		return help.Hrender(c, pages.OneProject(project))
	})

	app.Post("/moreProjects", func(c fiber.Ctx) error {
		page := new(strs.Page)
		if err:= c.Bind().Form(page); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		page.Offset = page.Offset + 1
		var list []strs.ProjectWithCategory
		if page.Category == 0 {
			query := `
			    SELECT 
			        p.*,
			        c.categoryName
			    FROM projects p
			    LEFT JOIN categories c ON p.categoryId = c.categoryId
			    ORDER BY p.projectId DESC LIMIT ? OFFSET ?
			`
			err := db.SelectContext(c.Context(), &list, query, page.Limit + 1, page.Limit * page.Offset)
			if err != nil {
			    return err
			}
		}else{
			query := `
			    SELECT 
			        p.*,
			        c.categoryName
			    FROM projects p
			    LEFT JOIN categories c ON p.categoryId = c.categoryId
				WHERE p.categoryId = ?
			    ORDER BY p.projectId DESC LIMIT ? OFFSET ?
			`
			err := db.SelectContext(c.Context(), &list, query, page.Category, page.Limit + 1, page.Limit * page.Offset)
			if err != nil {
			    return err
			}
		}
		con := len(list) > page.Limit
		if con {
			list = list[:page.Limit]
		}
		return help.Render(c, comp.ProjectSection(list, page.Limit, page.Offset, con, page.Category))
	})




	//book
	app.Get("/book", func(c fiber.Ctx) error {
		return help.Hrender(c, pages.Book())
	})

	//candidature
	app.Get("/candidature", func(c fiber.Ctx) error{
		var software []strs.Software
		err:= db.Select(&software, `SELECT * FROM Software`)
		if err!=nil{
			c.SendString(help.ShowError(err))
		}
		return help.Hrender(c, pages.Candidature(software))
	})

	//contact
	app.Get("/contact", func(c fiber.Ctx) error {
		contact := new(strs.Contact)

		// Use db.Get instead of db.Exec to query a single struct row
		err := db.Get(contact, `SELECT address, baladya, wilaya, email, number, fax, location, image FROM contact LIMIT 1`)
		if err != nil {
		    return c.SendString(help.ShowError(err))
		}
		return help.Hrender(c, pages.Contact(contact))
	})

	//rating
	app.Get("/star", func(c fiber.Ctx) error {
	    var stars []strs.Rating
	    err := db.SelectContext(c.Context(), &stars, `SELECT * FROM ratings WHERE approve = ? ORDER BY ratingId DESC LIMIT 3`, 1)
	    if err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	    var avg float64
	    err = db.GetContext(c.Context(), &avg, `SELECT COALESCE(AVG(rating), 0) FROM ratings WHERE approve = ?`, 1)
	    if err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	    hasMore := len(stars) > 2
		if hasMore {
    	    stars = stars[:2]
    	}
	    return help.Hrender(c, pages.Star(stars, hasMore, avg))
	})

	app.Post("/moreRatings", func(c fiber.Ctx) error {
		limit, err := strconv.Atoi(c.FormValue("limit"))
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		offset, err := strconv.Atoi(c.FormValue("offset"))
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		offset = offset + 1
	    var stars []strs.Rating
	    err = db.SelectContext(c.Context(), &stars, `SELECT * FROM ratings WHERE approve = ? ORDER BY ratingId DESC LIMIT ? OFFSET ?`, 1, limit+1, offset*limit)
	    if err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	    hasMore := len(stars) > limit
		if hasMore {
    	    stars = stars[:limit]
    	}
	    return help.Render(c, pageComp.StarsPage(stars, offset, limit, hasMore))
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
		if user.Access == "sss" {
			return help.Hrender(c, pages.Restart(user))
		}
		return help.Hrender(c, pages.Dash(user))
	})

	app.Post("/dashNav", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		target := c.FormValue("target")
		switch target {
		case "Contact":
			contact := new(strs.Contact)
			err := db.Get(contact, `SELECT address, baladya, wilaya, email, number, fax, location, image FROM contact LIMIT 1`)
			if err != nil{
				return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.ContactMod(user, contact))
		case "Candidats":
			var candidats []strs.JobApplication
			err := db.SelectContext(c, &candidats, "SELECT * FROM jobApplications ORDER BY apID DESC")
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
			var notSeen []strs.CanNot
			err = db.SelectContext(c, &notSeen, "SELECT * FROM canNot ORDER BY canNotId DESC")
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.CandidatsMod(user, candidats, notSeen))
		case "Atelier":
			var members []strs.Member
			err := db.SelectContext(c, &members, "SELECT * FROM members ORDER BY memberID DESC")
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.AtelierMod(user, members))
		case "Comptes":
			var users []strs.User
			err := db.SelectContext(c, &users, "SELECT * FROM users ORDER BY userID DESC")
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.ComptesMod(user, users))
		case "Projets":
			var categories []strs.Category
			var projects []strs.Project
			err := db.SelectContext(c, &categories, "SELECT * FROM categories")
			if err!=nil {
				return c.SendString(help.ShowError(err))
			}
			err = db.SelectContext(c, &projects, `SELECT * FROM projects`)
			if err!=nil {
				return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.ProjetsMod(user, categories, projects))
		case "Stars":
			var stars []strs.Rating
			err:= db.Select(&stars, `SELECT * FROM ratings ORDER BY ratingId DESC`)
			if err!=nil {
				return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.StarsMod(user, stars))
		case "Messages":
			var messages []strs.Message
			err := db.SelectContext(c, &messages, "SELECT * FROM messages ORDER BY messageID DESC")
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
			var notSeen []strs.MessNot
			err = db.SelectContext(c, &notSeen, "SELECT * FROM messNot ORDER BY messageID DESC")
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.MessagesMod(user, messages, notSeen))
		case "Logout":
			if err:=help.Logout(c); err!=nil{
				return c.SendString(help.ShowError(err))
			} else {
				return help.Redirect(c, "/")
			}
		case "Categories":
			var categories []strs.Category
			err := db.SelectContext(c, &categories, "SELECT * FROM categories")
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.CategoriesMod(user, categories))
		case "Partenaire":
			var partenairs []strs.Partenair
			err := db.SelectContext(c, &partenairs, "SELECT * FROM partenairs")
			if err != nil {
			    return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.PartenairsMod(user, partenairs))
		
		case "Logiciels":
			var soft []strs.Software
			err:= db.Select(&soft, `SELECT * FROM software ORDER BY softwareId DESC`)
			if err!=nil {
				return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.SoftwareMod(user, soft))	
		case "Logs":
			var logs []strs.Log
			err:= db.Select(&logs, `SELECT * FROM logs ORDER BY logId DESC`)
			if err!=nil {
				return c.SendString(help.ShowError(err))
			}
			return help.Render(c, comp.LogsMod(logs))	
		}
		return c.SendString("page introuvable")
	})

	//contact
	app.Post("/contact", help.Authmid, func(c fiber.Ctx) error {
	    user := c.Locals("user").(*strs.User)
	    if !strings.Contains(user.Access, "f") {
	        return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
	    }

	    contact := new(strs.Contact)
	    if err := c.Bind().Form(contact); err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	    if err := help.Validate.Struct(contact); err != nil {
	        return c.SendString(help.ShowError(err))
	    }

	    var oldpath string
	    var hasNewImage bool

	    if contact.Image != nil {
	        err := db.Get(&oldpath, `SELECT image FROM contact LIMIT 1`)
	        if err != nil  {
	            return c.SendString(help.ShowError(err))
	        }

	        path, err := help.SaveImage(c, contact.Image)
	        if err != nil {
	            return c.SendString(help.ShowError(err))
	        }
	        contact.ImagePath = path
	        hasNewImage = true
	    }

	    _, err := db.Exec(
	        `UPDATE contact SET address = ?, baladya = ?, wilaya = ?, email = ?, number = ?, fax = ?, location = ?, image = ?`,
	        contact.Address, contact.Baladya, contact.Wilaya, contact.Email, contact.Number, contact.Fax, contact.Location, contact.ImagePath,
	    )
	    if err != nil {
	        return c.SendString(help.ShowError(err))
	    }

	    if hasNewImage && oldpath != "" {
	        if err2 := help.DeleteImage(oldpath); err2 != nil {
	            return c.SendString(help.ShowError(err2))
	        }
	    }

	    c.Set("HX-Trigger", "success")
	    return c.SendString("")
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
		return c.SendString("")
	})

	app.Post("/deleteMessage", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		messageID:= c.FormValue("messageID")
		_,err:= db.Exec(`DELETE FROM messages WHERE messageID = ?`, messageID)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendString("")
	})

	app.Post("/viewMessage", help.Authmid, func (c fiber.Ctx) error  {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		notId:= c.FormValue("notId")
		_,err:= db.Exec(`DELETE FROM messNot WHERE messNotId = ?`, notId)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendString("")
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
		return c.SendString("")
	})

	app.Post("/updateUser", help.Authmid, func(c fiber.Ctx) error {
	    user := c.Locals("user").(*strs.User)
	    if !strings.Contains(user.Access, "f") {
	        return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
	    }
	
	    type UpdateUserReq struct {
	        UserId   int    `form:"UserId" validate:"required"`
	        UserName string `form:"UserName" validate:"required,min=3,max=50"`
	        Access   string `form:"Access" validate:"max=4"`
	    }
	
	    req := new(UpdateUserReq)
	    if err := c.Bind().Form(req); err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	    if err := help.Validate.Struct(req); err != nil {
	        return c.SendString(help.ShowError(err))
	    }

		if req.UserId == 1 {
			req.Access = "fcmd"
		}
	
	    _, err := db.Exec(`UPDATE users SET userName = ?, access = ? WHERE userId = ?`, req.UserName, req.Access, req.UserId)
	    if err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	
	    c.Set("HX-Trigger", "success")
	    return c.SendString("")
	})
	
	// 2. Update Password Only
	app.Post("/updateUserPassword", help.Authmid, func(c fiber.Ctx) error {
	    user := c.Locals("user").(*strs.User)
	    if !strings.Contains(user.Access, "f") {
	        return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
	    }
	
	    type UpdatePasswordReq struct {
	        UserId   int    `form:"UserId" validate:"required"`
	        Password string `form:"Password" validate:"required,min=6,max=12"`
	    }
	
	    req := new(UpdatePasswordReq)
	    if err := c.Bind().Form(req); err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	    if err := help.Validate.Struct(req); err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	
	    bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	    if err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	
	    _, err = db.Exec(`UPDATE users SET password = ? WHERE userId = ?`, string(bytes), req.UserId)
	    if err != nil {
	        return c.SendString(help.ShowError(err))
	    }
	
	    c.Set("HX-Trigger", "success")
	    return c.SendString("")
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
		return c.SendString("")
	})

	app.Post("/reset", help.Authmid, func (c fiber.Ctx) error {
		user:= c.Locals("user").(*strs.User)
		if user.Access != "sss" {
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
		if err:=help.Logout(c); err!=nil{
			return c.SendString(help.ShowError(err))
		} else {
			return help.Redirect(c, "/")
		}
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
		return c.SendString("")
	})

	app.Post("/deleteCategory", help.Authmid, func (c fiber.Ctx) error {
		user := c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f") {
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		categoryId := c.FormValue("categoryId")
	
		var imagePath string
		if err := db.Get(&imagePath, `SELECT imagePath FROM categories WHERE categoryId = ?`, categoryId); err != nil {
			return c.SendString(help.ShowError(err))
		}
	
		if _, err := db.Exec(`DELETE FROM categories WHERE categoryId = ?`, categoryId); err != nil {
			return c.SendString(help.ShowError(err))
		}
	
		if err := help.DeleteImage(imagePath); err != nil {
			return c.SendString(help.ShowError(err))
		}
	
		c.Set("HX-Trigger", "success")
		return c.SendString("")
	})

	app.Post("/updateCategory", help.Authmid, func (c fiber.Ctx) error {
		user := c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f") {
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		category := new(strs.Category)
		if err := c.Bind().Form(category); err != nil {
			return c.SendString(help.ShowError(err))
		}
		if err := help.Validate.Struct(category); err != nil {
			return c.SendString(help.ShowError(err))
		}

		var oldPath string
		if err := db.Get(&oldPath, `SELECT imagePath FROM categories WHERE categoryId = ?`, category.CategoryId); err != nil {
			return c.SendString(help.ShowError(err))
		}

		if category.Image != nil {
			path, err := help.SaveImage(c, category.Image)
			if err != nil {
				return c.SendString(help.ShowError(err))
			}
			category.ImagePath = path

			_, err = db.Exec(`UPDATE categories SET categoryName = ?, imagePath = ? WHERE categoryId = ?`,
				category.CategoryName, category.ImagePath, category.CategoryId)
			if err != nil {
				if err2 := help.DeleteImage(category.ImagePath); err2 != nil {
					return c.SendString(help.ShowError(err2))
				}
				return c.SendString(help.ShowError(err))
			}
			if err := help.DeleteImage(oldPath); err != nil {
				return c.SendString(help.ShowError(err))
			}
		} else {
			_, err := db.Exec(`UPDATE categories SET categoryName = ? WHERE categoryId = ?`,
				category.CategoryName, category.CategoryId)
			if err!=nil {
				return c.SendString(help.ShowError(err))
			}
		}
		c.Set("HX-Trigger", "success")
		return c.SendString("")
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
		return c.SendString("")
	})

	app.Post("/deletePartenair", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		partenairId := c.FormValue("partenairId")
		var partenairImagePath string
		if err:= db.Get(&partenairImagePath, `SELECT imagePath FROM partenairs WHERE partenairId = ?`, partenairId) ; err!= nil {
			return c.SendString(help.ShowError(err))
		}
		_,err := db.Exec(`DELETE FROM partenairs WHERE PartenairID = ?`, partenairId)
		if err!=nil{
			return c.SendString(help.ShowError(err))
		}
		if err := help.DeleteImage(partenairImagePath); err!=nil{
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendString("")
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
		_,err:= db.Exec(`INSERT INTO ratings (name, comment, rating, approve) VALUES (?,?,?,?)`,
						rating.Name, rating.Comment, rating.Rating, 0)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendString("")
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
		return c.SendString("")
	})

	app.Post("/approveRate", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		ratingId:= c.FormValue("ratingId")
		_,err := db.Exec(`UPDATE ratings SET approve = ? WHERE ratingId = ?`, 1, ratingId)
		if err!= nil {
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendString("")
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
		return c.SendString("")
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

		if member.MemberImage != nil {
			path, err:= help.SaveImage(c, member.MemberImage)
			if err!=nil {
				return c.SendString(help.ShowError(err))
			}
			member.MemberImagePath = path
			var oldPath string 
			err= db.Get(&oldPath, `SELECT memberImagePath FROM members WHERE memberID = ?`, member.MemberId)
			if err!=nil {
				return c.SendString(help.ShowError(err))
			}
			_,err = db.Exec(`UPDATE members SET memberName=?, memberTitle=?, memberDescription=?, memberImagePath=? WHERE memberId = ?`,
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

		} else {
			_,err:= db.Exec(`UPDATE members SET memberName=?, memberTitle=?, memberDescription=? WHERE memberId = ?`,
						member.MemberName, member.MemberTitle, member.MemberDescription, member.MemberId)
			if err != nil {
				return c.SendString(help.ShowError(err))
			}
		}

		c.Set("HX-Trigger", "success")
		return c.SendString("")
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
		return c.SendString("")
	})

	//job applications
	app.Post("/apply", func (c fiber.Ctx) error  {
		//logic
		apply:= new(strs.JobApplication)
		if err:=c.Bind().Form(apply); err!=nil {
			return c.SendString(err.Error())
		}
		if err:= help.Validate.Struct(apply); err!=nil {
			return c.SendString(err.Error())
		} 

		//pdf
		cvPath, err:= help.SavePDF(c, apply.Cv)
		if err !=nil {
			return c.SendString(err.Error())
		}
		letterPath, err := help.SavePDF(c, apply.Letter)
		if err !=nil {
			if err2:= help.DeletePDF(cvPath); err2!=nil {
				return c.SendString(help.ShowError(err2))
			}
			return c.SendString(err.Error())
		}
		
		apply.CvPath = cvPath
		apply.LetterPath = letterPath

		res,err := db.Exec(`INSERT INTO jobApplications (firstName, lastName, email, object, message, software, cv, letter) VALUES (?,?,?,?,?,?,?,?)`,
						apply.FirstName, apply.LastName, apply.Email, apply.Object, apply.Message, apply.Software, apply.CvPath,apply.LetterPath)
		if err!=nil {
			return c.SendString(err.Error())
		} else {
			id, err := res.LastInsertId()
			if err!=nil {
				return c.SendString(err.Error())
			}
			_,err = db.Exec(`INSERT INTO canNot (candidatId) VALUES (?)`, id)
			if err!=nil {
				return c.SendString(err.Error())
			}
		}
		c.Set("HX-Trigger", "success")
		return c.SendString("")
	})

	app.Post("/deleteApply", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		id:= c.FormValue("applyid")
		cv:= c.FormValue("cv")
		letter:=c.FormValue("letter")
		if _,err := db.Exec(`DELETE FROM jobApplications WHERE apId = ?`, id); err!=nil{
			return c.SendString(help.ShowError(err))
		}
		if err:= help.DeletePDF(cv); err!=nil {
			return c.SendString(help.ShowError(err))
		}

		if err:= help.DeletePDF(letter); err!=nil {
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendString("")
	})

	app.Post("/viewApply", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		id:= c.FormValue("notifId")
		if _,err := db.Exec(`DELETE FROM canNot WHERE canNotId = ?`, id); err!=nil{
			return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendString("")
	})

	//add software
	app.Post("/addSoftware", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		software:= new(strs.Software)
		if err:=c.Bind().Form(software); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(software); err!=nil {
			return c.SendString(help.ShowError(err))
		} 

		_,err := db.Exec(`INSERT INTO software(softwareName, required) VALUES (?,?)`,
						software.SoftwareName, software.Required)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		
		c.Set("HX-Trigger", "success")
		return c.SendString("")
	})


	//app.Post("/updateSoftware", help.Authmid, func (c fiber.Ctx) error {
	//	user:=c.Locals("user").(*strs.User)
	//	if !strings.Contains(user.Access, "f"){
	//		return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
	//	}
//
	//	software:= new(strs.Software)
	//	if err:=c.Bind().Form(software); err!=nil {
	//		return c.SendString(help.ShowError(err))
	//	}
	//	if err:= help.Validate.Struct(software); err!=nil {
	//		return c.SendString(help.ShowError(err))
	//	} 
//
	//	_,err := db.Exec(`UPDATE software SET softwareName=?, required=? WHERE SoftwareId = ?`,
	//					software.SoftwareName, software.Required, software.SoftwareId)
	//	if err!=nil {
	//		return c.SendString(help.ShowError(err))
	//	}
	//	
	//	c.Set("HX-Trigger", "success")
	//	return c.SendString("")
	//})

	app.Post("/deleteSoftware", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		id:= c.FormValue("softwareid")
		if _,err := db.Exec(`DELETE FROM software WHERE softwareId = ?`, id); err!=nil{
			return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendString("")
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
			project.ImagePaths = path + ";" + project.ImagePaths
		}
		project.ImagePaths = strings.TrimSuffix(project.ImagePaths, ";")

		//enter into te database
		_, err = db.ExecContext(
		    c.Context(),
		    `INSERT INTO projects (
		        userId, categoryId, projectName, description, 
		        imagePaths, mImagePath, date, maitre, emplacement, programme
		    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		    project.UserId,
		    project.CategoryId,
		    project.ProjectName,
		    project.Description,
		    project.ImagePaths,
		    project.MImagePath,
		    project.Date,
		    project.Maitre,
		    project.Emplacement,
		    project.Programme,
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

		message := fmt.Sprintf("%s a créé un nouveau projet : %s", user.UserName, project.ProjectName)
		if _, err := db.ExecContext(c.Context(), `INSERT INTO logs(message, type) VALUES (?,?)`, message,1); err != nil {
		    return c.SendString(help.ShowError(err))
		}

		c.Set("HX-Trigger", "success")
		return c.SendString("")

	})

	app.Post("/updateProject", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "m"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}
		
		//logic here
		project := new(strs.ProjectMod)
		if err:=c.Bind().Form(project); err!=nil {
			return c.SendString(help.ShowError(err))
		}
		if err:= help.Validate.Struct(project); err!=nil {
			return c.SendString(help.ShowError(err))
		} 

		oldProject:= new(strs.Project)
		err:= db.Get(oldProject, `SELECT * FROM projects WHERE projectId = ?`, project.ProjectId)
		if err !=nil {
			return c.SendString(help.ShowError(err))
		}
		var toDelete []string
		var comPaths []string
		//main image
		if project.MImage != nil {
			toDelete = append(toDelete, oldProject.MImagePath)
			project.MImagePath, err = help.SaveImage(c, project.MImage)
			if err !=nil {
				return c.SendString(help.ShowError(err))
			}
			comPaths = append(comPaths, project.MImagePath)
		} else {
			project.MImagePath = oldProject.MImagePath
		}
		//images
		oldPaths:= strings.Split(oldProject.ImagePaths, ";")
		newPaths:= strings.Split(project.ImagePaths, ";")
		
		for _,v := range oldPaths {
			found:=false
			for _,v2:=range newPaths {
				if v == v2{
					found = true
					break
				}
			}
			if found==false {
				toDelete = append(toDelete, v)
			}
		}
		
		for _, image := range project.Images {
			path, err := help.SaveImage(c, image)
			comPaths = append(comPaths, path)
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
			project.ImagePaths = path + ";" + project.ImagePaths
		}
		project.ImagePaths = strings.TrimSuffix(project.ImagePaths, ";")
		
		_, err = db.ExecContext(
		    c.Context(),
		    `UPDATE projects SET 
		        userId = ?,
		        categoryId = ?,
		        projectName = ?,
		        description = ?,
		        imagePaths = ?,
		        mImagePath = ?,
		        date = ?,
		        maitre = ?,
		        emplacement = ?,
		        programme = ?
		    WHERE projectId = ?`,
		    project.UserId,
		    project.CategoryId,
		    project.ProjectName,
		    project.Description,
		    project.ImagePaths,
		    project.MImagePath,
		    project.Date,
		    project.Maitre,
		    project.Emplacement,
		    project.Programme,
		    project.ProjectId,
		)
		if err != nil {
			for _,v:= range comPaths {
				if err2:= help.DeleteImage(v); err2!=nil {
					return c.SendString(help.ShowError(err2))
				}
			}
			return c.SendString(help.ShowError(err))
		}else {
			for _,v:=range toDelete {
				if err2:= help.DeleteImage(v); err2!=nil {
					return c.SendString(help.ShowError(err2))
				}
			}
		}

		message := fmt.Sprintf("%s a modifié un projet : %s", user.UserName, project.ProjectName)
		if _, err := db.ExecContext(c.Context(), `INSERT INTO logs(message, type) VALUES (?,?)`, message, 2); err != nil {
		    return c.SendString(help.ShowError(err))
		}
		c.Set("HX-Trigger", "success")
		return c.SendString("")
	})

	app.Post("/deleteProject", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "d"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		id:= c.FormValue("projectId")
		//get from database and delete
		var mPath string
		var sPaths string
		var projectName string
		err := db.Get(&mPath, `SELECT mImagePath FROM projects WHERE projectId = ?`, id)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		err = db.Get(&sPaths, `SELECT imagePaths FROM projects WHERE projectId = ?`, id)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		err = db.Get(&projectName, `SELECT projectName FROM projects WHERE projectId = ?`, id)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}

		paths:=strings.Split(sPaths, ";" )

		_,err= db.Exec(`DELETE FROM projects WHERE projectId=?`, id)
		if err!=nil{
			return c.SendString(help.ShowError(err))
		}
		err = help.DeleteImage(mPath)
		if err!=nil {
			return c.SendString(help.ShowError(err))
		}
		for _,v := range paths {
			err = help.DeleteImage(v)
			if err!=nil {
				return c.SendString(help.ShowError(err))
			}
		}

		message := fmt.Sprintf("%s a supprimé un projet : %s", user.UserName, projectName)
		if _, err := db.ExecContext(c.Context(), `INSERT INTO logs(message, type) VALUES (?,?)`, message, 3); err != nil {
		    return c.SendString(help.ShowError(err))
		}
		
		c.Set("HX-Trigger", "success")
		return c.SendString("")
	})


	app.Post("/updateNotif", help.Authmid, func (c fiber.Ctx) error {
		user:=c.Locals("user").(*strs.User)
		if !strings.Contains(user.Access, "f"){
			return c.SendString("Vous n'avez pas le droit d'effectuer cette opération.")
		}

		var messages int = 0
		if err := db.Get(&messages, `SELECT COUNT(*) FROM messNot`); err!=nil {
			return c.SendStatus(500)
		}
		var candidats int = 0
		if err := db.Get(&candidats, `SELECT COUNT(*) FROM canNot`); err!=nil {
			return c.SendStatus(500)
		}

		trigger := fmt.Sprintf(`{"notif": {"messages": %d, "candidats": %d}}`, messages, candidats)
		c.Set("HX-Trigger", trigger)
		return c.SendString("")
	})

	//app.Post("/paginate", func (c fiber.Ctx) error {
	//	page:= new(strs.Page)
	//	if err:=c.Bind().Form(page); err!=nil {
	//		return c.SendString(help.ShowError(err))
	//	}
	//	switch page.Table {
	//	case 1:
	//		var projects []strs.Project
	//		if page.Category == "" {
	//			err := db.Select(&projects, `SELECT * FROM projects ORDER BY projectId DESC LIMIT ? OFFSET ?`, page.Limit, page.Offset) 
	//			if err != nil {
	//				return c.SendString(help.ShowError(err))
	//			}
	//		} else {
	//			err := db.Select(&projects, `SELECT * FROM projects WHERE categoryId = ? ORDER BY projectId DESC LIMIT ? OFFSET ?`, page.Category, page.Limit, page.Offset) 
	//			if err != nil {
	//				return c.SendString(help.ShowError(err))
	//			}
	//		}
	//		return help.Hrender(PageProjects(projects))
	//	case 2:
	//		var ratings []strs.Rating
	//		err := db.Select(&ratings, `SELECT * FROM ratings ORDER BY ratingId DESC LIMIT ? OFFSET ?`, page.Limit, page.Offset) 
	//		if err != nil {
	//			return c.SendString(help.ShowError(err))
	//		}
	//		return help.Hrender(PageRatings(ratings))
	//	}
	//	return c.SendString("")
	//})


	app.Post("/alpinetest", func(c fiber.Ctx) error {
		return help.Render(c, comp.Alpinetest())
	})

	app.Post("/restartDatabase", func(c fiber.Ctx) error {
		help.DeleteTable(db)
		return c.SendString("restarted")
	})

	app.Listen(":3000")
}