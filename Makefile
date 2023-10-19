# .PHONY: start

# start:
# 	@echo
# 	/home/luka/go/bin/air -c .air.toml
# 	@echo
# 	npm run watch

start:
	make -j 3 tailwind-build tailwind-watch air
.PHONY: local

air:
	air -c .air.toml
.PHONY: air

tailwind-build:
	npm run build
.PHONY: tailwind-build

tailwind-watch:
	npm run watch
.PHONY: tailwind-watch