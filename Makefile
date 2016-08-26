all: doc

GITHUB_PATH:=https://github.com/CreatorKit/go-deviceserver-client/blob/master
doc:
	godoc2md github.com/CreatorKit/go-deviceserver-client \
	| sed -e 's,/src/github.com/CreatorKit/go-deviceserver-client,$(GITHUB_PATH),g' \
	      -e 's,/src/target/\([a-zA-Z._]*\)?s=\([0-9]*\):\([0-9]*\)#\(L[0-9]*\),$(GITHUB_PATH)/\1#\4,g' \
	| perl -pe 's/([0-9]+)/($$1+10)/e' > README.md
