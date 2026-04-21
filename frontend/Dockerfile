# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY . .

RUN npm install
RUN npm test

FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf

COPY --from=builder /app /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
