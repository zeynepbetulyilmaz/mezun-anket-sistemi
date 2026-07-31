# 1. Aþama: React/Vite uygulamasýný build etme
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# 2. Aþama: Nginx ile yayýnlama
FROM nginx:alpine
# Özel Nginx ayarýmýzý kopyalýyoruz
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Nginx'in varsayýlan "Welcome" sayfasýný siliyoruz
RUN rm -rf /usr/share/nginx/html/*

# ÝÞTE EKSÝK OLAN SATIR: React dosyalarýný Nginx'in içine aktarýyoruz
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]