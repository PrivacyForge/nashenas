FROM node:20

RUN npm install -g pnpm
WORKDIR /app
COPY package.json ./
COPY pnpm-lock.yaml ./
RUN pnpm i
COPY . .

EXPOSE 8000

RUN pnpm build 
CMD pnpm start