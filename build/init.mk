.PHONY: init
init: .env # Initialize

.env: .env.sample # Create local .env
	cp -i .env.sample .env
