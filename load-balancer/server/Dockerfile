FROM node:20.9.0-alpine3.18

WORKDIR /app/

COPY package*.json ./

COPY tsconfig.json ./

COPY ./src ./src

RUN npm install

RUN npm run build

RUN rm -rf node_modules

EXPOSE 3000

CMD ["node", "dist/index.js"]
