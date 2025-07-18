FROM eclipse-temurin:17-jdk-alpine

WORKDIR /opt/holodeckb2b

# Salin semua isi folder project ke dalam container
COPY . .

# Pastikan startServer.sh bisa dieksekusi
RUN chmod +x bin/startServer.sh

# Salin entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# (Opsional) Buat folder data/log jika perlu
RUN mkdir -p data msg_in msg_out logs

# Expose port aplikasi
EXPOSE 8080

# Gunakan entrypoint custom agar bisa copy file default ke volume host
ENTRYPOINT ["/entrypoint.sh"]
CMD ["sh", "bin/startServer.sh"]