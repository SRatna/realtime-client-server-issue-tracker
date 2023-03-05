# frontend build
FROM node:18-alpine AS frontendBuild

WORKDIR /app

COPY frontend/package.json ./
COPY frontend/package-lock.json ./

RUN npm install

COPY ./frontend ./

RUN npm run build

# backend build and run
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /task-processor

COPY --from=frontendBuild /app/dist ./dist

EXPOSE 3000

CMD [ "/task-processor" ]