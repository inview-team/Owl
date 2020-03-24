FROM node:lts-alpine
RUN npm install json-server -g
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
ENTRYPOINT npm run build
