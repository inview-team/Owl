FROM node:lts-alpine
RUN npm install json-server -g
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
EXPOSE 3000
