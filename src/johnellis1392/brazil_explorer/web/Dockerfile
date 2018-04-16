FROM node:alpine
# FROM node:latest

WORKDIR /usr/src/app
# RUN npm install --global --loglevel silent node-sass
RUN apk add --update git curl

RUN mkdir -p node_modules/node-sass/vendor/linux-x64-51
RUN curl -L https://github.com/sass/node-sass/releases/download/v4.5.0/linux-x64-51_binding.node \
    -o node_modules/node-sass/vendor/linux-x64-51/binding.node

COPY package.json package.json
RUN npm install --loglevel silent

COPY . .
CMD [ "npm", "run", "dev" ]
