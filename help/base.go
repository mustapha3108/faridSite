package help

import (
	"database/sql"
	"errors"
	"crow/frontend/mark/comp"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

//rendering
func Render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func Hrender(c fiber.Ctx, component templ.Component) error {
	if c.Get("HX-Request") == "true" {
		return Render(c, component)
	}
	return Render(c, comp.Layout(component))
}

func Redirect(c fiber.Ctx, path string) error {
    if c.Get("HX-Request") == "true" {
        c.Set("HX-Redirect", path)
        return c.SendStatus(200)
    }
    return c.Redirect().To(path)
}

func SaveImage(c fiber.Ctx, image *multipart.FileHeader) (string, error) {
    if !IsImage(image) {
        return "", errors.New("not supported image type")
    }
    ext := strings.ToLower(filepath.Ext(image.Filename))
    name := uuid.New().String() + ext
    path := filepath.Join("frontend", "glitter", "media", "dynamic", name)
    if err := c.SaveFile(image, path); err != nil {
        return "", err
    }
    return "media/dynamic/" + name, nil
}

//checking stuff and manual security 
const MaxImageSize = 6 * 1024 * 1024
func IsImage(file *multipart.FileHeader) bool {

	// ─── Size check ──────────────────────────────────────────────────────────
	if file.Size > MaxImageSize {
		return false
	}

	f, err := file.Open()
	if err != nil {
		return false
	}
	defer f.Close()

	// ─── Read enough bytes for all format checks ──────────────────────────────
	// 261 bytes covers all magic byte signatures including ZIP-based polyglots
	buf := make([]byte, 261)
	n, err := f.Read(buf)
	if err != nil || n < 8 {
		return false
	}
	buf = buf[:n]

	// ─── Magic byte checks ────────────────────────────────────────────────────
	allowed := []func([]byte) bool{
		func(b []byte) bool { // PNG
			return len(b) >= 8 &&
				b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 &&
				b[4] == 0x0D && b[5] == 0x0A && b[6] == 0x1A && b[7] == 0x0A
		},
		func(b []byte) bool { // JPEG / JPG
			return len(b) >= 3 &&
				b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF
		},
		//func(b []byte) bool { // GIF
		//	return len(b) >= 6 &&
		//		b[0] == 0x47 && b[1] == 0x49 && b[2] == 0x46 &&
		//		b[3] == 0x38 && (b[4] == 0x37 || b[4] == 0x39) && b[5] == 0x61
		//},
		func(b []byte) bool { // WEBP
			return len(b) >= 12 &&
				b[0] == 0x52 && b[1] == 0x49 && b[2] == 0x46 && b[3] == 0x46 &&
				b[8] == 0x57 && b[9] == 0x45 && b[10] == 0x42 && b[11] == 0x50
		},
		// func(b []byte) bool { // BMP
		// 	return len(b) >= 2 && b[0] == 0x42 && b[1] == 0x4D
		// },
		// func(b []byte) bool { // TIFF
		// 	return len(b) >= 4 &&
		// 		((b[0] == 0x49 && b[1] == 0x49 && b[2] == 0x2A && b[3] == 0x00) ||
		// 		 (b[0] == 0x4D && b[1] == 0x4D && b[2] == 0x00 && b[3] == 0x2A))
		// },
		// func(b []byte) bool { // SVG — disabled, can contain JS
		// 	return len(b) >= 5 &&
		// 		(string(b[:4]) == "<svg" || string(b[:5]) == "<?xml")
		// },
	}

	matched := false
	for _, check := range allowed {
		if check(buf) {
			matched = true
			break
		}
	}
	if !matched {
		return false
	}

	// ─── Polyglot checks ──────────────────────────────────────────────────────
	// reject ZIP-based polyglots (ZIP, JAR, DOCX etc embedded in image)
	if len(buf) >= 4 && buf[0] == 0x50 && buf[1] == 0x4B && buf[2] == 0x03 && buf[3] == 0x04 {
		return false
	}

	// reject PDF polyglots
	if len(buf) >= 4 && string(buf[:4]) == "%PDF" {
		return false
	}

	// reject PHP/script injection in the first bytes
	dangerousPatterns := []string{
		"<?php", "<?=", "<script", "<%", "#!/",
	}
	header := strings.ToLower(string(buf))
	for _, pattern := range dangerousPatterns {
		if strings.Contains(header, pattern) {
			return false
		}
	}

	// reject files with null bytes in unexpected places (common in exploit payloads)
	nullCount := 0
	for _, b := range buf[8:] {
		if b == 0x00 {
			nullCount++
		}
	}
	// allow some nulls (binary formats have them) but not suspiciously many in the header
	if nullCount > 128 {
		return false
	}

	return true
}

//error function
func ShowError(err error) string {
	if err == nil {
		return ""
	}

	if err.Error() == "not supported image type" {
		return "type d'image non pris en charge"
	}

	if err.Error() == "wrong credentials" {
		return "identifiants invalides"
	}

	// ─── Database ────────────────────────────────────────────────────────────
	if errors.Is(err, sql.ErrNoRows) {
		return "introuvable"
	}

	//postgres
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique violation
			return "existe déjà"
		case "23503": // foreign key violation
			return "l'enregistrement référencé n'existe pas"
		case "23502": // not null violation
			return "champ requis manquant"
		case "22001": // string too long
			return "valeur trop longue"
		case "28P01": // wrong password
			return "identifiants invalides"
		case "3D000": // database does not exist
			return "base de données introuvable"
		case "08006": // connection failure
			return "impossible de se connecter à la base de données"
		}
	}

	// ─── SQLite ──────────────────────────────────────────────────────────────
	// modernc.org/sqlite doesn't expose as clean a typed error as pgx does,
	// so we match on the message text — SQLite's constraint error strings
	// are stable across versions.
	msg := err.Error()
	switch {
	case strings.Contains(msg, "UNIQUE constraint failed"):
		return "existe déjà"
	case strings.Contains(msg, "FOREIGN KEY constraint failed"):
		return "l'enregistrement référencé n'existe pas"
	case strings.Contains(msg, "NOT NULL constraint failed"):
		return "champ requis manquant"
	case strings.Contains(msg, "CHECK constraint failed"):
		return "valeur invalide"
	case strings.Contains(msg, "database is locked"):
		return "base de données occupée, veuillez réessayer"
	}

	// ─── Bcrypt ──────────────────────────────────────────────────────────────
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "identifiants invalides"
	}
	if errors.Is(err, bcrypt.ErrHashTooShort) {
		return "identifiants invalides"
	}
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return "mot de passe trop long"
	}

	// ─── Validation (go-playground/validator) ────────────────────────────────
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var msgs []string
		for _, e := range validationErrs {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				msgs = append(msgs, field+" est requis")
			case "email":
				msgs = append(msgs, field+" doit être une adresse e-mail valide")
			case "min":
				msgs = append(msgs, field+" est trop court")
			case "max":
				msgs = append(msgs, field+" est trop long")
			case "len":
				msgs = append(msgs, field+" a une longueur invalide")
			case "numeric":
				msgs = append(msgs, field+" doit être un nombre")
			case "alpha":
				msgs = append(msgs, field+" ne doit contenir que des lettres")
			case "alphanum":
				msgs = append(msgs, field+" ne doit contenir que des lettres et des chiffres")
			case "url":
				msgs = append(msgs, field+" doit être une URL valide")
			case "uuid":
				msgs = append(msgs, field+" doit être un UUID valide")
			case "oneof":
				msgs = append(msgs, field+" a une valeur invalide")
			case "eqfield":
				msgs = append(msgs, field+" ne correspond pas")
			case "gtefield":
				msgs = append(msgs, field+" est hors limites")
			case "lt":
				msgs = append(msgs, field+" est trop grand")
			case "gt":
				msgs = append(msgs, field+" est trop petit")
			case "lte":
				msgs = append(msgs, field+" est trop grand")
			case "gte":
				msgs = append(msgs, field+" est trop petit")
			default:
				msgs = append(msgs, field+" est invalide")
			}
		}
		return strings.Join(msgs, ", ")
	}

	// ─── Fiber binding ───────────────────────────────────────────────────────
	if errors.Is(err, fiber.ErrUnprocessableEntity) {
		return "données de formulaire invalides"
	}
	if errors.Is(err, fiber.ErrBadRequest) {
		return "requête invalide"
	}

	// ─── File upload ─────────────────────────────────────────────────────────
	if errors.Is(err, fiber.ErrRequestEntityTooLarge) {
		return "fichier trop volumineux"
	}
	if strings.Contains(err.Error(), "no such file") {
		return "fichier introuvable"
	}
	if strings.Contains(err.Error(), "multipart") {
		return "téléversement de fichier invalide"
	}

	// ─── Session ─────────────────────────────────────────────────────────────
	if strings.Contains(err.Error(), "session") {
		return "erreur de session, veuillez vous reconnecter"
	}

	// ─── Redis ───────────────────────────────────────────────────────────────
	if strings.Contains(err.Error(), "redis") || strings.Contains(err.Error(), "EOF") {
		return "cache indisponible, veuillez réessayer"
	}

	// ─── Fallback ────────────────────────────────────────────────────────────
	return "une erreur s'est produite"
}