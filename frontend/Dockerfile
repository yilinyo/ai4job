FROM node:18-alpine AS build

RUN apk add --no-cache python3 make g++

WORKDIR /build/frontend
COPY frontend/package.json frontend/package-lock.json frontend/.npmrc ./
RUN npm ci --legacy-peer-deps

COPY frontend .
RUN npm run build:prod

FROM nginx:1.27-alpine

COPY --from=build /build/frontend/dist /usr/share/nginx/html
COPY deploy/nginx/default.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
