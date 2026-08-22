
.PHONY: dev devy start

project := crow

start:
	mkdir -p dev dev/glitter frontend frontend/glitter frontend/glitter/media frontend/glitter/media/dynamic frontend/glitter/media/static frontend/mark frontend/mark/comp frontend/mark/pages
	cd dev && npm init -y 
	(cd dev && npm install --save-dev vite tailwindcss@latest @tailwindcss/vite@latest daisyui@latest && npm install alpinejs && npm install htmx.org@2.0.8)
	mv vite.config.js dev/
	mv dev.js dev/glitter
	mv dev.css dev/glitter
	mv Layout.templ frontend/mark/comp
	mv welcome.templ frontend/mark/pages
	mv alpinetest.templ frontend/mark/comp
	touch frontend/glitter/app.css
	touch frontend/glitter/app.js

	go mod init $(project)
	go get github.com/gofiber/fiber/v3
	go get -tool github.com/a-h/templ/cmd/templ@latest
	go tool templ generate
	go get
	air init


devy:
	@trap 'kill 0' EXIT; \
	go tool templ generate --watch & \
	(cd dev && npx esbuild ./glitter/dev.js --bundle --outfile=./../frontend/glitter/app.js --watch=forever) & \
	air

dev:
	@trap 'kill 0' EXIT; \
	go tool templ generate --watch & \
	(cd dev && npx vite build --watch) & \
	air