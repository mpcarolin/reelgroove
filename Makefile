dev:
	set -a && source .env && set +a && /Users/mpcarolin/go/bin/air .

templ:
	cd server && make templ

dlv:
	cd server && /Users/mpcarolin/go/bin/dlv debug