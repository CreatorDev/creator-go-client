all: doc

GITHUB_PATH:=https://github.com/CreatorKit/go-deviceserver-client/blob/master
GITHUB_HATEOAS_PATH:=https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas

doc:
	godoc2md github.com/CreatorKit/go-deviceserver-client \
	| sed -e 's,/src/github.com/CreatorKit/go-deviceserver-client,$(GITHUB_PATH),g' \
	      -e 's,/src/target/\([a-zA-Z._]*\)?s=\([0-9]*\):\([0-9]*\)#\(L[0-9]*\),$(GITHUB_PATH)/\1#\4,g' \
	| perl -pe 's/([0-9]+)/($$1+10)/e' > README.md
	godoc2md github.com/CreatorKit/go-deviceserver-client/hateoas \
	| sed -e 's,/src/github.com/CreatorKit/go-deviceserver-client/hateoas,$(GITHUB_HATEOAS_PATH),g' \
	      -e 's,/src/target/\([a-zA-Z._]*\)?s=\([0-9]*\):\([0-9]*\)#\(L[0-9]*\),$(GITHUB_HATEOAS_PATH)/\1#\4,g' \
	| perl -pe 's/([0-9]+)/($$1+10)/e' > hateoas/README.md
