dev:
	set -a && source .env && set +a && /Users/mpcarolin/go/bin/air .

templ:
	cd server && make templ