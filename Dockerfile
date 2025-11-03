FROM golang:latest

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main .

# Creamos un directorio para montar los certificados
RUN mkdir -p /app/certs

EXPOSE 8080

# CMD por defecto sigue siendo ./main
# En tu código Go asegúrate de apuntar a /app/certs/fullchain.pem y /app/certs/privkey.pem
CMD ["./main"]
