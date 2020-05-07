#Dockerfile
#to build image: docker build -t frontend .
#to run: docker run -it --rm -p 8083:8080 frontend
FROM node:alpine AS builder
WORKDIR /app
COPY *.js *.json ./
COPY vue ./vue
RUN yarn install
RUN ./node_modules/.bin/webpack

#Dockerfile
FROM caddy:latest
COPY Caddyfile /etc/caddy/Caddyfile
COPY index.html /site/
COPY --from=builder /app/dist /site/dist/
