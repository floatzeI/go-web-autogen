go build;
swag init;
cp docs/swagger.json ./public/docs/docs.json;
rm docs/docs.go;