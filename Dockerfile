FROM golang:1.23

WORKDIR /app

COPY views views
COPY server server
COPY database database
COPY routes routes

COPY tailwind.config.js ./
COPY go.mod go.sum  ./
COPY Makefile ./

# não é legal
COPY .env ./

RUN wget https://github.com/tailwindlabs/tailwindcss/releases/download/v4.0.3/tailwindcss-linux-x64 \
        --quiet \
        -O tailwindcss
RUN chmod a+x tailwindcss

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate .

# make não consegue encontrar o tailwind sem que ajustemos o PATH
RUN PATH=$PATH:/app make build

EXPOSE 8888

CMD ["/app/base"]
