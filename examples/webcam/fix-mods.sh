#! /bin/bash
# Run the following once to pull correct dependencies
go get github.com/vladimirvivien/go4vl@latest
go get github.com/esimov/pigo/core@latest
go get github.com/fogleman/gg@8febc0f526adecda6f8ae80f3869b7cd77e52984

go mod tidy